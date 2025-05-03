// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package blockslices

import (
	"fmt"
	"testing"

	"github.com/mrtcx/plusdata/internal/assert"
)

func TestFront(t *testing.T) {
	testNums := []int{1, 2, 3, 4, _block, _block + 1, 2 * _block, 2*_block + 1}
	for _, v := range testNums {
		testnum := v
		t.Run(fmt.Sprintf("[num:%d]", testnum), func(t *testing.T) {
			q := New()
			assert.Equal(t, q.Front(), nil)
			for i := 0; i < testnum; i++ {
				q.PushBack(i)
			}
			assert.Equal(t, q.Front(), 0)
			for i := 0; i < testnum; i++ {
				q.PopBack()
			}
			assert.Equal(t, q.Front(), nil)
		})
	}
}

func TestBack(t *testing.T) {
	testNums := []int{1, 2, 3, 4, _block, _block + 1, 2 * _block, 2*_block + 1}
	for _, v := range testNums {
		testnum := v
		t.Run(fmt.Sprintf("[num:%d]", testnum), func(t *testing.T) {
			q := New()
			assert.Equal(t, q.Back(), nil)
			for i := 0; i < testnum; i++ {
				q.PushBack(i)
			}
			assert.Equal(t, q.Back(), testnum-1)
			for i := 0; i < testnum; i++ {
				q.PopBack()
			}
			assert.Equal(t, q.Back(), nil)
		})
	}
}

func TestPushBack(t *testing.T) {
	testNums := []int{1, 2, 3, 4, _block, _block + 1, 2 * _block, 2*_block + 1}
	for _, v := range testNums {
		testnum := v
		t.Run(fmt.Sprintf("[num:%d]", testnum), func(t *testing.T) {
			q := New()
			for i := 0; i < testnum; i++ {
				q.PushBack(i)
				assert.Equal(t, q.Back(), i)
				assert.Equal(t, q.Front(), 0)
				assert.Equal(t, q.size, i+1)
				assert.Equal(t, checkBlockCap(t, q), true)
			}
		})
	}
}

func TestPopBack(t *testing.T) {
	testNums := []int{1, 2, 3, 4, _block, _block + 1, 2 * _block, 2*_block + 1}
	for _, v := range testNums {
		testnum := v
		t.Run(fmt.Sprintf("[num:%d]", testnum), func(t *testing.T) {
			q := New()
			for i := 0; i < testnum; i++ {
				q.PushBack(i)
			}
			for i := testnum - 1; i >= 0; i-- {
				assert.Equal(t, q.Front(), 0)
				assert.Equal(t, q.Back(), i)
				assert.Equal(t, q.Size(), i+1)
				q.PopBack()
				assert.Equal(t, checkBlockCap(t, q), true)
			}
			assert.Equal(t, q.Front(), nil)
			assert.Equal(t, q.Back(), nil)
			assert.Equal(t, q.Size(), 0)
		})
	}
}

func TestGet(t *testing.T) {
	testNums := []int{1, 2, 3, 4, _block, _block + 1, 2 * _block, 2*_block + 1}
	for _, v := range testNums {
		testnum := v
		t.Run(fmt.Sprintf("[num:%d]", testnum), func(t *testing.T) {
			q := New()
			for i := 0; i < testnum; i++ {
				q.PushBack(i)
			}
			for i := 0; i < testnum; i++ {
				assert.Equal(t, q.Get(i), i)
			}
			q.PopBack()
			for i := 0; i < testnum-1; i++ {
				assert.Equal(t, q.Get(i), i)
			}
		})
	}
}

func TestSet(t *testing.T) {
	testNums := []int{1, 2, 3, 4, _block, _block + 1, 2 * _block, 2*_block + 1}
	for _, v := range testNums {
		testnum := v
		t.Run(fmt.Sprintf("[num:%d]", testnum), func(t *testing.T) {
			q := New()
			for i := 0; i < testnum; i++ {
				q.PushBack(i)
			}
			for i := 0; i < testnum; i++ {
				q.Set(i, i*-1)
			}
			for i := 0; i < testnum; i++ {
				assert.Equal(t, q.Get(i), i*-1)
			}
		})
	}
}

func TestClean(t *testing.T) {
	q := New()
	q.PushBack(1)
	q.Clean()
	assert.Equal(t, q.Size(), 0)
	assert.Equal(t, q.Empty(), true)
	assert.Equal(t, q.Front(), nil)
	assert.Equal(t, q.Back(), nil)
	assert.Equal(t, q.PopBack(), nil)
	assert.Equal(t, q.blocks, nil)
}

func checkBlockCap(t *testing.T, arr *blockSlice) bool {
	if cap(arr.blocks)/4 > len(arr.blocks) {
		t.Logf("fail:%d %d", len(arr.blocks)*4, cap(arr.blocks))
		return false
	}
	for _, b := range arr.blocks {
		_ = b
		if cap(b)/4 > len(b) {
			t.Logf("fail:%d %d", len(b)*4, cap(b))
			return false
		}
	}
	return true
}

func BenchmarkPushBack(b *testing.B) {
	testsizes := []int{100000, 500000, 1000000, 2000000}
	testNames := []string{"10w", "50w", "100w", "200w"}
	for i, v := range testsizes {
		size := v
		name := testNames[i]
		b.Run(fmt.Sprintf("slice[%s]", name), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var arr []interface{}
				for j := 0; j < size; j++ {
					arr = append(arr, i)
				}
			}
		})
		b.Run(fmt.Sprintf("blockslice[%s]", name), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				arr := New()
				for j := 0; j < size; j++ {
					arr.PushBack(j)
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
		b.Run(fmt.Sprintf("slice[%s]", name), func(b *testing.B) {
			var arr []interface{}
			for i := 0; i < size; i++ {
				arr = append(arr, i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < size; j++ {
					_ = arr[j]
				}
			}
		})
		b.Run(fmt.Sprintf("blockslice[%s]", name), func(b *testing.B) {
			var arr = New()
			for i := 0; i < size; i++ {
				arr.PushBack(i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < size; j++ {
					_ = arr.Get(j)
				}
			}
		})
	}
}

func BenchmarkPopBack(b *testing.B) {
	testsizes := []int{100000, 500000, 1000000, 2000000}
	testNames := []string{"10w", "50w", "100w", "200w"}
	for i, v := range testsizes {
		size := v
		name := testNames[i]
		b.Run(fmt.Sprintf("[%s]", name), func(b *testing.B) {
			arr := New()
			for i := 1; i <= size; i++ {
				arr.PushBack(i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for k := 0; k < size; k++ {
					arr.PopBack()
				}
			}
		})
	}
}
