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
	re1 = regexp.MustCompile(`Array\((.+)\)`) // Array() types
	r   = strings.NewReplacer(
		"Boolean", "bool", // boolean type
		"Bool", "bool", // boolean type
		"String", "string", // string type
		"{", "struct {", // object types
	)
	re2  = regexp.MustCompile(`\[\s*|]`)                   // array types
	re3  = regexp.MustCompile(`Int(\d+)`)                  // int types
	re4  = regexp.MustCompile(`Float(\d+)`)                // float types
	re5  = regexp.MustCompile(`: Number // Integer`)       // Number type
	re6  = regexp.MustCompile(`: (.+)\?`)                  // optional types
	re7  = regexp.MustCompile(`: "(.+)"(?: // Constant)?`) // string constants
	re8  = regexp.MustCompile(`: -?\d+`)                   // numeric literal values
	re9  = regexp.MustCompile(`: (?:true|false)`)          // boolean literal values
	re10 = regexp.MustCompile(`(?m): \[]\s*(?://.*)?$`)    // empty arrays
	re11 = regexp.MustCompile(`".+":`)                     // field names p1
	re12 = regexp.MustCompile(`".+:"`)                     // field names p2
	re13 = regexp.MustCompile(`\n?\s*\n`)                  // excessive newlines
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	s := bytes2string(data)

	s = strings.ReplaceAll(s, ",", "")
	s = re1.ReplaceAllString(s, `[]$1`)
	s = r.Replace(s)
	s = re2.ReplaceAllStringFunc(s, func(s string) string {
		switch s[0] {
		case '[':
			return "[]"
		case ']':
			return ""
		default:
			panic("unreachable")
		}
	})
	s = re3.ReplaceAllString(s, `int$1`)
	s = re4.ReplaceAllString(s, `float$1`)
	s = re5.ReplaceAllLiteralString(s, ": Number")
	s = strings.ReplaceAll(s, "Number", "int64")
	s = strings.ReplaceAll(s, ": int64 // Unix timestamp", ": time.Time")
	s = re6.ReplaceAllString(s, `: option.Option[$1]`)
	s = re7.ReplaceAllString(s, `: string // "$1"`)
	s = re8.ReplaceAllString(s, `: int64`)
	s = re9.ReplaceAllString(s, `: bool`)
	s = re10.ReplaceAllLiteralString(s, ": []any")
	s = strings.ReplaceAll(s, ": null", ": any")
	s = re11.ReplaceAllStringFunc(s, func(s string) string {
		runes := []rune(s)
		runes = runes[1 : len(runes)-2]
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	})
	s = re12.ReplaceAllStringFunc(s, func(s string) string {
		runes := []rune(s)
		runes = runes[1 : len(runes)-3]
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	})
	s = re13.ReplaceAllLiteralString(s, "\n")

	data, err = format.Source(string2bytes(s))
	if err != nil {
		os.Stderr.WriteString(err.Error()) //nolint:errcheck
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
