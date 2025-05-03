// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package doublelinkedlist

type List struct {
	head *Element
	size int
}

type Element struct {
	Value interface{}
	prev  *Element
	next  *Element
	l     *List
}

func (e *Element) Next() *Element {
	nxt := e.next
	if nxt == e.l.head {
		return nil
	}
	return nxt
}

func (e *Element) Prev() *Element {
	prv := e.prev
	if prv == e.l.head {
		return nil
	}
	return prv
}

// 哨兵
func sentry(l *List) *Element {
	e := &Element{}
	e.prev, e.next, e.l = e, e, l
	return e
}

func New() *List {
	l := &List{}
	l.head = sentry(l)
	return l
}

func (l *List) Size() int {
	return l.size
}

func (l *List) Front() *Element {
	if l.size == 0 {
		return nil
	}
	return l.head.next
}

func (l *List) Back() *Element {
	if l.size == 0 {
		return nil
	}
	return l.head.prev
}

func (l *List) InsertAfter(e *Element, value interface{}) *Element {
	if e.l != l {
		panic("element not in list")
	}
	ne := &Element{Value: value, l: l}
	e.next, e.next.prev, ne.next, ne.prev = ne, ne, e.next, e
	l.size++
	return ne
}

func (l *List) InsertBefore(e *Element, value interface{}) *Element {
	return l.InsertAfter(e.prev, value)
}

func (l *List) PushFront(value interface{}) *Element {
	return l.InsertAfter(l.head, value)
}

func (l *List) PushBack(value interface{}) *Element {
	return l.InsertAfter(l.head.prev, value)
}

func (l *List) MoveAfter(e, m *Element) *Element {
	if e.l != l || m.l != l {
		panic("element not in list")
	}
	if e == m || e.next == m {
		return e
	}
	e.next, e.next.prev, m.next, m.prev, m.prev.next, m.next.prev = m, m, e.next, e, m.next, m.prev
	return e
}

func (l *List) MoveBefore(e, m *Element) *Element {
	return l.MoveAfter(e.prev, m)
}

func (l *List) MoveToFront(e *Element) *Element {
	return l.MoveAfter(l.head, e)
}

func (l *List) MoveToBack(e *Element) *Element {
	return l.MoveBefore(l.head, e)
}

func (l *List) Remove(e *Element) {
	if l != e.l {
		panic("element not in list")
	}
	e.prev.next, e.next.prev = e.next, e.prev
	e.next, e.prev, e.l = nil, nil, nil
	l.size--
}

func (l *List) PopFront() interface{} {
	if l.size == 0 {
		return nil
	}
	e := l.head.next
	l.Remove(e)
	return e.Value
}

func (l *List) PopBack() interface{} {
	if l.size == 0 {
		return nil
	}
	e := l.head.prev
	l.Remove(e)
	return e.Value
}

func (l *List) Merge(delete *List) {
	if delete.size == 0 {
		return
	}
	if l.size == 0 {
		l.head, l.size = delete.head, delete.size
	} else {
		l.head.prev.next, delete.head.next.prev = delete.head.next, l.head.next
		l.head.prev, delete.head.prev.next = delete.head.prev, l.head
		l.size++
	}
}
