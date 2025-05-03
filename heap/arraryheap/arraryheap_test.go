// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package arraryheap

import (
	"fmt"
	"testing"

	"github.com/mrtcx/plusdata/heap"
	"github.com/mrtcx/plusdata/internal/assert"
)

var testLess = func(a, b interface{}) bool { return a.(int) < b.(int) }

func TestPushPop(t *testing.T) {
	_blockCap := 1024 * 64
	testNum := []int{1, 2, 3, 4, 2*_blockCap + 2}
	for _, v := range testNum {
		num := v
		t.Run(fmt.Sprintf("[num:%d]", num), func(t *testing.T) {
			q := New(func(a, b interface{}) bool { return a.(int) < b.(int) })
			testPushPop(t, q, num)
		})
	}
}

func TestClean(t *testing.T) {
	_blockCap := 1024 * 64
	testNums := []int{1, 2, 3, 4, 2*_blockCap + 2}
	for _, v := range testNums {
		num := v
		t.Run(fmt.Sprintf("[num:%d]", num), func(t *testing.T) {
			q := New(testLess)
			for i := 0; i < num; i++ {
				q.Push(i)
			}
			q.Clean()
			assert.Equal(t, q.Size(), 0)
			assert.Equal(t, q.Empty(), true)
			assert.Equal(t, q.Top(), nil)
			assert.Equal(t, q.Pop(), nil)
		})
	}
}

func testPushPop(t *testing.T, q *ArraryHeap, num int) {
	for k := 0; k < 3; k++ {
		pushVal := num
		switch k {
		case 0:
			for i := 1; i <= pushVal; i++ {
				q.Push(i)
			}
		case 1:
			for i := pushVal; i >= 1; i-- {
				q.Push(i)
			}
		case 2:
			for i, j := 1, pushVal; i <= j; {
				q.Push(i)
				if j != i {
					q.Push(j)
				}
				i++
				j--
			}
		}

		assert.Equal(t, q.Top(), 1)
		assert.Equal(t, q.Size(), pushVal)
		for i := 1; i <= pushVal; i++ {
			assert.Equal(t, q.Top(), i)
			assert.Equal(t, q.Pop(), i)
			assert.Equal(t, q.Size(), pushVal-i)
		}
		assert.Equal(t, q.Top(), nil)
		assert.Equal(t, q.Empty(), true)
		assert.Equal(t, q.Size(), 0)
	}
}

func BenchmarkPush(b *testing.B) {
	testsizes := []int{100000, 500000, 1000000, 2000000}
	testNames := []string{"10w", "50w", "100w", "200w"}
	for i, v := range testsizes {
		size := v
		name := testNames[i]
		b.Run(name, func(b *testing.B) {
			h := New(heap.IntLess)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < size; j++ {
					h.Push(j)
				}
			}
		})
	}
}
