// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package skiplist

import "github.com/mrtcx/plusdata/tree"

type element struct {
	sk   *skipList
	node *Node
}

func (e *element) Key() interface{} {
	return e.node.key
}

func (e *element) Value() interface{} {
	return e.node.value
}

func (e *element) SetValue(value interface{}) {
	e.node.value = value
}

func (e *element) Prev() tree.Element {
	return e.sk.Prev(e.node.key)
}

func (e *element) Next() tree.Element {
	if e.node.nexts[0] == nil {
		return nil
	}
	return &element{
		sk:   e.sk,
		node: e.node.nexts[0],
	}
}

func (s *skipList) Find(key interface{}) tree.Element {
	node := s.findNode(key)
	if node == nil {
		return nil
	}
	return &element{sk: s, node: node}
}

func (s *skipList) Left() tree.Element {
	node := s.frontNode()
	if node == nil {
		return nil
	}
	return &element{sk: s, node: node}
}

func (s *skipList) Right() tree.Element {
	node := s.backNode()
	if node == nil {
		return nil
	}
	return &element{sk: s, node: node}
}

func (s *skipList) Prev(key interface{}) tree.Element {
	node := s.findPrevNode(key)
	if node == nil {
		return nil
	}
	return &element{sk: s, node: node}
}

func (s *skipList) Next(key interface{}) tree.Element {
	node := s.findNextNode(key)
	if node == nil {
		return nil
	}
	return &element{sk: s, node: node}
}
