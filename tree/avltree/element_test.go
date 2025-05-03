// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package avltree

import (
	"fmt"
	"testing"

	"github.com/mrtcx/plusdata/internal/assert"
	"github.com/mrtcx/plusdata/tree"
)

func TestFind(t *testing.T) {
	nums := []int{1, 2, 4, 8, 9, 1024, 1024 + 1}
	for _, num := range nums {
		tnum := num

		t.Run(fmt.Sprintf("[num:%d]", tnum), func(t *testing.T) {
			av := New(intcmp)
			for i := 1; i <= tnum; i++ {
				av.Insert(2*i, 2*i)
			}
			for i := 1; i <= tnum; i++ {
				e := av.Find(2 * i)
				assert.Equal(t, e.Key(), 2*i)
				assert.Equal(t, e.Value(), 2*i)
				assert.Equal(t, av.Find(2*i-1), nil)
			}
			for i := 1; i <= tnum; i++ {
				av.Remove(2 * i)
				e := av.Find(2 * i)
				assert.Equal(t, e, nil)
			}
		})
	}
}

func TestLeft(t *testing.T) {
	nums := []int{1, 2, 4, 8, 9, 1024, 1024 + 1}
	for _, num := range nums {
		tnum := num
		t.Run(fmt.Sprintf("[num:%d]", tnum), func(t *testing.T) {
			av := New(intcmp)
			assert.Equal(t, av.Left(), nil)
			for i := 1; i <= tnum; i++ {
				av.Insert(i, i)
			}
			left := av.Left()
			for i := 1; i <= tnum; i++ {
				assert.Equal(t, left.Key(), i)
				assert.Equal(t, left.Value(), i)
				left = left.Next()
			}
		})
	}
}

func TestRight(t *testing.T) {
	nums := []int{1, 2, 4, 8, 9, 1024, 1024 + 1}
	for _, num := range nums {
		tnum := num
		t.Run(fmt.Sprintf("[num:%d]", tnum), func(t *testing.T) {
			av := New(intcmp)
			assert.Equal(t, av.Right(), nil)
			for i := 1; i <= tnum; i++ {
				av.Insert(i, i)
			}
			right := av.Right()
			for i := tnum; i >= 1; i-- {
				assert.Equal(t, right.Key(), i)
				assert.Equal(t, right.Value(), i)
				right = right.Prev()
			}
		})
	}
}

func TestPrev(t *testing.T) {
	nums := []int{1, 2, 4, 8, 9, 1024, 1024 + 1}
	for _, num := range nums {
		tnum := num
		t.Run(fmt.Sprintf("[num:%d]", tnum), func(t *testing.T) {
			av := New(intcmp)
			assert.Equal(t, av.Prev(2), nil)
			for i := 1; i <= tnum; i++ {
				av.Insert(2*i, 2*i)
			}
			assert.Equal(t, av.Prev(0), nil)
			assert.Equal(t, av.Prev(2*1), nil)
			for i := 2; i <= tnum; i++ {
				e := av.Prev(2 * i)
				assert.Equal(t, e.Key(), (i-1)*2)
				assert.Equal(t, e.Value(), (i-1)*2)

				e = av.Prev(2*i + 1)
				assert.Equal(t, e.Key(), i*2)
				assert.Equal(t, e.Value(), i*2)
			}

			l, r := (tnum+1)/2, (tnum+1)/2+1
			for ; r <= tnum; r++ {
				av.Remove(2 * r)
				e := av.Prev(2 * r)
				assert.Equal(t, e.Key(), l*2)
				assert.Equal(t, e.Value(), l*2)

				e = av.Prev(2 * (r + 1))
				assert.Equal(t, e.Key(), l*2)
				assert.Equal(t, e.Value(), l*2)
			}
			for ; l > 1; l-- {
				av.Remove(2 * l)
				e := av.Prev(2 * l)
				assert.Equal(t, e.Key(), (l-1)*2)
				assert.Equal(t, e.Value(), (l-1)*2)

				e = av.Prev(2 * (l + 1))
				assert.Equal(t, e.Key(), (l-1)*2)
				assert.Equal(t, e.Value(), (l-1)*2)
			}
			av.Remove(2 * l)
			assert.Equal(t, av.Prev(r), nil)
			assert.Equal(t, av.Prev(l), nil)
		})
	}
}

func TestNext(t *testing.T) {
	nums := []int{1, 2, 4, 8, 9, 1024, 1024 + 1}
	for _, num := range nums {
		tnum := num
		t.Run(fmt.Sprintf("[num:%d]", tnum), func(t *testing.T) {
			av := New(intcmp)
			assert.Equal(t, av.Next(0), nil)
			for i := 1; i <= tnum; i++ {
				av.Insert(2*i, 2*i)
			}
			assert.Equal(t, av.Next(2*tnum), nil)
			assert.Equal(t, av.Next(2*tnum+1), nil)
			for i := 1; i <= tnum-1; i++ {
				e := av.Next(2 * i)
				assert.Equal(t, e.Key(), (i+1)*2)
				assert.Equal(t, e.Value(), (i+1)*2)

				e = av.Next(2*i + 1)
				assert.Equal(t, e.Key(), (i+1)*2)
				assert.Equal(t, e.Value(), (i+1)*2)
			}

			l, r := (tnum+1)/2, (tnum+1)/2+1
			for ; l >= 1 && r <= tnum; l-- {
				av.Remove(2 * l)
				e := av.Next(2 * l)
				assert.Equal(t, e.Key(), r*2)
				assert.Equal(t, e.Value(), r*2)

				e = av.Next((l - 1) * 2)
				assert.Equal(t, e.Key(), r*2)
				assert.Equal(t, e.Value(), r*2)
			}
			for ; r < tnum; r++ {
				av.Remove(2 * r)
				e := av.Next(2 * r)
				assert.Equal(t, e.Key(), (r+1)*2)
				assert.Equal(t, e.Value(), (r+1)*2)
			}
			av.Remove(2 * r)
			assert.Equal(t, av.Prev(r), nil)
			assert.Equal(t, av.Prev(l), nil)
		})
	}
}

func TestElementPrev(t *testing.T) {
	nums := []int{1, 2, 4, 8, 9, 1024, 1024 + 1}
	for _, num := range nums {
		tnum := num
		t.Run(fmt.Sprintf("[num:%d]", tnum), func(t *testing.T) {
			av := New(intcmp)
			for i := 1; i <= tnum; i++ {
				av.Insert(i, i)
			}
			for i := 1; i <= tnum; i++ {
				e := av.Find(i)
				for j := i; j >= 1; j-- {
					assert.Equal(t, e.Key(), j)
					assert.Equal(t, e.Value(), j)
					e = e.Prev()
				}
				assert.Equal(t, e, nil)
			}
		})
	}
}

func TestElementNext(t *testing.T) {
	nums := []int{1, 2, 4, 8, 9, 1024, 1024 + 1}
	for _, num := range nums {
		tnum := num
		t.Run(fmt.Sprintf("[num:%d]", tnum), func(t *testing.T) {
			av := New(intcmp)
			for i := 1; i <= tnum; i++ {
				av.Insert(i, i)
			}
			for i := 1; i <= tnum; i++ {
				e := av.Find(i)
				for j := i; j <= tnum; j++ {
					assert.Equal(t, e.Key(), j)
					assert.Equal(t, e.Value(), j)
					e = e.Next()
				}
				assert.Equal(t, e, nil)
			}
		})
	}
}
func TestElementSet(t *testing.T) {
	av := New(intcmp)
	av.Insert(1, 1)
	e := av.Find(1)
	e.SetValue(2)
	assert.Equal(t, e.Value(), 2)
	geti, _ := av.Get(1)
	assert.Equal(t, geti, 2)
}

func BenchmarkPrev(b *testing.B) {
	testsizes := []int{100000, 500000, 1000000, 2000000}
	testNames := []string{"10w", "50w", "100w", "200w"}
	for i, v := range testsizes {
		size := v
		name := testNames[i]
		b.Run(name, func(b *testing.B) {
			tr := New(tree.IntComparator)
			for j := 0; j < size; j++ {
				tr.Insert(j, nil)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < size; j++ {
					_ = tr.Prev(j)
				}
			}
		})
	}
}

func BenchmarkNext(b *testing.B) {
	testsizes := []int{100000, 500000, 1000000, 2000000}
	testNames := []string{"10w", "50w", "100w", "200w"}
	for i, v := range testsizes {
		size := v
		name := testNames[i]
		b.Run(name, func(b *testing.B) {
			tr := New(tree.IntComparator)
			for j := 0; j < size; j++ {
				tr.Insert(j, nil)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < size; j++ {
					_ = tr.Next(j)
				}
			}
		})
	}
}

func BenchmarkElementPrev(b *testing.B) {
	testsizes := []int{100000, 500000, 1000000, 2000000}
	testNames := []string{"10w", "50w", "100w", "200w"}
	for i, v := range testsizes {
		size := v
		name := testNames[i]
		b.Run(name, func(b *testing.B) {
			tr := New(tree.IntComparator)
			for j := 0; j < size; j++ {
				tr.Insert(j, nil)
			}
			right := tr.Right()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r := right
				for j := 0; j < size; j++ {
					r = r.Prev()
				}
			}
		})
	}
}

func BenchmarkElementNext(b *testing.B) {
	testsizes := []int{100000, 500000, 1000000, 2000000}
	testNames := []string{"10w", "50w", "100w", "200w"}
	for i, v := range testsizes {
		size := v
		name := testNames[i]
		b.Run(name, func(b *testing.B) {
			tr := New(tree.IntComparator)
			for j := 0; j < size; j++ {
				tr.Insert(j, nil)
			}
			left := tr.Left()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r := left
				for j := 0; j < size; j++ {
					r = r.Next()
				}
			}
		})
	}
}
