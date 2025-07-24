// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package invidious

import (
	"strconv"
	"unsafe"

	"golang.org/x/exp/constraints"
)

func itoa[T constraints.Signed](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

func quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	} else {
		return strconv.QuoteToGraphic(s)
	}
}

func string2bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func bytes2string(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
