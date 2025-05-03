// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package skiplist

import (
	"math/rand"

	"github.com/mrtcx/plusdata/tree"
)

var _ tree.Tree = (*skipList)(nil)

const (
	_maxLevel    = 32     // 最大层数
	_probability = 0x4000 // 25%的概率阈值 (0xFFFF * 0.25)
)

type skipList struct {
	size int
	cmp  tree.Comparator
	head *Node
	tail *Node
}

type Node struct {
	key   interface{}
	value interface{}
	nexts []*Node
}

func New(cmp tree.Comparator) *skipList {
	return &skipList{
		cmp: cmp,
		head: &Node{
			nexts: make([]*Node, 1),
		},
	}
}

func (s *skipList) Clean() {
	s.size, s.head = 0, &Node{nexts: make([]*Node, 1)}
}

func (s *skipList) Size() int {
	return s.size
}

func (s *skipList) Empty() bool {
	return s.size == 0
}

func (s *skipList) Insert(key interface{}, value interface{}) {
	pres := s.levelPreNodes(key)
	if pres[0].nexts[0] != nil && s.cmp(pres[0].nexts[0].key, key) == 0 {
		pres[0].nexts[0].value = value
		return
	}
	level := randomLevel()
	insertNode := &Node{key: key, value: value, nexts: make([]*Node, level)}
	if level < len(pres) {
		pres = pres[:level]
	}
	for i := len(pres) - 1; i >= 0; i-- {
		pres[i].nexts[i], insertNode.nexts[i] = insertNode, pres[i].nexts[i]
	}
	for i := len(s.head.nexts); i < level; i++ {
		s.head.nexts = append(s.head.nexts, insertNode)
		s.head.nexts[i] = insertNode
	}
	if insertNode.nexts[0] == nil {
		s.tail = insertNode
	}
	s.size++
}

func (s *skipList) Remove(key interface{}) {
	levelPres := s.levelPreNodes(key)
	if levelPres[0].nexts[0] == nil || s.cmp(levelPres[0].nexts[0].key, key) != 0 {
		return
	}
	for i := len(levelPres) - 1; i >= 0; i-- {
		if levelPres[i].nexts[i] != nil {
			levelPres[i].nexts[i] = levelPres[i].nexts[i].nexts[i]
		}
	}
	if levelPres[0].nexts[0] == nil {
		s.tail = levelPres[0]
	}
	s.size--
}

func (s *skipList) Get(key interface{}) (interface{}, bool) {
	node := s.findNode(key)
	if node == nil {
		return nil, false
	}
	return node.value, true
}

func (s *skipList) frontNode() *Node {
	if s.size == 0 {
		return nil
	}
	return s.head.nexts[0]
}

func (s *skipList) backNode() *Node {
	if s.size == 0 {
		return nil
	}
	return s.tail
}

func (s *skipList) findNode(key interface{}) *Node {
	pre := s.preLocate(key)
	if pre.nexts[0] != nil && s.cmp(pre.nexts[0].key, key) == 0 {
		return pre.nexts[0]
	}
	return nil
}

func (s *skipList) findPrevNode(key interface{}) *Node {
	pre := s.preLocate(key)
	if pre == s.head {
		return nil
	}
	return pre
}

func (s *skipList) findNextNode(key interface{}) *Node {
	pre := s.preLocate(key)
	if pre.nexts[0] != nil {
		if s.cmp(pre.nexts[0].key, key) > 0 {
			return pre.nexts[0]
		}
		return pre.nexts[0].nexts[0]
	}
	return nil
}

func (s *skipList) levelPreNodes(key interface{}) []*Node {
	pre := s.head
	var pres []*Node = make([]*Node, len(pre.nexts))
	for i := len(pre.nexts) - 1; i >= 0; i-- {
		for pre.nexts[i] != nil && s.cmp(key, pre.nexts[i].key) > 0 {
			pre = pre.nexts[i]
		}
		pres[i] = pre
	}
	return pres
}

func (s *skipList) preLocate(key interface{}) *Node {
	pre := s.head
	for i := len(pre.nexts) - 1; i >= 0; i-- {
		for pre.nexts[i] != nil && s.cmp(key, pre.nexts[i].key) > 0 {
			pre = pre.nexts[i]
		}
	}
	return pre
}

func randomLevel() int {
	level := 1
	for rand.Int31n(0xFFFF) < _probability && level < _maxLevel {
		level++
	}
	return level
}
