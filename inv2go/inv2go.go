// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"go/format"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func main() {
	output, err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if output != "" {
		fmt.Fprintln(os.Stdout, output) //nolint:errcheck
	}
}

func run() (string, error) {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}

	processed := string(input)
	processed = r.Replace(processed)
	processed = re1.ReplaceAllString(processed, `: []$1`)
	processed = re2.ReplaceAllString(processed, `: []$1`)
	processed = re3.ReplaceAllString(processed, `[]$1`)
	processed = re4.ReplaceAllStringFunc(processed, func(s string) string {
		switch s[0] {
		case '[':
			return "[]"
		case ']':
			return ""
		default:
			panic("unreachable")
		}
	})
	processed = re5.ReplaceAllString(processed, `int$1`)
	processed = re6.ReplaceAllString(processed, `float$1`)
	processed = strings.ReplaceAll(processed, ": Number // Integer", ": Number")
	processed = strings.ReplaceAll(processed, "Number", "int64")
	processed = strings.ReplaceAll(processed, ": int64 // Unix timestamp", ": time.Time")
	processed = re7.ReplaceAllString(processed, `: option.Option[$1]`)
	processed = re8.ReplaceAllString(processed, `: string // "$1"`)
	processed = re9.ReplaceAllLiteralString(processed, ": int64")
	processed = re10.ReplaceAllLiteralString(processed, ": bool")
	processed = re11.ReplaceAllLiteralString(processed, ": []any")
	processed = strings.ReplaceAll(processed, ": null", ": any")
	processed = re12.ReplaceAllStringFunc(processed, func(s string) string {
		runes := []rune(s)
		runes = runes[1 : len(runes)-2]
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	})
	processed = re13.ReplaceAllStringFunc(processed, func(s string) string {
		runes := []rune(s)
		runes = runes[1 : len(runes)-3]
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	})
	processed = re14.ReplaceAllLiteralString(processed, "\n")

	formatted, err := format.Source([]byte(processed))
	if err != nil {
		return processed, err
	}
	return string(formatted), nil
}

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
	re7  = regexp.MustCompile(`: (\w+)\?`)                       // optional types
	re8  = regexp.MustCompile(`: "(.+)"(?: // Constant)?`)       // string constants
	re9  = regexp.MustCompile(`: -?\d+`)                         // numeric literal values
	re10 = regexp.MustCompile(`: (?:true|false)`)                // boolean literal values
	re11 = regexp.MustCompile(`(?m): \[]\s*(?://.*)?$`)          // empty arrays
	re12 = regexp.MustCompile(`".+":`)                           // field names p1
	re13 = regexp.MustCompile(`".+:"`)                           // field names p2 (typo)
	re14 = regexp.MustCompile(`\n?\s*\n`)                        // excessive newlines
)
