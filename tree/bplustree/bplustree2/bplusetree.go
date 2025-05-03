// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bplustree

import (
	"github.com/mrtcx/plusdata/tree"
)

type BPlusTree struct {
	root  *Node
	cmp   tree.Comparator
	size  int
	order int
}

type Node struct {
	keys   []interface{}
	vals   []interface{}
	childs []*Node
	next   *Node
}

func (node *Node) isLeaf() bool {
	return len(node.childs) == 0
}

func New(cmp tree.Comparator, order int) *BPlusTree {
	return &BPlusTree{
		cmp:   cmp,
		order: order,
	}
}

func (bp *BPlusTree) Clean() {
	bp.root, bp.size = nil, 0
}

func (bp *BPlusTree) Insert(key, val interface{}) {
	root := bp.insert(bp.root, key, val)
	if len(root.keys) > int(bp.maxKeys()) {
		key, left, right := split(root)
		root = &Node{
			keys:   []interface{}{key, right.keys[len(right.keys)-1]},
			childs: []*Node{left, right},
		}
	}
	bp.root = root
}

func (bp *BPlusTree) insert(root *Node, key, val interface{}) *Node {
	if root == nil {
		bp.size++
		return &Node{
			keys: []interface{}{key},
			vals: []interface{}{val},
		}
	}
	idx := 0
	for ; idx < len(root.keys) && bp.cmp(key, root.keys[idx]) > 0; idx++ {
	}
	if root.isLeaf() {
		if idx < len(root.keys) && bp.cmp(root.keys[idx], key) == 0 {
			root.vals[idx] = val
			return root
		}
		bp.size++
		root.keys = increaseSpace(root.keys, idx)
		root.vals = increaseSpace(root.vals, idx)
		root.keys[idx], root.vals[idx] = key, val
		return root
	} else {
		if idx == len(root.keys) && bp.cmp(key, root.keys[idx-1]) > 0 {
			idx -= 1
			root.keys[idx] = key
		}
		retNode := bp.insert(root.childs[idx], key, val)
		if len(retNode.keys) > bp.maxKeys() {
			key, left, right := split(retNode)
			root.keys = increaseSpace(root.keys, idx)
			root.childs = increaseSpace2(root.childs, idx)
			root.keys[idx] = key
			root.childs[idx], root.childs[idx+1] = left, right
		}
		return root
	}
}

func (bp *BPlusTree) Remove(key interface{}) {
	if bp.root == nil {
		return
	}
	root := bp.remove(bp.root, key)
	if len(root.keys) == 0 {
		bp.root = nil
		return
	}
	if root.isLeaf() {
		return
	}
	if len(root.keys) == 1 {
		bp.root = root.childs[0]
	}
}

func (bp *BPlusTree) remove(root *Node, key interface{}) *Node {
	if root == nil {
		return nil
	}
	idx := 0
	for ; idx < len(root.keys) && bp.cmp(key, root.keys[idx]) > 0; idx++ {
	}
	if root.isLeaf() {
		if idx < len(root.keys) && bp.cmp(root.keys[idx], key) == 0 {
			bp.size--
			root.keys = decreaseSpace(root.keys, idx)
			root.vals = decreaseSpace(root.vals, idx)
		}
	} else {
		if idx < len(root.keys) {
			retNode := bp.remove(root.childs[idx], key)
			if len(retNode.keys) < bp.minKeys() {
				minkeys, maxkeys := bp.minKeys(), bp.maxKeys()
				if !mergeRight(root, idx, maxkeys) &&
					!mergeLeft(root, idx, maxkeys) &&
					!borrowFromLeft(root, idx, minkeys) &&
					!borrowFromRight(root, idx, minkeys) {
					root.keys = decreaseSpace(root.keys, idx)
					root.childs = decreaseSpace2(root.childs, idx)
				}
			} else {
				root.keys[idx] = root.childs[idx].keys[len(root.childs[idx].keys)-1]
			}
		}
	}
	return root
}

func (bp *BPlusTree) Get(key interface{}) (interface{}, bool) {
	root := bp.root
	if root == nil {
		return nil, false
	}
	for {
		if root.isLeaf() {
			for i := 0; i < len(root.keys); i++ {
				if bp.cmp(root.keys[i], key) == 0 {
					return root.vals[i], true
				}
			}
			return nil, false
		}
		idx := 0
		for ; idx < len(root.keys) && bp.cmp(key, root.keys[idx]) > 0; idx++ {
		}
		if idx == len(root.keys) {
			return nil, false
		}
		root = root.childs[idx]
	}
}

func split(node *Node) (key interface{}, left *Node, right *Node) {
	pivot := len(node.keys) / 2
	key = node.keys[pivot-1]
	right = &Node{}
	right.keys = append(right.keys, node.keys[pivot:]...)
	node.keys = node.keys[:pivot]
	if node.isLeaf() {
		right.vals = append(right.vals, node.vals[pivot:]...)
		node.vals = node.vals[:pivot]
	} else {
		right.childs = append(right.childs, node.childs[pivot:]...)
		node.childs = node.childs[:pivot]
	}
	left = node
	left.next, right.next = right, left.next
	return key, left, right
}

func borrowFromRight(parent *Node, idx int, minKeys int) bool {
	if idx == len(parent.keys)-1 {
		return false
	}
	child := parent.childs[idx]
	rightChild := parent.childs[idx+1]
	if len(rightChild.keys) <= minKeys {
		return false
	}
	child.keys = append(child.keys, rightChild.keys[0])
	rightChild.keys = rightChild.keys[1:]
	if child.isLeaf() {
		child.vals = append(child.vals, rightChild.vals[0])
		rightChild.vals = rightChild.vals[1:]
	} else {
		child.childs = append(child.childs, rightChild.childs[0])
		rightChild.childs = rightChild.childs[1:]
	}
	parent.keys[idx] = child.keys[len(child.keys)-1]
	return true
}

func borrowFromLeft(parent *Node, idx int, minKeys int) bool {
	if idx == 0 {
		return false
	}
	leftChild := parent.childs[idx-1]
	child := parent.childs[idx]
	if len(leftChild.keys) <= minKeys {
		return false
	}
	child.keys = increaseSpace(child.keys, 0)
	child.keys[0] = leftChild.keys[len(leftChild.keys)-1]
	leftChild.keys = leftChild.keys[:len(leftChild.keys)-1]
	if child.isLeaf() {
		child.vals = increaseSpace(child.vals, 0)
		child.vals[0] = leftChild.vals[len(leftChild.vals)-1]
		leftChild.vals = leftChild.vals[:len(leftChild.vals)-1]
	} else {
		child.childs = increaseSpace2(child.childs, 0)
		child.childs[0] = leftChild.childs[len(leftChild.childs)-1]
		leftChild.childs = leftChild.childs[:len(leftChild.childs)-1]
	}
	parent.keys[idx] = child.keys[len(child.keys)-1]
	parent.keys[idx-1] = leftChild.keys[len(leftChild.keys)-1]
	return true
}

func mergeLeft(parnet *Node, idx int, maxKeys int) bool {
	if idx == 0 {
		return false
	}
	leftChild := parnet.childs[idx-1]
	child := parnet.childs[idx]
	if len(leftChild.keys)+len(child.keys) > maxKeys {
		return false
	}
	leftChild.keys = append(leftChild.keys, child.keys...)
	if child.isLeaf() {
		leftChild.vals = append(leftChild.vals, child.vals...)
	} else {
		leftChild.childs = append(leftChild.childs, child.childs...)
	}
	leftChild.next = child.next
	parnet.keys[idx-1] = leftChild.keys[len(leftChild.keys)-1]
	parnet.keys = decreaseSpace(parnet.keys, idx)
	parnet.childs = decreaseSpace2(parnet.childs, idx)
	return true
}

func mergeRight(parnet *Node, idx int, maxKeys int) bool {
	if idx == len(parnet.keys)-1 {
		return false
	}
	child := parnet.childs[idx]
	rightChild := parnet.childs[idx+1]
	if len(rightChild.keys)+len(child.keys) > maxKeys {
		return false
	}
	child.keys = append(child.keys, rightChild.keys...)
	if child.isLeaf() {
		child.vals = append(child.vals, rightChild.vals...)
	} else {
		child.childs = append(child.childs, rightChild.childs...)
	}
	child.next = rightChild.next
	parnet.childs[idx+1] = child
	parnet.keys = decreaseSpace(parnet.keys, idx)
	parnet.childs = decreaseSpace2(parnet.childs, idx)
	return true
}

func (bp *BPlusTree) maxKeys() int {
	return bp.order - 1
}

func (bp *BPlusTree) minKeys() int {
	return (bp.order+1)/2 - 1
}

func increaseSpace(arr []interface{}, idx int) []interface{} {
	arr = append(arr, nil)
	if idx < len(arr)-1 {
		copy(arr[idx+1:], arr[idx:])
	}
	return arr
}

func increaseSpace2(arr []*Node, idx int) []*Node {
	arr = append(arr, nil)
	if idx < len(arr)-1 {
		copy(arr[idx+1:], arr[idx:])
	}
	return arr
}

func decreaseSpace(arr []interface{}, idx int) []interface{} {
	copy(arr[idx:], arr[idx+1:])
	arr = arr[:len(arr)-1]
	return arr
}

func decreaseSpace2(arr []*Node, idx int) []*Node {
	copy(arr[idx:], arr[idx+1:])
	arr = arr[:len(arr)-1]
	return arr
}
