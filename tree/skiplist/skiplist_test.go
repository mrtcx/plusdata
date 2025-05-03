// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package skiplist

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
			sk := New(intcmp)
			assert.Equal(t, sk.head.nexts[0], nil)
			for i := 1; i <= tnum; i++ {
				sk.Insert(i, i)
				sk.Insert(i, i)
				assert.Equal(t, checkOrder(t, sk), true)
			}
			assert.Equal(t, sk.size, tnum)

			sk = New(intcmp)
			for i := tnum; i >= 1; i-- {
				sk.Insert(i, i)
				sk.Insert(i, i)
				assert.Equal(t, checkOrder(t, sk), true)
			}
			assert.Equal(t, sk.size, tnum)

			sk = New(intcmp)
			for i, j := 1, tnum; i <= j; i, j = i+1, j-1 {
				sk.Insert(i, i)
				sk.Insert(j, j)
				assert.Equal(t, checkOrder(t, sk), true)
			}
			assert.Equal(t, sk.size, tnum)

			sk = New(intcmp)
			for i := 1; i <= tnum; i++ {
				randnum := int(rand.Int31() % int32(tnum))
				sk.Insert(randnum, randnum)
				sk.Insert(randnum*-1, randnum*-1)
				assert.Equal(t, checkOrder(t, sk), true)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	nums := []int{1, 2, 3, 4, 8, 1024, 1024 + 1}
	for _, num := range nums {
		tnum := num
		t.Run(fmt.Sprintf("[num:%d]", tnum), func(t *testing.T) {
			sk := newInsertNum(tnum)
			sk.Remove(0)
			assert.Equal(t, checkOrder(t, sk), true)
			assert.Equal(t, sk.size, tnum)
			for i := 1; i <= tnum; i++ {
				sk.Remove(i)
				assert.Equal(t, checkOrder(t, sk), true)
				assert.Equal(t, sk.size, tnum-i)
			}
			assert.Equal(t, sk.head.nexts[0], nil)
			assert.Equal(t, sk.size, 0)

			sk = newInsertNum(tnum)
			for i := tnum; i >= 1; i-- {
				sk.Remove(i)
				assert.Equal(t, checkOrder(t, sk), true)
				assert.Equal(t, sk.size, i-1)
			}
			assert.Equal(t, sk.head.nexts[0], nil)
			assert.Equal(t, sk.size, 0)

			sk = newInsertNum(tnum)
			left, right := tnum/2, tnum/2+1
			for ; left >= 1 && right <= tnum; left, right = left-1, right+1 {
				sk.Remove(left)
				sk.Remove(right)
				assert.Equal(t, checkOrder(t, sk), true)
				assert.Equal(t, sk.size, tnum-(right-left+1))
			}
			sk.Remove(left)
			sk.Remove(right)
			assert.Equal(t, sk.head.nexts[0], nil)
			assert.Equal(t, sk.size, 0)
		})
	}
}

func newInsertNum(num int) *skipList {
	sk := New(intcmp)
	for i := 1; i <= num; i++ {
		sk.Insert(i, i)
	}
	return sk
}

func TestGet(t *testing.T) {
	num := 1024
	sk := New(intcmp)
	geti, ei := sk.Get(1)
	assert.Equal(t, geti, nil)
	assert.Equal(t, ei, false)
	for i := 1; i <= num; i++ {
		sk.Insert(i, i)
		geti, ei = sk.Get(i)
		assert.Equal(t, geti, i)
		assert.Equal(t, ei, true)
		geti, ei = sk.Get(i + 1)
		assert.Equal(t, geti, nil)
		assert.Equal(t, ei, false)
	}
	for i := 1; i <= num; i++ {
		sk.Remove(i)
		_, ei = sk.Get(i)
		assert.Equal(t, ei, false)
		for j := i + 1; j <= num; j++ {
			geti, ei = sk.Get(j)
			assert.Equal(t, ei, true)
			assert.Equal(t, geti, j)
		}
	}
}

func checkOrder(t *testing.T, sk *skipList) bool {
	if sk.head.nexts[0] == nil {
		return sk.size == 0
	}
	var arr []interface{}
	head := sk.head.nexts[0]
	for head != nil {
		arr = append(arr, head.key)
		if sk.cmp(head.key, head.value) != 0 {
			t.Logf("fail node:%v", head)
			return false
		}
		head = head.nexts[0]
	}
	if len(arr) != sk.size {
		t.Logf("fail size %d %d", len(arr), sk.size)
		return false
	}
	for i := 1; i < sk.size; i++ {
		if sk.cmp(arr[i], arr[i-1]) <= 0 {
			t.Logf("fail order %d %d", arr[i], arr[i-1])
			return false
		}
	}
	return true
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
		b.Run(fmt.Sprintf("skiplist-%s", name), func(b *testing.B) {
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
		b.Run(fmt.Sprintf("skiplist-%s", name), func(b *testing.B) {
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
		b.Run(fmt.Sprintf("skiplist-%s", name), func(b *testing.B) {
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
