// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package doublelinkedlist

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	l := New()
	if l.head.next != l.head || l.head.prev != l.head {
		t.Errorf("Expected empty list")
	}
}

func nexti(e *Element, i int) *Element {
	for i > 0 {
		e = e.next
		i--
	}
	return e
}

func toStringByNext(l *List) string {
	s := ""
	e := l.head
	for {
		s += fmt.Sprintf("%d->", e.Value)
		e = e.next
		if e == l.head {
			s += fmt.Sprintf("%d", e.Value)
			break
		}
	}
	return s
}

func toStringByPrev(l *List) string {
	s := ""
	e := l.head
	for {
		s += fmt.Sprintf("%d->", e.Value)
		e = e.prev
		if e == l.head {
			s += fmt.Sprintf("%d", e.Value)
			break
		}
	}
	return s
}

func TestBack(t *testing.T) {
	l := New()
	if l.Back() != nil {
		t.Errorf("Expected %v, got %v", nil, l.Back())
	}
	l.PushBack(1)
	if l.Back() != l.head.prev || l.Back().Value != 1 {
		t.Errorf("Expected %d, got %d", 2, l.Back().Value)
	}
}

func TestFront(t *testing.T) {
	l := New()
	if l.Front() != nil {
		t.Errorf("Expected %v, got %v", nil, l.Front())
	}
	l.PushFront(1)
	if l.Front() != l.head.next || l.Front().Value != 1 {
		t.Errorf("Expected %d, got %d", 1, l.Front().Value)
	}
}

func TestInsertAfter(t *testing.T) {
	l := New()
	l.head.Value = 0
	testCase := []struct {
		insertValue        int
		expectStringOfNext string
		expectStringOfPrev string
		expectSize         int
	}{
		{
			insertValue:        1,
			expectStringOfNext: "0->1->0",
			expectStringOfPrev: "0->1->0",
			expectSize:         1,
		},
		{
			insertValue:        2,
			expectStringOfNext: "0->2->1->0",
			expectStringOfPrev: "0->1->2->0",
			expectSize:         2,
		},
		{
			insertValue:        3,
			expectStringOfNext: "0->3->2->1->0",
			expectStringOfPrev: "0->1->2->3->0",
			expectSize:         3,
		},
	}
	for i, test := range testCase {
		l.InsertAfter(l.head, test.insertValue)
		if toStringByNext(l) != test.expectStringOfNext {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfNext, toStringByNext(l))
		}
		if toStringByPrev(l) != test.expectStringOfPrev {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfPrev, toStringByPrev(l))
		}
		if l.size != test.expectSize {
			t.Errorf("%d Expected %d, got %d", i, test.expectSize, l.size)
		}
	}
}

func TestInsertBefore(t *testing.T) {
	l := New()
	l.head.Value = 0
	testCase := []struct {
		insertValue        int
		expectStringOfNext string
		expectStringOfPrev string
		expectSize         int
	}{
		{
			insertValue:        1,
			expectStringOfNext: "0->1->0",
			expectStringOfPrev: "0->1->0",
			expectSize:         1,
		},
		{
			insertValue:        2,
			expectStringOfNext: "0->1->2->0",
			expectStringOfPrev: "0->2->1->0",
			expectSize:         2,
		},
		{
			insertValue:        3,
			expectStringOfNext: "0->1->2->3->0",
			expectStringOfPrev: "0->3->2->1->0",
			expectSize:         3,
		},
	}
	for i, test := range testCase {
		l.InsertBefore(l.head, test.insertValue)
		if toStringByNext(l) != test.expectStringOfNext {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfNext, toStringByNext(l))
		}
		if toStringByPrev(l) != test.expectStringOfPrev {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfPrev, toStringByPrev(l))
		}
		if l.size != test.expectSize {
			t.Errorf("%d Expected %d, got %d", i, test.expectSize, l.size)
		}
	}
}

func TestPushFront(t *testing.T) {
	l := New()
	l.head.Value = 0
	testCase := []struct {
		pushValue          int
		expectStringOfNext string
		expectStringOfPrev string
		expectSize         int
	}{
		{
			pushValue:          1,
			expectStringOfNext: "0->1->0",
			expectStringOfPrev: "0->1->0",
			expectSize:         1,
		},
		{
			pushValue:          2,
			expectStringOfNext: "0->2->1->0",
			expectStringOfPrev: "0->1->2->0",
			expectSize:         2,
		},
		{
			pushValue:          3,
			expectStringOfNext: "0->3->2->1->0",
			expectStringOfPrev: "0->1->2->3->0",
			expectSize:         3,
		},
	}
	for i, test := range testCase {
		l.PushFront(test.pushValue)
		if toStringByNext(l) != test.expectStringOfNext {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfNext, toStringByNext(l))
		}
		if toStringByPrev(l) != test.expectStringOfPrev {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfPrev, toStringByPrev(l))
		}
		if l.size != test.expectSize {
			t.Errorf("%d Expected %d, got %d", i, test.expectSize, l.size)
		}
	}
}

func TestPushBack(t *testing.T) {
	l := New()
	l.head.Value = 0
	testCase := []struct {
		pushValue          int
		expectStringOfNext string
		expectStringOfPrev string
		expectSize         int
	}{
		{
			pushValue:          1,
			expectStringOfNext: "0->1->0",
			expectStringOfPrev: "0->1->0",
			expectSize:         1,
		},
		{
			pushValue:          2,
			expectStringOfNext: "0->1->2->0",
			expectStringOfPrev: "0->2->1->0",
			expectSize:         2,
		},
		{
			pushValue:          3,
			expectStringOfNext: "0->1->2->3->0",
			expectStringOfPrev: "0->3->2->1->0",
			expectSize:         3,
		},
	}
	for i, test := range testCase {
		l.PushBack(test.pushValue)
		if toStringByNext(l) != test.expectStringOfNext {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfNext, toStringByNext(l))
		}
		if toStringByPrev(l) != test.expectStringOfPrev {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfPrev, toStringByPrev(l))
		}
		if l.size != test.expectSize {
			t.Errorf("%d Expected %d, got %d", i, test.expectSize, l.size)
		}
	}
}

func TestMoveAfter(t *testing.T) {
	l := New()
	l.head.Value = 0
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	testcase := []struct {
		baseIdx            int
		moveIdx            int
		expectStringOfNext string
		expectStringOfPrev string
	}{
		{
			baseIdx:            0,
			moveIdx:            0,
			expectStringOfNext: "0->1->2->3->0",
			expectStringOfPrev: "0->3->2->1->0",
		},
		{
			baseIdx:            0,
			moveIdx:            1,
			expectStringOfNext: "0->1->2->3->0",
			expectStringOfPrev: "0->3->2->1->0",
		},
		{
			baseIdx:            2,
			moveIdx:            1,
			expectStringOfNext: "0->2->1->3->0",
			expectStringOfPrev: "0->3->1->2->0",
		},
		{
			baseIdx:            2,
			moveIdx:            1,
			expectStringOfNext: "0->1->2->3->0",
			expectStringOfPrev: "0->3->2->1->0",
		},
		{
			baseIdx:            1,
			moveIdx:            3,
			expectStringOfNext: "0->1->3->2->0",
			expectStringOfPrev: "0->2->3->1->0",
		},
	}
	for i, test := range testcase {
		l.MoveAfter(nexti(l.head, test.baseIdx), nexti(l.head, test.moveIdx))
		if toStringByNext(l) != test.expectStringOfNext {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfNext, toStringByNext(l))
		}
		if toStringByPrev(l) != test.expectStringOfPrev {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfPrev, toStringByPrev(l))
		}
	}
}

func TestMoveBefore(t *testing.T) {
	l := New()
	l.head.Value = 0
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	testcase := []struct {
		baseIdx            int
		moveIdx            int
		expectStringOfNext string
		expectStringOfPrev string
	}{
		{
			baseIdx:            0,
			moveIdx:            0,
			expectStringOfNext: "0->1->2->3->0",
			expectStringOfPrev: "0->3->2->1->0",
		},
		{
			baseIdx:            2,
			moveIdx:            1,
			expectStringOfNext: "0->1->2->3->0",
			expectStringOfPrev: "0->3->2->1->0",
		},
		{
			baseIdx:            1,
			moveIdx:            2,
			expectStringOfNext: "0->2->1->3->0",
			expectStringOfPrev: "0->3->1->2->0",
		},
		{
			baseIdx:            1,
			moveIdx:            2,
			expectStringOfNext: "0->1->2->3->0",
			expectStringOfPrev: "0->3->2->1->0",
		},
		{
			baseIdx:            1,
			moveIdx:            3,
			expectStringOfNext: "0->3->1->2->0",
			expectStringOfPrev: "0->2->1->3->0",
		},
	}
	for i, test := range testcase {
		l.MoveBefore(nexti(l.head, test.baseIdx), nexti(l.head, test.moveIdx))
		if toStringByNext(l) != test.expectStringOfNext {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfNext, toStringByNext(l))
		}
		if toStringByPrev(l) != test.expectStringOfPrev {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfPrev, toStringByPrev(l))
		}
	}
}

func TestMoveToFront(t *testing.T) {
	l := New()
	l.head.Value = 0
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	testcase := []struct {
		moveIdx            int
		expectStringOfNext string
		expectStringOfPrev string
	}{
		{
			moveIdx:            0,
			expectStringOfNext: "0->1->2->3->0",
			expectStringOfPrev: "0->3->2->1->0",
		},
		{
			moveIdx:            1,
			expectStringOfNext: "0->1->2->3->0",
			expectStringOfPrev: "0->3->2->1->0",
		},
		{
			moveIdx:            2,
			expectStringOfNext: "0->2->1->3->0",
			expectStringOfPrev: "0->3->1->2->0",
		},
		{
			moveIdx:            3,
			expectStringOfNext: "0->3->2->1->0",
			expectStringOfPrev: "0->1->2->3->0",
		},
	}
	for i, test := range testcase {
		l.MoveToFront(nexti(l.head, test.moveIdx))
		if toStringByNext(l) != test.expectStringOfNext {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfNext, toStringByNext(l))
		}
		if toStringByPrev(l) != test.expectStringOfPrev {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfPrev, toStringByPrev(l))
		}
	}
}

func TestMoveToBack(t *testing.T) {
	l := New()
	l.head.Value = 0
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	testcase := []struct {
		moveIdx            int
		expectStringOfNext string
		expectStringOfPrev string
	}{
		{
			moveIdx:            0,
			expectStringOfNext: "0->1->2->3->0",
			expectStringOfPrev: "0->3->2->1->0",
		},
		{
			moveIdx:            1,
			expectStringOfNext: "0->2->3->1->0",
			expectStringOfPrev: "0->1->3->2->0",
		},
		{
			moveIdx:            2,
			expectStringOfNext: "0->2->1->3->0",
			expectStringOfPrev: "0->3->1->2->0",
		},
		{
			moveIdx:            3,
			expectStringOfNext: "0->2->1->3->0",
			expectStringOfPrev: "0->3->1->2->0",
		},
	}
	for i, test := range testcase {
		l.MoveToBack(nexti(l.head, test.moveIdx))
		if toStringByNext(l) != test.expectStringOfNext {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfNext, toStringByNext(l))
		}
		if toStringByPrev(l) != test.expectStringOfPrev {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfPrev, toStringByPrev(l))
		}
	}
}

func TestRemove(t *testing.T) {
	l := New()
	l.head.Value = 0
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	testcase := []struct {
		removeIdx          int
		expectStringOfNext string
		expectStringOfPrev string
		expectSize         int
	}{
		{
			removeIdx:          1,
			expectStringOfNext: "0->2->3->0",
			expectStringOfPrev: "0->3->2->0",
			expectSize:         2,
		},
		{
			removeIdx:          2,
			expectStringOfNext: "0->2->0",
			expectStringOfPrev: "0->2->0",
			expectSize:         1,
		},
		{
			removeIdx:          1,
			expectStringOfNext: "0->0",
			expectStringOfPrev: "0->0",
			expectSize:         0,
		},
	}
	for i, test := range testcase {
		l.Remove(nexti(l.head, test.removeIdx))
		if toStringByNext(l) != test.expectStringOfNext {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfNext, toStringByNext(l))
		}
		if toStringByPrev(l) != test.expectStringOfPrev {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfPrev, toStringByPrev(l))
		}
		if l.size != test.expectSize {
			t.Errorf("%d Expected %d, got %d", i, test.expectSize, l.size)
		}
	}
}

func TestPopFront(t *testing.T) {
	l := New()
	l.head.Value = 0
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	testcase := []struct {
		expectValue        interface{}
		expectStringOfNext string
		expectStringOfPrev string
		expectSize         int
	}{
		{
			expectValue:        1,
			expectStringOfNext: "0->2->3->0",
			expectStringOfPrev: "0->3->2->0",
			expectSize:         2,
		},
		{
			expectValue:        2,
			expectStringOfNext: "0->3->0",
			expectStringOfPrev: "0->3->0",
			expectSize:         1,
		},
		{
			expectValue:        3,
			expectStringOfNext: "0->0",
			expectStringOfPrev: "0->0",
			expectSize:         0,
		},
		{
			expectValue:        nil,
			expectStringOfNext: "0->0",
			expectStringOfPrev: "0->0",
			expectSize:         0,
		},
	}
	for i, test := range testcase {
		e := l.PopFront()
		if e != test.expectValue {
			t.Errorf("%d Expected %d, got %v", i, test.expectValue, e)
		}
		if toStringByNext(l) != test.expectStringOfNext {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfNext, toStringByNext(l))
		}
		if toStringByPrev(l) != test.expectStringOfPrev {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfPrev, toStringByPrev(l))
		}
		if l.size != test.expectSize {
			t.Errorf("%d Expected %d, got %d", i, test.expectSize, l.size)
		}
	}
}

func TestPopBack(t *testing.T) {
	l := New()
	l.head.Value = 0
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	testcase := []struct {
		expectValue        interface{}
		expectStringOfNext string
		expectStringOfPrev string
		expectSize         int
	}{
		{
			expectValue:        3,
			expectStringOfNext: "0->1->2->0",
			expectStringOfPrev: "0->2->1->0",
			expectSize:         2,
		},
		{
			expectValue:        2,
			expectStringOfNext: "0->1->0",
			expectStringOfPrev: "0->1->0",
			expectSize:         1,
		},
		{
			expectValue:        1,
			expectStringOfNext: "0->0",
			expectStringOfPrev: "0->0",
			expectSize:         0,
		},
		{
			expectValue:        nil,
			expectStringOfNext: "0->0",
			expectStringOfPrev: "0->0",
			expectSize:         0,
		},
	}
	for i, test := range testcase {
		e := l.PopBack()
		if e != test.expectValue {
			t.Errorf("%d Expected %d, got %v", i, test.expectValue, e)
		}
		if toStringByNext(l) != test.expectStringOfNext {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfNext, toStringByNext(l))
		}
		if toStringByPrev(l) != test.expectStringOfPrev {
			t.Errorf("%d Expected %s, got %s", i, test.expectStringOfPrev, toStringByPrev(l))
		}
		if l.size != test.expectSize {
			t.Errorf("%d Expected %d, got %d", i, test.expectSize, l.size)
		}
	}
}
