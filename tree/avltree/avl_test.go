// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package avltree

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/mrtcx/plusdata/internal/assert"
	"github.com/mrtcx/plusdata/tree"
)

var intcmp = tree.IntComparator

func TestInsert(t *testing.T) {
	nums := []int{1, 2, 4, 8, 1024, 1024 + 1}
	for _, num := range nums {
		tnum := num
		t.Run(fmt.Sprintf("[num:%d]", tnum), func(t *testing.T) {
			rb := New(intcmp)
			assert.Equal(t, rb.root, _nil)
			for i := 1; i <= tnum; i++ {
				rb.Insert(i, i)
				rb.Insert(i, i)
				assert.Equal(t, checkHeightBalance(t, rb, rb.root), true)
				assert.Equal(t, checkOrder(t, rb, rb.root), true)
			}
			assert.Equal(t, rb.size, tnum)

			rb = New(intcmp)
			for i := tnum; i >= 1; i-- {
				rb.Insert(i, i)
				rb.Insert(i, i)
				assert.Equal(t, checkHeightBalance(t, rb, rb.root), true)
				assert.Equal(t, checkOrder(t, rb, rb.root), true)
			}
			assert.Equal(t, rb.size, tnum)

			rb = New(intcmp)
			for i, j := 1, tnum; i <= j; i, j = i+1, j-1 {
				rb.Insert(i, i)
				rb.Insert(j, j)
				assert.Equal(t, checkHeightBalance(t, rb, rb.root), true)
				assert.Equal(t, checkOrder(t, rb, rb.root), true)
			}
			assert.Equal(t, rb.size, tnum)

			rb = New(intcmp)
			for i := 1; i <= tnum; i++ {
				randnum := int(rand.Int31() % int32(tnum))
				rb.Insert(randnum, randnum)
				rb.Insert(randnum*-1, randnum*-1)
				assert.Equal(t, checkHeightBalance(t, rb, rb.root), true)
				assert.Equal(t, checkOrder(t, rb, rb.root), true)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	nums := []int{1, 2, 3, 4, 8, 1024, 1024 + 1}
	for _, num := range nums {
		tnum := num
		t.Run(fmt.Sprintf("[num:%d]", tnum), func(t *testing.T) {
			rb := newInsertNum(tnum)
			rb.Remove(0)
			assert.Equal(t, checkHeightBalance(t, rb, rb.root), true)
			assert.Equal(t, checkOrder(t, rb, rb.root), true)
			assert.Equal(t, rb.size, tnum)
			for i := 1; i <= tnum; i++ {
				rb.Remove(i)
				assert.Equal(t, checkHeightBalance(t, rb, rb.root), true)
				assert.Equal(t, checkOrder(t, rb, rb.root), true)
				assert.Equal(t, rb.size, tnum-i)
			}
			assert.Equal(t, rb.root, _nil)
			assert.Equal(t, rb.size, 0)

			rb = newInsertNum(tnum)
			for i := tnum; i >= 1; i-- {
				rb.Remove(i)
				assert.Equal(t, checkHeightBalance(t, rb, rb.root), true)
				assert.Equal(t, checkOrder(t, rb, rb.root), true)
				assert.Equal(t, rb.size, i-1)
			}
			assert.Equal(t, rb.root, _nil)
			assert.Equal(t, rb.size, 0)

			rb = newInsertNum(tnum)
			left, right := tnum/2, tnum/2+1
			for ; left >= 1 && right <= tnum; left, right = left-1, right+1 {
				rb.Remove(left)
				rb.Remove(right)
				assert.Equal(t, checkHeightBalance(t, rb, rb.root), true)
				assert.Equal(t, checkOrder(t, rb, rb.root), true)
				assert.Equal(t, rb.size, tnum-(right-left+1))
			}
			rb.Remove(left)
			rb.Remove(right)
			assert.Equal(t, rb.root, _nil)
			assert.Equal(t, rb.size, 0)
		})
	}
}

func newInsertNum(num int) *avlTree {
	rb := New(intcmp)
	for i := 1; i <= num; i++ {
		rb.Insert(i, i)
	}
	return rb
}

func TestGet(t *testing.T) {
	num := 1024
	rb := New(intcmp)
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
}

func checkHeightBalance(t *testing.T, rb *avlTree, root *Node) bool {
	balance := checkHeightBalance2(t, root)
	return balance
}

func checkHeightBalance2(t *testing.T, root *Node) bool {
	if root == _nil {
		return true
	}
	diff := root.lchild.h - root.rchild.h
	if diff > 1 || diff < -1 {
		t.Logf("fail node:%v", root)
		return false
	}
	if !checkHeightBalance2(t, root.lchild) {
		return false
	}
	if !checkHeightBalance2(t, root.rchild) {
		return false
	}
	return true
}

func checkOrder(t *testing.T, rb *avlTree, root *Node) bool {
	if root == _nil {
		return rb.size == 0
	}
	var arr []int = make([]int, 0)
	seqence(root, &arr)
	if len(arr) != rb.size {
		t.Logf("fail seq %d %d", rb.size, len(arr))
		return false
	}
	for i := 1; i < rb.size; i++ {
		if rb.cmp(arr[i], arr[i-1]) <= 0 {
			t.Logf("fail seq %v %v", arr[i-1], arr[i])
			return false
		}
	}
	return true
}

func seqence(root *Node, arr *[]int) {
	if root == _nil {
		return
	}
	seqence(root.lchild, arr)
	*arr = append(*arr, root.key.(int))
	seqence(root.rchild, arr)
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
		b.Run(fmt.Sprintf("avltree-%s", name), func(b *testing.B) {
			tr := New(tree.IntComparator)
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
		b.Run(fmt.Sprintf("avltree-%s", name), func(b *testing.B) {
			tr := New(tree.IntComparator)
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
		b.Run(fmt.Sprintf("avltree-%s", name), func(b *testing.B) {
			tr := New(tree.IntComparator)
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
