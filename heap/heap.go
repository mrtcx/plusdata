// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package heap

import (
	"strings"
)

type Less func(a, b interface{}) bool

var IntLess = func(top, bottom interface{}) bool {
	return top.(int) < bottom.(int)
}

var StringLess = func(top, bottom interface{}) bool {
	return strings.Compare(top.(string), bottom.(string)) < 0
}

type Heap interface {
	Size() int
	Empty() bool
	Clean()
	Top() interface{}
	Push(interface{})
	Pop() interface{}
}
