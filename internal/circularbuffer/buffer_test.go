// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package circularbuffer

import (
	"fmt"
	"testing"
)

func TestPushPop(t *testing.T) {
	testcapacity := []int{1, 2, 3, 4}
	t.Parallel()
	for _, v := range testcapacity {
		capacity := v
		t.Run(fmt.Sprintf("[capacity:%d]", capacity), func(t *testing.T) {
			testPushPopWithCapactiy(t, capacity)
		})
	}
}

func TestResetCapacity(t *testing.T) {
	testcapacity := []int{1, 2, 3, 4}
	for _, v := range testcapacity {
		capacity := v
		t.Run(fmt.Sprintf("[capacity:%d]", v), func(t *testing.T) {
			testResetCapacityWithCapaction(t, capacity)
		})
	}
}

func TestGetSet(t *testing.T) {
	b := New(1)
	b.PushBack(1)
	if b.Get(0) != 1 {
		t.Errorf("expected %d got %d", 1, b.Get(0))
	}
	b.Set(0, 2)
	if b.Get(0) != 2 {
		t.Errorf("expected %d got %d", 2, b.Get(0))
	}
}

func testPushPopWithCapactiy(t *testing.T, capacity int) {
	step := (capacity + 1) / 2
	minxPos, maxPos := -1*(capacity+step), capacity+step
	for i := minxPos; i <= maxPos; i += step {
		block := New(capacity)
		block.start, block.end = i, i
		t.Run(fmt.Sprintf("[position:%d]", i), func(t *testing.T) {
			testPushFrontPopFront(t, block, capacity)
			testPushBackPopBack(t, block, capacity)
			testPushFrontPopBack(t, block, capacity)
			testPushBackPopFront(t, block, capacity)
		})
	}
}

func testPushFrontPopFront(t *testing.T, block *Buffer, capacity int) {
	pushFull := capacity
	for i := 1; i <= pushFull; i++ {
		block.PushFront(i)
	}
	for pos, targetNum := 1, pushFull; pos <= pushFull; pos++ {
		if block.Get(pos-1).(int) != targetNum {
			t.Errorf("expected %d got %d", targetNum, block.Get(pos).(int))
		}
		targetNum--
	}
	if block.Front().(int) != pushFull {
		t.Errorf("expected %d got %d", capacity, block.Front().(int))
	}
	if block.Back().(int) != 1 {
		t.Errorf("expected %d got %d", 1, block.Back().(int))
	}
	if block.Size() != pushFull {
		t.Errorf("expected %d got %d", capacity, block.Size())
	}
	if !block.IsFull() {
		t.Errorf("expected true got false")
	}
	for i := pushFull; i >= 1; i-- {
		v := block.PopFront().(int)
		if v != i {
			t.Errorf("expected %d got %d", i, v)
		}
	}
	if block.Size() != 0 {
		t.Errorf("expected %d got %d", 0, block.Size())
	}
	if !block.IsEmpty() {
		t.Errorf("expected true got false")
	}
	if block.PopBack() != nil || block.PopFront() != nil || block.Back() != nil || block.Front() != nil {
		t.Errorf("expected nil")
	}
}

func testPushBackPopBack(t *testing.T, block *Buffer, capacity int) {
	pushFull := capacity
	for i := 1; i <= pushFull; i++ {
		block.PushBack(i)
	}
	for pos, targetNum := 1, 1; pos <= pushFull; pos++ {
		if block.Get(pos-1).(int) != targetNum {
			t.Errorf("expected %d got %d", targetNum, block.Get(pos).(int))
		}
		targetNum++
	}
	if block.Front().(int) != 1 {
		t.Errorf("expected %d got %d", 1, block.Front().(int))
	}
	if block.Back().(int) != pushFull {
		t.Errorf("expected %d got %d", capacity, block.Back().(int))
	}
	if block.Size() != pushFull {
		t.Errorf("expected %d got %d", capacity, block.Size())
	}
	if !block.IsFull() {
		t.Errorf("expected true got false")
	}
	for i := pushFull; i >= 1; i-- {
		v := block.PopBack().(int)
		if v != i {
			t.Errorf("expected %d got %d", i, v)
		}
		if block.Size() != i-1 {
			t.Errorf("expected %d got %d", i-1, block.Size())
		}
	}
	if !block.IsEmpty() {
		t.Errorf("expected true got false")
	}
	if block.PopBack() != nil || block.PopFront() != nil || block.Front() != nil || block.Back() != nil {
		t.Errorf("expected nil")
	}
}

func testPushFrontPopBack(t *testing.T, block *Buffer, capacity int) {
	for times := 1; times <= 2; times++ {
		pushFull := capacity
		for i := 1; i <= pushFull; i++ {
			block.PushFront(i)
		}
		for i := 1; i <= pushFull; i++ {
			v := block.PopBack().(int)
			if v != i {
				t.Errorf("expected %d got %d", i, v)
			}
		}
		if !block.IsEmpty() {
			t.Errorf("expected true got false")
		}
		if block.PopBack() != nil || block.PopFront() != nil || block.Front() != nil || block.Back() != nil {
			t.Errorf("expected nil")
		}
	}
}

func testPushBackPopFront(t *testing.T, block *Buffer, capacity int) {
	for times := 1; times <= 2; times++ {
		pushFull := capacity
		for i := 1; i <= pushFull; i++ {
			block.PushBack(i)
		}
		for i := 1; i <= pushFull; i++ {
			v := block.PopFront().(int)
			if v != i {
				t.Errorf("expected %d got %d", i, v)
			}
		}
		if !block.IsEmpty() {
			t.Errorf("expected true got false")
		}
		if block.PopBack() != nil || block.PopFront() != nil || block.Front() != nil || block.Back() != nil {
			t.Errorf("expected nil")
		}
	}
}

func testResetCapacityWithCapaction(t *testing.T, capacity int) {
	t.Parallel()
	step := (capacity + 1) / 2
	minxPos, maxPos := -1*(capacity+step), capacity+step
	for pos := minxPos; pos <= maxPos; pos += step {
		p := pos
		t.Run(fmt.Sprintf("postion:%d", p), func(t *testing.T) {
			testExpandCapacityWithPosion(t, capacity, step, p)
			testReduceCapacityWithPosion(t, capacity, step, p)
		})
	}
}

func testExpandCapacityWithPosion(t *testing.T, capacity int, stepLen int, pos int) {
	for pushVal := 1; pushVal <= capacity; pushVal += stepLen {
		block := New(capacity)
		block.start, block.end = pos, pos
		for j := 1; j <= pushVal; j++ {
			block.PushBack(j)
		}
		block.ResetCapacity(capacity + pushVal)
		if block.size != pushVal {
			t.Errorf("expected %d got %d", pushVal, block.size)
		}
		if block.Capacity() != capacity+pushVal {
			t.Errorf("expected %d got %d", pushVal+capacity, block.Capacity())
		}
		for num, idx := 1, 0; num <= pushVal; {
			if block.Get(idx).(int) != num {
				t.Errorf("expected %d got %d", num, block.Get(num-1).(int))
			}
			num++
			idx++
		}
		block.PushFront(0)
		if v := block.PopFront().(int); v != 0 {
			t.Errorf("expected 0 get %d", v)
		}
		block.PushBack(0)
		if v := block.PopBack().(int); v != 0 {
			t.Errorf("expected 0 get %d", v)
		}
		for num, idx := 1, 0; num <= pushVal; {
			if block.Get(idx).(int) != num {
				t.Errorf("expected %d got %d", num, block.Get(idx).(int))
			}
			num++
			idx++
		}
	}
}

func testReduceCapacityWithPosion(t *testing.T, capacity int, stepLen int, pos int) {
	for pushVal := 1; pushVal <= capacity; pushVal += stepLen {
		block := New(capacity)
		block.start, block.end = pos, pos
		for j := 1; j <= pushVal; j++ {
			block.PushBack(j)
		}
		block.ResetCapacity(pushVal + 1)
		if block.size != pushVal {
			t.Errorf("expected %d got %d", pushVal, block.size)
		}
		if block.Capacity() != pushVal+1 {
			t.Errorf("expected %d got %d", pushVal+capacity, block.Capacity())
		}
		if block.IsFull() {
			t.Errorf("expected false got true")
		}
		for num := 1; num <= pushVal; num++ {
			if block.Get(num-1).(int) != num {
				t.Errorf("expected %d got %d", num, block.Get(num-1).(int))
			}
		}
		block.PushFront(0)
		if v := block.PopFront().(int); v != 0 {
			t.Errorf("expected 0 get %d", v)
		}
		block.PushBack(0)
		if v := block.PopBack().(int); v != 0 {
			t.Errorf("expected 0 get %d", v)
		}
		block.ResetCapacity(pushVal)
		for num, idx := 1, 0; num <= pushVal; {
			if block.Get(idx).(int) != num {
				t.Errorf("expected %d got %d", num, block.Get(idx).(int))
			}
			num++
			idx++
		}
	}
}
