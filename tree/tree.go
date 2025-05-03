// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package tree

import (
	"strings"
)

type Comparator func(a, b interface{}) int

var IntComparator = func(left, right interface{}) int {
	return left.(int) - right.(int)
}

var StringComparator = func(left, right interface{}) int {
	return strings.Compare(left.(string), right.(string))
}

type Tree interface {
	Size() int
	Empty() bool
	Clean()
	Insert(key interface{}, value interface{})
	Remove(key interface{})
	Get(key interface{}) (interface{}, bool)
	Find(key interface{}) Element
	Left() Element
	Right() Element
	Prev(key interface{}) Element
	Next(key interface{}) Element
}

type Element interface {
	Key() interface{}
	Value() interface{}
	SetValue(value interface{})
	Next() Element
	Prev() Element
}
