// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package treeheap

import (
	"github.com/mrtcx/plusdata/heap"
	"github.com/mrtcx/plusdata/tree"
)

var _ heap.Heap = (*TreeHeap)(nil)

type TreeHeap struct {
	top  interface{}
	tree tree.Tree
}

func New(tree tree.Tree) *TreeHeap {
	return &TreeHeap{
		tree: tree,
	}
}

func (h *TreeHeap) Clean() {
	h.tree.Clean()
	h.top = nil
}

func (h *TreeHeap) Size() int {
	return h.tree.Size()
}

func (h *TreeHeap) Empty() bool {
	return h.tree.Empty()
}

func (h *TreeHeap) Top() interface{} {
	return h.top
}

func (h *TreeHeap) Push(val interface{}) {
	h.tree.Insert(val, nil)
	h.top = h.tree.Left().Key()
}

func (h *TreeHeap) Pop() interface{} {
	if h.top == nil {
		return nil
	}
	h.tree.Remove(h.top)
	top := h.top
	h.top = nil
	if left := h.tree.Left(); left != nil {
		h.top = left.Key()
	}
	return top
}
