// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package arraryheap

import (
	"github.com/mrtcx/plusdata/arrary"
	"github.com/mrtcx/plusdata/arrary/blockslices"
	"github.com/mrtcx/plusdata/heap"
)

var _ heap.Heap = (*ArraryHeap)(nil)

type ArraryHeap struct {
	container arrary.Arrary
	cmp       heap.Less
}

func New(cmp heap.Less) *ArraryHeap {
	return &ArraryHeap{
		container: blockslices.New(),
		cmp:       cmp,
	}
}

func (h *ArraryHeap) Clean() {
	h.container.Clean()
}

func (h *ArraryHeap) Size() int {
	return h.container.Size()
}

func (h *ArraryHeap) Empty() bool {
	return h.container.Empty()
}

func (h *ArraryHeap) Top() interface{} {
	if h.container.Empty() {
		return nil
	}
	return h.container.Get(0)
}

func (h *ArraryHeap) Push(val interface{}) {
	h.container.PushBack(val)
	h.downToTop()
}

func (h *ArraryHeap) Pop() interface{} {
	if h.container.Empty() {
		return nil
	}
	head := h.container.Get(0)
	tail := h.container.PopBack()
	if !h.container.Empty() {
		h.container.Set(0, tail)
		h.topToDown()
	}
	return head
}

func (h *ArraryHeap) topToDown() {
	parent := 0
	size := h.container.Size()
	for parent*2+1 < size {
		lchild := parent*2 + 1
		rchild := parent*2 + 2
		parentVal := h.container.Get(parent)
		lchildval := h.container.Get(lchild)
		if rchild >= size || h.cmp(lchildval, h.container.Get(rchild)) {
			if h.cmp(lchildval, parentVal) {
				h.container.Set(lchild, parentVal)
				h.container.Set(parent, lchildval)
				parent = lchild
				continue
			}
		} else {
			rchildval := h.container.Get(rchild)
			if h.cmp(rchildval, parentVal) {
				h.container.Set(rchild, parentVal)
				h.container.Set(parent, rchildval)
				parent = rchild
				continue
			}
		}
		break
	}
}
func (h *ArraryHeap) downToTop() {
	child := h.Size() - 1
	for child > 0 {
		parent := (child - 1) / 2
		childval := h.container.Get(child)
		parentval := h.container.Get(parent)
		if h.cmp(childval, parentval) {
			h.container.Set(parent, childval)
			h.container.Set(child, parentval)
			child = parent
			continue
		}
		break
	}
}
