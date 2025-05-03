// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/mrtcx/plusdata/heap"
	"github.com/mrtcx/plusdata/heap/arraryheap"
	"github.com/mrtcx/plusdata/heap/treeheap"
	"github.com/mrtcx/plusdata/tree"
	"github.com/mrtcx/plusdata/tree/rbtree"
)

func main() {
	//数组堆，堆内元素不去重
	var heap heap.Heap = arraryheap.New(heap.IntLess)
	heap.Push(4)
	heap.Push(2)
	heap.Push(2)
	heap.Push(3)
	heap.Push(1)
	heap.Size() //5
	for !heap.Empty() {
		_ = heap.Top().(int) //1, 2, 2, 3, 4, 数组堆内元素没有去重
		heap.Pop()
	}

	//排序树堆
	heap = treeheap.New(rbtree.New(tree.IntComparator)) //使用红黑树
	heap.Push(4)
	heap.Push(2)
	heap.Push(2)
	heap.Push(3)
	heap.Push(1)
	heap.Size() //4
	for !heap.Empty() {
		heap.Top() //1, 2, 3, 4, 排序树堆内元素会去重
		heap.Pop()
	}

	heap.Push(1) //[1]
	heap.Clean() //[]
	heap.Size()  //0
}
