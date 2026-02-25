// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package invidious

import (
	"strconv"
	"unsafe"
)

func itoa[T ~int | ~int8 | ~int16 | ~int32 | ~int64](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

func appendQuote(dst []byte, s string) []byte {
	if strconv.CanBackquote(s) {
		dst = append(dst, '`')
		dst = append(dst, s...)
		dst = append(dst, '`')
		return dst
	} else {
		return strconv.AppendQuoteToGraphic(dst, s)
	}
}

func quotedLen(s string) int {
	return 1 + len(s) + 1 // best-case scenario
}

func string2bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func bytes2string(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
