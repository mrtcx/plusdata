// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bplustree

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/mrtcx/plusdata/internal/assert"
)

var intcmp func(i, j interface{}) int = func(i, j interface{}) int {
	return i.(int) - j.(int)
}

func TestInsert(t *testing.T) {
	nums := []int{1, 2, 4, 8, 1024}
	orders := []int{3}
	for _, num := range nums {
		tnum := num
		for _, order := range orders {
			torder := order
			t.Run(fmt.Sprintf("[num:%d order:%d]", tnum, torder), func(t *testing.T) {
				rb := New(intcmp, torder)
				assert.Equal(t, rb.root, nil)
				for i := 1; i <= tnum; i++ {
					rb.Insert(i, i)
					rb.Insert(i, i)
				}
				assert.Equal(t, checkBalance(t, rb, rb.root), true)
				assert.Equal(t, checkOrder(t, rb, rb.root), true)
				assert.Equal(t, rb.size, tnum)

				rb = New(intcmp, torder)
				for i := tnum; i >= 1; i-- {
					rb.Insert(i, i)
					rb.Insert(i, i)
				}
				assert.Equal(t, checkBalance(t, rb, rb.root), true)
				assert.Equal(t, checkOrder(t, rb, rb.root), true)
				assert.Equal(t, rb.size, tnum)

				rb = New(intcmp, torder)
				for i, j := 1, tnum; i <= j; i, j = i+1, j-1 {
					rb.Insert(i, i)
					rb.Insert(j, j)
				}
				assert.Equal(t, checkBalance(t, rb, rb.root), true)
				assert.Equal(t, checkOrder(t, rb, rb.root), true)
				assert.Equal(t, rb.size, tnum)

				rb = New(intcmp, torder)
				for i := 1; i <= tnum; i++ {
					randnum := int(rand.Int31() % int32(tnum))
					rb.Insert(randnum, randnum)
					rb.Insert(randnum*-1, randnum*-1)
				}
				assert.Equal(t, checkBalance(t, rb, rb.root), true)
				assert.Equal(t, checkOrder(t, rb, rb.root), true)
			})
		}
	}
}

func TestRemove(t *testing.T) {
	nums := []int{4}
	orders := []int{3}
	for _, num := range nums {
		tnum := num
		for _, order := range orders {
			torder := order
			t.Run(fmt.Sprintf("[num:%d order:%d]", tnum, torder), func(t *testing.T) {
				rb := New(intcmp, torder)
				for i := 1; i <= tnum; i++ {
					rb.Insert(i, i)
				}
				rb.Remove(0)
				assert.Equal(t, checkBalance(t, rb, rb.root), true)
				assert.Equal(t, checkOrder(t, rb, rb.root), true)
				assert.Equal(t, rb.size, tnum)
				for i := 1; i <= tnum; i++ {
					rb.Remove(i)
					assert.Equal(t, checkBalance(t, rb, rb.root), true)
					assert.Equal(t, checkOrder(t, rb, rb.root), true)
					assert.Equal(t, rb.size, tnum-i)
				}
				assert.Equal(t, rb.root, nil)

				rb = New(intcmp, torder)
				for i := 1; i <= tnum; i++ {
					rb.Insert(i, i)
				}
				for i := tnum; i >= 1; i-- {
					rb.Remove(i)
					assert.Equal(t, checkBalance(t, rb, rb.root), true)
					assert.Equal(t, checkOrder(t, rb, rb.root), true)
					assert.Equal(t, rb.size, i-1)
				}
				assert.Equal(t, rb.root, nil)

				rb = New(intcmp, torder)
				for i := 1; i <= tnum; i++ {
					rb.Insert(i, i)
				}
				for left, right := tnum/2, tnum/2+1; left >= 1 && right <= tnum; left, right = left-1, right+1 {
					rb.Remove(left)
					rb.Remove(right)
					assert.Equal(t, checkBalance(t, rb, rb.root), true)
					assert.Equal(t, checkOrder(t, rb, rb.root), true)
					assert.Equal(t, rb.size, tnum-(right-left+1))
				}
			})
		}
	}
}

func TestGet(t *testing.T) {
	orders := []int{3, 4, 5}
	for _, v := range orders {
		order := v
		t.Run(fmt.Sprintf("order:%d", order), func(t *testing.T) {
			rb := New(intcmp, order)
			geti, ei := rb.Get(1)
			assert.Equal(t, geti, nil)
			assert.Equal(t, ei, false)
			for i := 1; i < 8; i++ {
				rb.Insert(i, i)
				geti, ei = rb.Get(i)
				assert.Equal(t, geti, i)
				assert.Equal(t, ei, true)
				geti, ei = rb.Get(i + 1)
				assert.Equal(t, geti, nil)
				assert.Equal(t, ei, false)
			}
		})
	}
}

func checkBalance(t *testing.T, bp *BPlusTree, root *Node) bool {
	if root == bp.root {
		if root == nil {
			return bp.size == 0
		}
		if !root.isLeaf() {
			if len(root.childs) < 2 || len(root.childs) > int(bp.maxKeys()) || len(root.childs) != len(root.keys) {
				t.Logf("fail node:%v", root)
				return false
			}
		} else {
			if len(root.keys) > int(bp.maxKeys()) || len(root.keys) != len(root.vals) {
				t.Logf("fail node:%v", root)
				return false
			}
		}
		for i := 0; i < len(root.childs); i++ {
			child := root.childs[i]
			if bp.cmp(child.keys[len(child.keys)-1], root.keys[i]) != 0 {
				t.Logf("fail node:%v", root)
				return false
			}
			if !checkBalance(t, bp, child) {
				return false
			}
		}
		return true
	}
	if len(root.keys) < int(bp.minKeys()) || len(root.keys) > int(bp.maxKeys()) {
		t.Logf("fail node:%v", root)
		return false
	}
	if root.isLeaf() {
		if len(root.childs) != 0 || len(root.vals) != len(root.keys) {
			t.Logf("fail node:%v", root)
			return false
		}
		return true
	} else {
		if len(root.vals) != 0 || len(root.childs) != len(root.keys) {
			t.Logf("fail node:%v", root)
			return false
		}
		for i := 0; i < len(root.childs); i++ {
			child := root.childs[i]
			if bp.cmp(child.keys[len(child.keys)-1], root.keys[i]) != 0 {
				t.Logf("fail node:%v child:%v", root, child)
				return false
			}
			if !checkBalance(t, bp, child) {
				return false
			}
		}
		return true
	}
}

func checkOrder(t *testing.T, bp *BPlusTree, root *Node) bool {
	var arr []interface{} = make([]interface{}, 0)
	seqence(root, &arr)
	for i := 1; i < bp.size; i++ {
		if bp.cmp(arr[i], arr[i-1]) < 0 {
			t.Logf("fail seq:%v", arr)
			return false
		}
	}
	return true
}

func seqence(root *Node, arr *[]interface{}) {
	if root == nil {
		return
	}
	for !root.isLeaf() {
		root = root.childs[0]
	}
	for root != nil {
		*arr = append(*arr, root.keys...)
		root = root.next
	}
}
