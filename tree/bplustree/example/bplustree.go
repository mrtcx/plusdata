// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
)

type Comparator func(a, b interface{}) int

type BPlusTree struct {
	root       *Node
	order      int
	minKeys    int
	maxKeys    int
	comparator Comparator
}

type Node struct {
	isLeaf   bool
	keys     []interface{}
	values   []interface{}
	children []*Node
	next     *Node
	prev     *Node
	parent   *Node
}

func NewBPlusTree(order int, comparator Comparator) *BPlusTree {
	return &BPlusTree{
		order:      order,
		minKeys:    int(math.Ceil(float64(order)/2)) - 1,
		maxKeys:    order - 1,
		comparator: comparator,
	}
}

func (t *BPlusTree) Insert(key, value interface{}) {
	if t.root == nil {
		t.root = t.createLeafNode()
		t.root.keys = append(t.root.keys, key)
		t.root.values = append(t.root.values, value)
		return
	}

	leaf := t.findLeafNode(key)
	if t.insertIntoLeaf(leaf, key, value) {
		if len(leaf.keys) > t.maxKeys {
			t.splitLeafNode(leaf)
		}
	}
}

func (t *BPlusTree) Delete(key interface{}) {
	if t.root == nil {
		return
	}

	leaf := t.findLeafNode(key)
	if idx, found := t.findKeyIndex(leaf, key); found {
		t.deleteFromLeaf(leaf, idx)
		if len(leaf.keys) < t.minKeys {
			t.handleUnderflow(leaf)
		}
	}
}

func (t *BPlusTree) Search(key interface{}) (interface{}, bool) {
	if t.root == nil {
		return nil, false
	}

	leaf := t.findLeafNode(key)
	if idx, found := t.findKeyIndex(leaf, key); found {
		return leaf.values[idx], true
	}
	return nil, false
}

func (t *BPlusTree) Update(key, value interface{}) {
	t.Delete(key)
	t.Insert(key, value)
}

// 内部实现方法
func (t *BPlusTree) createLeafNode() *Node {
	return &Node{
		isLeaf: true,
		keys:   make([]interface{}, 0),
		values: make([]interface{}, 0),
	}
}

func (t *BPlusTree) findLeafNode(key interface{}) *Node {
	current := t.root
	for !current.isLeaf {
		idx := 0
		for idx < len(current.keys) && t.comparator(key, current.keys[idx]) >= 0 {
			idx++
		}
		current = current.children[idx]
	}
	return current
}

func (t *BPlusTree) findKeyIndex(node *Node, key interface{}) (int, bool) {
	for i, k := range node.keys {
		if t.comparator(k, key) == 0 {
			return i, true
		}
	}
	return -1, false
}

func (t *BPlusTree) insertIntoLeaf(leaf *Node, key, value interface{}) bool {
	idx := 0
	for idx < len(leaf.keys) && t.comparator(key, leaf.keys[idx]) > 0 {
		idx++
	}

	leaf.keys = append(leaf.keys[:idx], append([]interface{}{key}, leaf.keys[idx:]...)...)
	leaf.values = append(leaf.values[:idx], append([]interface{}{value}, leaf.values[idx:]...)...)
	return true
}

func (t *BPlusTree) splitLeafNode(original *Node) {
	splitIdx := len(original.keys) / 2
	newNode := t.createLeafNode()

	newNode.keys = append(newNode.keys, original.keys[splitIdx:]...)
	newNode.values = append(newNode.values, original.values[splitIdx:]...)
	original.keys = original.keys[:splitIdx]
	original.values = original.values[:splitIdx]

	newNode.next = original.next
	if original.next != nil {
		original.next.prev = newNode
	}
	original.next = newNode
	newNode.prev = original

	t.insertIntoParent(original, newNode.keys[0], newNode)
}

func (t *BPlusTree) splitInternalNode(original *Node) {
	splitIdx := len(original.keys) / 2
	pivotKey := original.keys[splitIdx]

	newNode := &Node{
		isLeaf:   false,
		keys:     original.keys[splitIdx+1:],
		children: original.children[splitIdx+1:],
		parent:   original.parent,
	}

	for _, child := range newNode.children {
		child.parent = newNode
	}

	original.keys = original.keys[:splitIdx]
	original.children = original.children[:splitIdx+1]

	t.insertIntoParent(original, pivotKey, newNode)
}

func (t *BPlusTree) insertIntoParent(left *Node, key interface{}, right *Node) {
	parent := left.parent
	if parent == nil {
		t.root = &Node{
			keys:     []interface{}{key},
			children: []*Node{left, right},
		}
		left.parent = t.root
		right.parent = t.root
		return
	}

	idx := 0
	for idx < len(parent.keys) && t.comparator(key, parent.keys[idx]) >= 0 {
		idx++
	}

	parent.keys = append(parent.keys[:idx], append([]interface{}{key}, parent.keys[idx:]...)...)
	parent.children = append(parent.children[:idx+1], append([]*Node{right}, parent.children[idx+1:]...)...)

	if len(parent.keys) > t.maxKeys {
		t.splitInternalNode(parent)
	}
}

func (t *BPlusTree) deleteFromLeaf(leaf *Node, idx int) {
	leaf.keys = append(leaf.keys[:idx], leaf.keys[idx+1:]...)
	leaf.values = append(leaf.values[:idx], leaf.values[idx+1:]...)
}

func (t *BPlusTree) handleUnderflow(node *Node) {
	if node == t.root {
		if len(node.keys) == 0 && len(node.children) > 0 {
			t.root = node.children[0]
			t.root.parent = nil
		}
		return
	}

	parent := node.parent
	idx := 0
	for ; idx < len(parent.children) && parent.children[idx] != node; idx++ {
	}

	// Try borrow left
	if idx > 0 {
		leftSibling := parent.children[idx-1]
		if len(leftSibling.keys) > t.minKeys {
			t.borrowFromLeft(node, leftSibling, idx-1)
			return
		}
	}

	// Try borrow right
	if idx < len(parent.children)-1 {
		rightSibling := parent.children[idx+1]
		if len(rightSibling.keys) > t.minKeys {
			t.borrowFromRight(node, rightSibling, idx)
			return
		}
	}

	// Merge nodes
	if idx > 0 {
		t.mergeNodes(parent.children[idx-1], node, idx-1)
	} else {
		t.mergeNodes(node, parent.children[idx+1], idx)
	}
}

func (t *BPlusTree) borrowFromLeft(node, leftSibling *Node, parentIdx int) {
	if node.isLeaf {
		lastKey := leftSibling.keys[len(leftSibling.keys)-1]
		lastValue := leftSibling.values[len(leftSibling.values)-1]

		leftSibling.keys = leftSibling.keys[:len(leftSibling.keys)-1]
		leftSibling.values = leftSibling.values[:len(leftSibling.values)-1]

		node.keys = append([]interface{}{lastKey}, node.keys...)
		node.values = append([]interface{}{lastValue}, node.values...)

		node.parent.keys[parentIdx] = node.keys[0]
	} else {
		lastKey := leftSibling.keys[len(leftSibling.keys)-1]
		lastChild := leftSibling.children[len(leftSibling.children)-1]

		leftSibling.keys = leftSibling.keys[:len(leftSibling.keys)-1]
		leftSibling.children = leftSibling.children[:len(leftSibling.children)-1]

		node.keys = append([]interface{}{node.parent.keys[parentIdx]}, node.keys...)
		node.children = append([]*Node{lastChild}, node.children...)
		lastChild.parent = node

		node.parent.keys[parentIdx] = lastKey
	}
}

func (t *BPlusTree) borrowFromRight(node, rightSibling *Node, parentIdx int) {
	if node.isLeaf {
		firstKey := rightSibling.keys[0]
		firstValue := rightSibling.values[0]

		rightSibling.keys = rightSibling.keys[1:]
		rightSibling.values = rightSibling.values[1:]

		node.keys = append(node.keys, firstKey)
		node.values = append(node.values, firstValue)

		node.parent.keys[parentIdx] = rightSibling.keys[0]
	} else {
		firstKey := rightSibling.keys[0]
		firstChild := rightSibling.children[0]

		rightSibling.keys = rightSibling.keys[1:]
		rightSibling.children = rightSibling.children[1:]

		node.keys = append(node.keys, node.parent.keys[parentIdx])
		node.children = append(node.children, firstChild)
		firstChild.parent = node

		node.parent.keys[parentIdx] = firstKey
	}
}

func (t *BPlusTree) mergeNodes(left, right *Node, parentIdx int) {
	parent := left.parent
	if left.isLeaf {
		left.keys = append(left.keys, right.keys...)
		left.values = append(left.values, right.values...)
		left.next = right.next
		if right.next != nil {
			right.next.prev = left
		}
	} else {
		left.keys = append(left.keys, parent.keys[parentIdx])
		left.keys = append(left.keys, right.keys...)
		left.children = append(left.children, right.children...)
		for _, child := range right.children {
			child.parent = left
		}
	}

	parent.keys = append(parent.keys[:parentIdx], parent.keys[parentIdx+1:]...)
	parent.children = append(parent.children[:parentIdx+1], parent.children[parentIdx+2:]...)

	if len(parent.keys) < t.minKeys {
		t.handleUnderflow(parent)
	}
}

// 测试用比较器
func intComparator(a, b interface{}) int {
	ai := a.(int)
	bi := b.(int)
	switch {
	case ai < bi:
		return -1
	case ai > bi:
		return 1
	default:
		return 0
	}
}

func main() {
	bptree := NewBPlusTree(3, intComparator)

	// 测试插入和递归分裂
	for i := 1; i <= 10; i++ {
		bptree.Insert(i, fmt.Sprintf("V%d", i))
	}

	// 验证查询
	for i := 1; i <= 10; i++ {
		if val, found := bptree.Search(i); found {
			fmt.Printf("Key %d: %s\n", i, val)
		}
	}

	// 测试删除和递归合并
	for i := 10; i > 3; i-- {
		bptree.Delete(i)
	}

	// 验证删除结果
	fmt.Println("\nAfter deletion:")
	for i := 1; i <= 10; i++ {
		if val, found := bptree.Search(i); found {
			fmt.Printf("Key %d: %s\n", i, val)
		}
	}
}
