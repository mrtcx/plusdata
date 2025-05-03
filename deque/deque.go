// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deque

type Deque interface {
	Size() int
	Empty() bool
	Clean()
	Get(int) interface{}
	Set(int, interface{})
	Front() interface{}
	Back() interface{}
	PushBack(interface{})
	PopBack() interface{}
	PushFront(value interface{})
	PopFront() interface{}
}
