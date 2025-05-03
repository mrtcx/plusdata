// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package singlelinkedlist

type List struct {
	head *Element
	tail *Element
	size int
}

type Element struct {
	Value interface{}
	next  *Element
}

func New() *List {
	return &List{}
}

func (l *List) Clean() {
	l.head, l.tail, l.size = nil, nil, 0
}

func (l *List) Size() int {
	return l.size
}

func (l *List) Empty() bool {
	return l.size == 0
}

func (l *List) Front() *Element {
	return l.head
}

func (l *List) Back() *Element {
	return l.tail
}

func (l *List) PushFront(value interface{}) *Element {
	e := &Element{
		Value: value,
	}
	if l.size == 0 {
		l.head, l.tail = e, e
	} else {
		e.next, l.head = l.head, e
	}
	l.size++
	return e
}

func (l *List) PushBack(value interface{}) *Element {
	if l.size == 0 {
		return l.PushFront(value)
	}
	e := &Element{
		Value: value,
	}
	l.tail.next, l.tail = e, e
	l.size++
	return e
}

func (l *List) PopFront() interface{} {
	if l.size != 0 {
		e := l.head
		l.head = l.head.next
		if l.head == nil {
			l.tail = nil
		}
		e.next = nil
		l.size--
		return e.Value
	}
	return nil
}

func (l *List) InsertAfter(pos *Element, value interface{}) *Element {
	newe := &Element{Value: value}
	pos.next, newe.next = newe, pos.next
	if newe.next == nil {
		l.tail = newe
	}
	l.size++
	return newe
}

func (l *List) Merge(delete *List) {
	if delete.size == 0 {
		return
	}
	if l.size == 0 {
		l.head, l.tail, l.size = delete.head, delete.tail, delete.size
	} else {
		l.tail.next = delete.head
		l.tail = delete.tail
		l.size += delete.size
	}
}

func (e *Element) Next() *Element {
	return e.next
}
