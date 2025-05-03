// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package circularblocks

import (
	"fmt"
	"testing"

	"github.com/mrtcx/plusdata/internal/assert"
	"github.com/mrtcx/plusdata/internal/circularbuffer"
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
			for i := 0; i < testnum; i++ {
				q.PushFront(i)
			}
			assert.Equal(t, q.Front(), testnum-1)
			for i := 0; i < testnum; i++ {
				q.PopFront()
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
			assert.Equal(t, q.Front(), nil)
			for i := 0; i < testnum; i++ {
				q.PushBack(i)
			}
			assert.Equal(t, q.Back(), testnum-1)
			for i := 0; i < testnum; i++ {
				q.PopBack()
			}
			assert.Equal(t, q.Back(), nil)
			for i := 0; i < testnum; i++ {
				q.PushFront(i)
			}
			assert.Equal(t, q.Back(), 0)
			for i := 0; i < testnum; i++ {
				q.PopFront()
			}
			assert.Equal(t, q.Back(), nil)
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
			for i := 0; i < testnum; i++ {
				q.PushFront(i)
			}
			for i := 0; i < testnum; i++ {
				assert.Equal(t, q.Get(i), testnum-1-i)
			}
			q.PopFront()
			for i := 0; i < testnum-1; i++ {
				assert.Equal(t, q.Get(i), testnum-2-i)
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
			q.Clean()
			for i := 0; i < testnum; i++ {
				q.PushFront(i)
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

func TestPushPop(t *testing.T) {
	testNum := []int{1, 2, 3, 4, _block, _block + 2, 2 * _block, 2*_block + 2}
	for _, v := range testNum {
		num := v
		t.Run(fmt.Sprintf("[num:%d]", num), func(t *testing.T) {
			q := New()
			testPushFrontPopFront(t, q, num)
			testPushBackPopBack(t, q, num)
			testPushFrontPopBack(t, q, num)
			testPushBackPopFront(t, q, num)
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
	assert.Equal(t, q.PopFront(), nil)
}

func testPushFrontPopFront(t *testing.T, q *circularBlocks, num int) {
	pushFull := num
	for i := 1; i <= pushFull; i++ {
		q.PushFront(i)
		assert.Equal(t, checkBlockCap(t, q), true)
	}
	for pos, targetNum := 1, pushFull; pos <= pushFull; pos++ {
		assert.Equal(t, q.Get(pos-1), targetNum)
		targetNum--
	}
	assert.Equal(t, q.Front(), pushFull)
	assert.Equal(t, q.Back(), 1)
	assert.Equal(t, q.Size(), pushFull)
	for i := pushFull; i >= 1; i-- {
		assert.Equal(t, q.PopFront(), i)
		assert.Equal(t, checkBlockCap(t, q), true)
	}
	assert.Equal(t, q.Size(), 0)
	if q.PopBack() != nil || q.PopFront() != nil || q.Back() != nil || q.Front() != nil {
		t.Errorf("expected nil")
	}
}

func testPushBackPopBack(t *testing.T, q *circularBlocks, num int) {
	pushFull := num
	for i := 1; i <= pushFull; i++ {
		q.PushBack(i)
		assert.Equal(t, checkBlockCap(t, q), true)
	}
	for pos, targetNum := 1, 1; pos <= pushFull; pos++ {
		assert.Equal(t, q.Get(pos-1), targetNum)
		targetNum++
	}
	assert.Equal(t, q.Front(), 1)
	assert.Equal(t, q.Back(), pushFull)
	assert.Equal(t, q.Size(), pushFull)
	for i := pushFull; i >= 1; i-- {
		assert.Equal(t, q.PopBack(), i)
		assert.Equal(t, q.Size(), i-1)
		assert.Equal(t, checkBlockCap(t, q), true)
	}
	if q.PopBack() != nil || q.PopFront() != nil || q.Front() != nil || q.Back() != nil {
		t.Errorf("expected nil")
	}
}

func testPushFrontPopBack(t *testing.T, q *circularBlocks, num int) {
	for times := 1; times <= 2; times++ {
		pushFull := num
		for i := 1; i <= pushFull; i++ {
			q.PushFront(i)
			assert.Equal(t, checkBlockCap(t, q), true)
		}
		for i := 1; i <= pushFull; i++ {
			assert.Equal(t, q.PopBack(), i)
			assert.Equal(t, checkBlockCap(t, q), true)
		}
		if q.PopBack() != nil || q.PopFront() != nil || q.Front() != nil || q.Back() != nil {
			t.Errorf("expected nil")
		}
	}
}

func testPushBackPopFront(t *testing.T, q *circularBlocks, num int) {
	for times := 1; times <= 2; times++ {
		pushFull := num
		for i := 1; i <= pushFull; i++ {
			q.PushBack(i)
			assert.Equal(t, checkBlockCap(t, q), true)
		}
		for i := 1; i <= pushFull; i++ {
			assert.Equal(t, q.PopFront(), i)
			assert.Equal(t, checkBlockCap(t, q), true)
		}
		if q.PopBack() != nil || q.PopFront() != nil || q.Front() != nil || q.Back() != nil {
			t.Errorf("expected nil")
		}
	}
}

func checkBlockCap(t *testing.T, b *circularBlocks) bool {
	if b.bbs.Capacity() > 1 && b.bbs.Capacity()/4 > b.bbs.Size() {
		return false
	}
	for i := 0; i < b.bbs.Size(); i++ {
		bu := b.bbs.Get(i).(*circularbuffer.Buffer)
		if bu.Capacity() > _initBlockCap && bu.Capacity()/4 > bu.Size() {
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
		b.Run(fmt.Sprintf("circularblocks[%s]", name), func(b *testing.B) {
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

func BenchmarkPushFront(b *testing.B) {
	testsizes := []int{100000, 500000, 1000000, 2000000}
	testNames := []string{"10w", "50w", "100w", "200w"}
	for i, v := range testsizes {
		size := v
		name := testNames[i]

		b.Run(fmt.Sprintf("circularblocks[%s]", name), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				arr := New()
				for j := 0; j < size; j++ {
					arr.PushFront(j)
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

func BenchmarkPopFront(b *testing.B) {
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
					arr.PopFront()
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
		b.Run(fmt.Sprintf("circularblocks[%s]", name), func(b *testing.B) {
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
