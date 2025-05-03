// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package singlelinkedlist

import (
	"fmt"
	"testing"
)

func toString(l *List) string {
	s := ""
	for e := l.head; e != nil; e = e.next {
		if e != l.head {
			s += "->"
		}
		s += fmt.Sprintf("%d", e.Value)
	}
	return s
}

func TestNew(t *testing.T) {
	list := New()
	if list.head != nil || list.tail != nil || list.size != 0 {
		t.Errorf("Expected empty")
	}
}

func TestPushFront(t *testing.T) {
	testCase := []struct {
		pushValue    int
		expectString string
		expectSize   int
	}{
		{
			pushValue:    1,
			expectString: "1",
			expectSize:   1,
		},
		{
			pushValue:    2,
			expectString: "2->1",
			expectSize:   2,
		},
		{
			pushValue:    3,
			expectString: "3->2->1",
			expectSize:   3,
		},
		{
			pushValue:    4,
			expectString: "4->3->2->1",
			expectSize:   4,
		},
	}
	list := New()
	for _, test := range testCase {
		list.PushFront(test.pushValue)
		if toString(list) != test.expectString {
			t.Errorf("TestPushFront failed for value %d", test.pushValue)
		}
	}
}

func TestPushBack(t *testing.T) {
	testCase := []struct {
		pushValue    int
		expectString string
		expectSize   int
	}{
		{
			pushValue:    1,
			expectString: "1",
			expectSize:   1,
		},
		{
			pushValue:    2,
			expectString: "1->2",
			expectSize:   2,
		},
		{
			pushValue:    3,
			expectString: "1->2->3",
			expectSize:   3,
		},
		{
			pushValue:    4,
			expectString: "1->2->3->4",
			expectSize:   4,
		},
	}
	list := New()
	for _, test := range testCase {
		list.PushBack(test.pushValue)
		if toString(list) != test.expectString {
			t.Errorf("TestPushBack failed for value %d", test.pushValue)
		}
		if list.size != test.expectSize {
			t.Errorf("TestPushBack failed for value %d", test.pushValue)
		}
	}
}

func TestPopFront(t *testing.T) {
	l := New()
	l.PushFront(1)
	l.PushFront(2)
	l.PushFront(3)
	testcase := []struct {
		expectPopValue interface{}
		expectString   string
		expectSize     int
	}{
		{
			expectPopValue: 3,
			expectString:   "2->1",
			expectSize:     2,
		},
		{
			expectPopValue: 2,
			expectString:   "1",
			expectSize:     1,
		},
		{
			expectPopValue: 1,
			expectString:   "",
			expectSize:     0,
		},
		{
			expectPopValue: nil,
			expectString:   "",
			expectSize:     0,
		},
	}
	for _, test := range testcase {
		p := l.PopFront()
		if p != test.expectPopValue {
			t.Errorf("TestPopFront failed for value %d", test.expectPopValue)
		}
		if toString(l) != test.expectString {
			t.Errorf("TestPopFront failed for value %d", test.expectPopValue)
		}
		if l.size != test.expectSize {
			t.Errorf("TestPopFront failed for value %d", test.expectPopValue)
		}
	}
}

func TestFront(t *testing.T) {
	l := New()
	if l.Front() != nil {
		t.Errorf("TestFront nil")
	}
	l.PushFront(1)
	if l.Front() != l.head || l.Front().Value != 1 {
		t.Errorf("TestFront failed for value %d", l.Front().Value)
	}
}

func TestBack(t *testing.T) {
	l := New()
	if l.Back() != nil {
		t.Errorf("TestBack nil")
	}
	l.PushBack(1)
	if l.Back() != l.tail || l.Back().Value != 1 {
		t.Errorf("TestBack failed")
	}
}

func TestMerge(t *testing.T) {
	l1 := New()
	l2 := New()
	l1.Merge(l2)
	if l1.head != nil || l1.tail != nil || l1.size != 0 {
		t.Errorf("expected empty")
	}
	l2.PushFront(1)
	l1.Merge(l2)
	if toString(l1) != "1" {
		t.Errorf("expected 1 got %s", toString(l1))
	}
	if l2.head != nil || l2.tail != nil || l2.size != 0 {
		t.Errorf("expected empty")
	}
	l3 := New()
	l3.PushBack(2)
	l1.Merge(l3)
	if toString(l1) != "1->2" {
		t.Errorf("expected 1->2 got %s", toString(l1))
	}
}
