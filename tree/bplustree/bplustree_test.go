// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bplustree

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/mrtcx/plusdata/internal/assert"
	"github.com/mrtcx/plusdata/tree"
)

var intcmp = tree.IntComparator

func TestInsert(t *testing.T) {
	nums := []int{1, 2, 4, 8, 9, 1024, 1024 + 1}
	orders := []int{3, 4, 5}
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
					assert.Equal(t, checkBalance(t, rb, rb.root), true)
					assert.Equal(t, checkOrder(t, rb, rb.root), true)
				}
				assert.Equal(t, rb.size, tnum)

				rb = New(intcmp, torder)
				for i := tnum; i >= 1; i-- {
					rb.Insert(i, i)
					rb.Insert(i, i)
					assert.Equal(t, checkBalance(t, rb, rb.root), true)
					assert.Equal(t, checkOrder(t, rb, rb.root), true)
				}
				assert.Equal(t, rb.size, tnum)

				rb = New(intcmp, torder)
				for i, j := 1, tnum; i <= j; i, j = i+1, j-1 {
					rb.Insert(i, i)
					rb.Insert(j, j)
					assert.Equal(t, checkBalance(t, rb, rb.root), true)
					assert.Equal(t, checkOrder(t, rb, rb.root), true)
				}
				assert.Equal(t, rb.size, tnum)

				rb = New(intcmp, torder)
				for i := 1; i <= tnum; i++ {
					randnum := int(rand.Int31() % int32(tnum))
					rb.Insert(randnum, randnum)
					rb.Insert(randnum*-1, randnum*-1)
					assert.Equal(t, checkBalance(t, rb, rb.root), true)
					assert.Equal(t, checkOrder(t, rb, rb.root), true)
				}
			})
		}
	}
}

func TestRemove(t *testing.T) {
	nums := []int{1, 2, 3, 4, 8, 1024, 1024 + 1}
	orders := []int{3, 4, 5, 8}

	for _, num := range nums {
		tnum := num
		for _, order := range orders {
			torder := order
			t.Run(fmt.Sprintf("[num:%d order:%d]", tnum, torder), func(t *testing.T) {
				rb := newBpNum(tnum, torder)
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
				assert.Equal(t, rb.size, 0)

				rb = newBpNum(tnum, torder)
				for i := tnum; i >= 1; i-- {
					rb.Remove(i)
					assert.Equal(t, checkBalance(t, rb, rb.root), true)
					assert.Equal(t, checkOrder(t, rb, rb.root), true)
					assert.Equal(t, rb.size, i-1)
				}
				assert.Equal(t, rb.root, nil)
				assert.Equal(t, rb.size, 0)

				rb = newBpNum(tnum, torder)
				left, right := tnum/2, tnum/2+1
				for ; left >= 1 && right <= tnum; left, right = left-1, right+1 {
					rb.Remove(left)
					rb.Remove(right)
					assert.Equal(t, checkBalance(t, rb, rb.root), true)
					assert.Equal(t, checkOrder(t, rb, rb.root), true)
					assert.Equal(t, rb.size, tnum-(right-left+1))
				}
				rb.Remove(left)
				rb.Remove(right)
				assert.Equal(t, rb.root, nil)
				assert.Equal(t, rb.size, 0)
			})
		}
	}
}

func newBpNum(insertNum, testorder int) *bplusTree {
	rb := New(intcmp, testorder)
	for i := 1; i <= insertNum; i++ {
		rb.Insert(i, i)
	}
	return rb
}

func TestGet(t *testing.T) {
	orders := []int{3, 4, 5}
	num := 1024
	for _, v := range orders {
		testOrder := v
		t.Run(fmt.Sprintf("order:%d", testOrder), func(t *testing.T) {
			rb := New(intcmp, testOrder)
			geti, ei := rb.Get(1)
			assert.Equal(t, geti, nil)
			assert.Equal(t, ei, false)
			for i := 1; i <= num; i++ {
				rb.Insert(i, i)
				geti, ei = rb.Get(i)
				assert.Equal(t, geti, i)
				assert.Equal(t, ei, true)
				geti, ei = rb.Get(i + 1)
				assert.Equal(t, geti, nil)
				assert.Equal(t, ei, false)
			}
			for i := 1; i <= num; i++ {
				rb.Remove(i)
				_, ei = rb.Get(i)
				assert.Equal(t, ei, false)
				for j := i + 1; j <= num; j++ {
					geti, ei = rb.Get(j)
					assert.Equal(t, ei, true)
					assert.Equal(t, geti, j)
				}
			}
		})
	}
}

func checkBalance(t *testing.T, bp *bplusTree, root *Node) bool {
	if root == nil {
		return root == bp.root && bp.size == 0
	}
	for i := 1; i < len(root.keys); i++ {
		if bp.cmp(root.keys[i], root.keys[i-1]) <= 0 {
			t.Logf("fail node:%v", root)
			return false
		}
	}
	if root.isLeaf() {
		if len(root.valus) != len(root.keys) || len(root.childs) != 0 {
			t.Logf("fail node:%v", root)
			return false
		}
		minkeys, maxkeys := bp.minKeys(), bp.maxKeys()
		if root == bp.root {
			minkeys = 1
		}
		if len(root.keys) < minkeys || len(root.keys) > maxkeys {
			t.Logf("fail node:%v", root)
			return false
		}
	} else {
		if len(root.childs)-len(root.keys) != 1 || len(root.valus) != 0 || root.next != nil {
			t.Logf("fail node:%v", root)
			return false
		}
		minkeys, maxkeys := bp.minKeys(), bp.maxKeys()
		if root == bp.root {
			minkeys = 1
		}
		if len(root.keys) < minkeys || len(root.keys) > maxkeys {
			t.Logf("fail node:%v", root)
			return false
		}
		for i := 0; i < len(root.keys); i++ {
			lchild := root.childs[i]
			rchild := root.childs[i+1]
			if lchild.isLeaf() {
				if !rchild.isLeaf() {
					t.Logf("fail node:%v child:%v", root, rchild)
					return false
				}
				if bp.cmp(lchild.keys[len(lchild.keys)-1], root.keys[i]) > 0 {
					t.Logf("fail node:%v child:%v", root, lchild)
					return false
				}
			} else {
				if rchild.isLeaf() {
					t.Logf("fail node:%v child:%v", root, rchild)
					return false
				}
				if bp.cmp(lchild.keys[len(lchild.keys)-1], root.keys[i]) >= 0 {
					t.Logf("fail node:%v child:%v", root, lchild)
					return false
				}
			}
			if bp.cmp(rchild.keys[0], root.keys[i]) <= 0 {
				t.Logf("fail node:%v child:%v", root, rchild)
				return false
			}
		}
		for i := 0; i < len(root.childs); i++ {
			child := root.childs[i]
			if !checkBalance(t, bp, child) {
				return false
			}
		}
	}
	return true
}

func checkOrder(t *testing.T, bp *bplusTree, root *Node) bool {
	var leftarr []interface{} = make([]interface{}, 0)
	var rightarr []interface{} = make([]interface{}, 0)
	leftToRightseqence(root, &leftarr)
	rightToLeftseqence(root, &rightarr)
	if len(leftarr) != bp.size {
		t.Logf("fail seq %d %d", bp.size, len(leftarr))
		return false
	}
	if len(rightarr) != bp.size {
		t.Logf("fail seq %d %d", bp.size, len(rightarr))
		return false
	}
	for i := 1; i < bp.size; i++ {
		if bp.cmp(leftarr[i], leftarr[i-1]) <= 0 {
			t.Logf("fail seq:%v", leftarr)
			return false
		}
		if bp.cmp(rightarr[i], rightarr[i-1]) >= 0 {
			t.Logf("fail seq:%v", leftarr)
			return false
		}
		if bp.cmp(leftarr[i], rightarr[len(rightarr)-1-i]) != 0 {
			t.Logf("fail seq:%v %v %v %v", rightarr, leftarr, rightarr[len(rightarr)-i], leftarr[i])
			return false
		}
	}
	return true
}

func leftToRightseqence(root *Node, arr *[]interface{}) {
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

func rightToLeftseqence(root *Node, arr *[]interface{}) {
	if root == nil {
		return
	}
	for !root.isLeaf() {
		root = root.childs[len(root.childs)-1]
	}
	for root != nil {
		for i := len(root.keys) - 1; i >= 0; i-- {
			*arr = append(*arr, root.keys[i])
		}
		root = root.prev
	}
}

func BenchmarkInsert(b *testing.B) {
	testsizes := []int{100000, 500000, 1000000, 2000000}
	testNames := []string{"10w", "50w", "100w", "200w"}
	for i, v := range testsizes {
		size := v
		name := testNames[i]
		b.Run(fmt.Sprintf("map-%s", name), func(b *testing.B) {
			tr := make(map[int]interface{})
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < size; j++ {
					tr[j] = nil
				}
			}
		})
		b.Run(fmt.Sprintf("bplustree-%s", name), func(b *testing.B) {
			tr := New(tree.IntComparator, 32)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < size; j++ {
					tr.Insert(j, nil)
				}
			}
		})
	}
}

func BenchmarkRemove(b *testing.B) {
	testsizes := []int{100000, 500000, 1000000, 2000000}
	testNames := []string{"10w", "50w", "100w", "200w"}
	for i, v := range testsizes {
		size := v
		name := testNames[i]
		b.Run(fmt.Sprintf("map-%s", name), func(b *testing.B) {
			tr := make(map[int]interface{})
			for j := 0; j < size; j++ {
				tr[j] = nil
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < size; j++ {
					delete(tr, j)
				}
			}
		})
		b.Run(fmt.Sprintf("bplustree-%s", name), func(b *testing.B) {
			tr := New(tree.IntComparator, 32)
			for j := 0; j < size; j++ {
				tr.Insert(j, nil)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < size; j++ {
					tr.Remove(j)
				}
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	testsizes := []int{100000, 500000, 1000000, 2000000}
	testNames := []string{"10w", "50w", "100w", "200w"}
	for i, v := range testsizes {
		size := v
		name := testNames[i]
		b.Run(fmt.Sprintf("map-%s", name), func(b *testing.B) {
			tr := make(map[int]interface{})
			for j := 0; j < size; j++ {
				tr[j] = nil
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < size; j++ {
					_, _ = tr[j]
				}
			}
		})
		b.Run(fmt.Sprintf("bplustree-%s", name), func(b *testing.B) {
			tr := New(tree.IntComparator, 32)
			for j := 0; j < size; j++ {
				tr.Insert(j, nil)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < size; j++ {
					_, _ = tr.Get(j)
				}
			}
		})
	}
}
