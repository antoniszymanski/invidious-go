// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"go/format"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"
	"unsafe"
)

var (
	r = strings.NewReplacer(
		"Boolean", "bool", // boolean type
		"Bool", "bool", // boolean type
		"String", "string", // string type
		"{", "struct {", // object types
		",", "",
	)
	re1  = regexp.MustCompile(`: \[\s*// One or more (\w+)\s*]`) // array types p1
	re2  = regexp.MustCompile(`: (\w+)\[\]`)                     // array types p2
	re3  = regexp.MustCompile(`Array\((\w+)\)`)                  // array types p3
	re4  = regexp.MustCompile(`\[\s*|]`)                         // array types p4
	re5  = regexp.MustCompile(`Int(\d+)`)                        // int types
	re6  = regexp.MustCompile(`Float(\d+)`)                      // float types
	re7  = regexp.MustCompile(`: Number // Integer`)             // Number type
	re8  = regexp.MustCompile(`: (\w+)\?`)                       // optional types
	re9  = regexp.MustCompile(`: "(.+)"(?: // Constant)?`)       // string constants
	re10 = regexp.MustCompile(`: -?\d+`)                         // numeric literal values
	re11 = regexp.MustCompile(`: (?:true|false)`)                // boolean literal values
	re12 = regexp.MustCompile(`(?m): \[]\s*(?://.*)?$`)          // empty arrays
	re13 = regexp.MustCompile(`".+":`)                           // field names p1
	re14 = regexp.MustCompile(`".+:"`)                           // field names p2 (typo)
	re15 = regexp.MustCompile(`\n?\s*\n`)                        // excessive newlines
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	s := bytes2string(data)

	s = r.Replace(s)
	s = re1.ReplaceAllString(s, `: []$1`)
	s = re2.ReplaceAllString(s, `: []$1`)
	s = re3.ReplaceAllString(s, `[]$1`)
	s = re4.ReplaceAllStringFunc(s, func(s string) string {
		switch s[0] {
		case '[':
			return "[]"
		case ']':
			return ""
		default:
			panic("unreachable")
		}
	})
	s = re5.ReplaceAllString(s, `int$1`)
	s = re6.ReplaceAllString(s, `float$1`)
	s = re7.ReplaceAllLiteralString(s, ": Number")
	s = strings.ReplaceAll(s, "Number", "int64")
	s = strings.ReplaceAll(s, ": int64 // Unix timestamp", ": time.Time")
	s = re8.ReplaceAllString(s, `: option.Option[$1]`)
	s = re9.ReplaceAllString(s, `: string // "$1"`)
	s = re10.ReplaceAllString(s, `: int64`)
	s = re11.ReplaceAllString(s, `: bool`)
	s = re12.ReplaceAllLiteralString(s, ": []any")
	s = strings.ReplaceAll(s, ": null", ": any")
	s = re13.ReplaceAllStringFunc(s, func(s string) string {
		runes := []rune(s)
		runes = runes[1 : len(runes)-2]
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	})
	s = re14.ReplaceAllStringFunc(s, func(s string) string {
		runes := []rune(s)
		runes = runes[1 : len(runes)-3]
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	})
	s = re15.ReplaceAllLiteralString(s, "\n")

	data, err = format.Source(string2bytes(s))
	if err != nil {
		os.Stderr.WriteString(err.Error()) //nolint:errcheck
		os.Stderr.WriteString("\n")        //nolint:errcheck
		data = string2bytes(s)
	}
	os.Stdout.Write(data) //nolint:errcheck
}

func string2bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func bytes2string(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
