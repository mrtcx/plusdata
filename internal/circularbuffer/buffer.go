// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package circularbuffer

import "fmt"

type Buffer struct {
	start int
	end   int
	size  int
	array []interface{}
}

func New(capacity int) *Buffer {
	return &Buffer{
		start: 0,
		end:   0,
		array: make([]interface{}, capacity),
	}
}

func (b *Buffer) arraryIndex(point int) int {
	idx := (point%len(b.array) + len(b.array)) % len(b.array)
	return idx
}

func (b *Buffer) Size() int {
	return b.size
}

func (b *Buffer) Capacity() int {
	return len(b.array)
}

func (b *Buffer) IsFull() bool {
	return b.size == len(b.array)
}

func (b *Buffer) IsEmpty() bool {
	return b.size == 0
}

func (b *Buffer) Back() interface{} {
	if b.size == 0 {
		return nil
	}
	return b.array[b.arraryIndex(b.end)]
}

func (b *Buffer) Front() interface{} {
	if b.size == 0 {
		return nil
	}
	return b.array[b.arraryIndex(b.start)]
}

func (b *Buffer) PushFront(value interface{}) {
	if b.size == len(b.array) {
		panic("is full")
	}
	if b.size != 0 {
		b.start--
	}
	b.array[b.arraryIndex(b.start)] = value
	b.size++
}

func (b *Buffer) PopFront() interface{} {
	if b.size == 0 {
		return nil
	}
	idx := b.arraryIndex(b.start)
	pv := b.array[idx]
	b.array[idx] = nil
	if b.size > 1 {
		b.start++
	}
	b.size--
	return pv
}

func (b *Buffer) PushBack(value interface{}) {
	if b.size == len(b.array) {
		panic("is full")
	}
	if b.size != 0 {
		b.end++
	}
	b.array[b.arraryIndex(b.end)] = value
	b.size++
}

func (b *Buffer) PopBack() interface{} {
	if b.size == 0 {
		return nil
	}
	idx := b.arraryIndex(b.end)
	pv := b.array[idx]
	b.array[idx] = nil
	if b.size > 1 {
		b.end--
	}
	b.size--
	return pv
}

func (b *Buffer) Get(index int) interface{} {
	if index >= b.size || index < 0 {
		panic(fmt.Sprintf("index[%d] beyond bound [%d:%d)", index, 0, b.size))
	}
	return b.array[b.arraryIndex(b.start+index)]
}

func (b *Buffer) Set(index int, val interface{}) {
	if index >= b.size || index < 0 {
		panic(fmt.Sprintf("index[%d] beyond bound [%d:%d)", index, 0, b.size))
	}
	b.array[b.arraryIndex(b.start+index)] = val
}

func (b *Buffer) ResetCapacity(capacity int) {
	if capacity < b.size {
		panic(fmt.Sprintf("cpacity[%d] less size[%d]", capacity, b.size))
	}
	if capacity == len(b.array) {
		return
	}
	newarrary := make([]interface{}, capacity)
	oldarrary := b.array
	sidx, eidx := b.arraryIndex(b.start), b.arraryIndex(b.end)
	if sidx <= eidx {
		copy(newarrary[0:b.size], oldarrary[sidx:sidx+b.size])
		b.start, b.end = 0, eidx-sidx
		b.array = newarrary
		return
	} else {
		toplen := len(oldarrary) - sidx
		copy(newarrary[len(newarrary)-toplen:], b.array[sidx:sidx+toplen])
		copy(newarrary[:b.size-toplen], oldarrary[:b.size-toplen])
		b.array = newarrary
		b.start, b.end = len(newarrary)-toplen, eidx
		return
	}
}
