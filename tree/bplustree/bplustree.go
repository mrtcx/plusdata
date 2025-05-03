// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bplustree

import (
	"github.com/mrtcx/plusdata/tree"
)

var _ tree.Tree = (*bplusTree)(nil)

type bplusTree struct {
	root  *Node
	cmp   tree.Comparator
	size  int
	order int
}

type Node struct {
	keys   []interface{}
	valus  []interface{}
	childs []*Node
	next   *Node
	prev   *Node
}

func (node *Node) isLeaf() bool {
	return len(node.childs) == 0
}

func New(cmp tree.Comparator, order int) *bplusTree {
	return &bplusTree{
		cmp:   cmp,
		order: order,
	}
}

func (bp *bplusTree) Clean() {
	bp.root, bp.size = nil, 0
}

func (bp *bplusTree) Size() int {
	return bp.size
}

func (bp *bplusTree) Empty() bool {
	return bp.size == 0
}

func (bp *bplusTree) Insert(key, val interface{}) {
	root := bp.insert(bp.root, key, val)
	if len(root.keys) > int(bp.maxKeys()) {
		skey, sleft, sright := split(root)
		root = &Node{
			keys:   []interface{}{skey},
			childs: []*Node{sleft, sright},
		}
	}
	bp.root = root
}

func (bp *bplusTree) insert(root *Node, key, val interface{}) *Node {
	if root == nil {
		bp.size++
		return &Node{
			keys:  []interface{}{key},
			valus: []interface{}{val},
		}
	}
	idx := bp.binarySearchIdx(key, root)
	if root.isLeaf() {
		if idx < len(root.keys) && bp.cmp(root.keys[idx], key) == 0 {
			root.valus[idx] = val
			return root
		}
		bp.size++
		root.keys = increaseSpace(root.keys, idx)
		root.valus = increaseSpace(root.valus, idx)
		root.keys[idx], root.valus[idx] = key, val
		return root
	} else {
		retNode := bp.insert(root.childs[idx], key, val)
		if len(retNode.keys) > bp.maxKeys() {
			skey, sleft, sright := split(retNode)
			root.keys = increaseSpace(root.keys, idx)
			root.childs = increaseSpace2(root.childs, idx)
			root.keys[idx] = skey
			root.childs[idx], root.childs[idx+1] = sleft, sright
		}
		return root
	}
}

func (bp *bplusTree) Remove(key interface{}) {
	if bp.root == nil {
		return
	}
	root := bp.remove(bp.root, key)
	if root.isLeaf() {
		if len(root.keys) == 0 {
			bp.root = nil
		}
		return
	}
	if len(root.childs) < 2 {
		bp.root = root.childs[0]
		return
	}
}

func (bp *bplusTree) remove(root *Node, key interface{}) *Node {
	if root == nil {
		return nil
	}
	idx := bp.binarySearchIdx(key, root)
	if root.isLeaf() {
		if idx < len(root.keys) && bp.cmp(root.keys[idx], key) == 0 {
			bp.size--
			root.keys = decreaseSpace(root.keys, idx)
			root.valus = decreaseSpace(root.valus, idx)
		}
	} else {
		retNode := bp.remove(root.childs[idx], key)
		if len(retNode.keys) < bp.minKeys() {
			minkeys, maxkeys := bp.minKeys(), bp.maxKeys()
			if !borrowFromRight(root, idx, minkeys) &&
				!mergeRight(root, idx, maxkeys) &&
				!borrowFromLeft(root, idx, minkeys) &&
				!mergeLeft(root, idx, maxkeys) {
				panic("")
			}
		}
	}
	return root
}

func (bp *bplusTree) Get(key interface{}) (interface{}, bool) {
	node, idx := bp.getNodeIdx(key)
	if node == nil {
		return nil, false
	}
	return node.valus[idx], true
}

func (bp *bplusTree) getNodeIdx(key interface{}) (*Node, int) {
	root := bp.root
	if root == nil {
		return nil, 0
	}
	for {
		idx := bp.binarySearchIdx(key, root)
		if root.isLeaf() {
			if idx < len(root.keys) && bp.cmp(root.keys[idx], key) == 0 {
				return root, idx
			} else {
				return nil, 0
			}
		}
		root = root.childs[idx]
	}
}

func (bp *bplusTree) prevNodeIdx(key interface{}) (*Node, int) {
	root := bp.root
	if root == nil {
		return nil, 0
	}
	for {
		idx := bp.binarySearchIdx(key, root)
		if root.isLeaf() {
			if idx == len(root.keys) || bp.cmp(root.keys[idx], key) >= 0 {
				idx--
			}
			if idx >= 0 {
				return root, idx
			} else {
				if root.prev == nil {
					return nil, 0
				} else {
					return root.prev, len(root.prev.keys) - 1
				}
			}
		}
		root = root.childs[idx]
	}
}

func (bp *bplusTree) nextNodeIdx(key interface{}) (*Node, int) {
	root := bp.root
	if root == nil {
		return nil, 0
	}
	for {
		idx := bp.binarySearchIdx(key, root)
		if root.isLeaf() {
			if idx == len(root.keys) || bp.cmp(key, root.keys[idx]) >= 0 {
				idx++
			}
			if idx < len(root.keys) {
				return root, idx
			} else {
				if root.next == nil {
					return nil, 0
				} else {
					return root.next, 0
				}
			}
		}
		root = root.childs[idx]
	}
}

func (bp *bplusTree) binarySearchIdx(key interface{}, root *Node) int {
	if len(root.keys) <= 4 {
		idx := 0
		for ; idx < len(root.keys) && bp.cmp(key, root.keys[idx]) > 0; idx++ {
		}
		return idx
	}
	l, r := 0, len(root.keys)-1
	for l < r {
		mid := (l + r) / 2
		if bp.cmp(key, root.keys[mid]) > 0 {
			l = mid + 1
		} else {
			r = mid
		}
	}
	if bp.cmp(key, root.keys[r]) <= 0 {
		return r
	} else {
		return r + 1
	}
}

func (bp *bplusTree) mostLeft() *Node {
	if bp.root == nil {
		return nil
	}
	root := bp.root
	for !root.isLeaf() {
		root = root.childs[0]
	}
	return root
}

func (bp *bplusTree) mostRight() *Node {
	if bp.root == nil {
		return nil
	}
	root := bp.root
	for !root.isLeaf() {
		root = root.childs[len(root.childs)-1]
	}
	return root
}

func split(node *Node) (key interface{}, left *Node, right *Node) {
	pivot := len(node.keys) / 2
	right = &Node{}
	if node.isLeaf() {
		key = node.keys[pivot-1]
		right.keys = append(right.keys, node.keys[pivot:]...)
		node.keys = node.keys[:pivot]
		right.valus = append(right.valus, node.valus[pivot:]...)
		node.valus = node.valus[:pivot]
	} else {
		key = node.keys[pivot]
		right.keys = append(right.keys, node.keys[pivot+1:]...)
		node.keys = node.keys[:pivot]
		right.childs = append(right.childs, node.childs[pivot+1:]...)
		node.childs = node.childs[:pivot+1]
	}
	left = node
	if node.isLeaf() {
		if left.next != nil {
			left.next.prev = right
		}
		left.next, right.next, right.prev = right, left.next, left
	}
	return key, left, right
}

func mergeRight(parnet *Node, idx int, maxKeys int) bool {
	if idx == len(parnet.keys) {
		return false
	}
	idxChild := parnet.childs[idx]
	idxRightChild := parnet.childs[idx+1]
	if len(idxRightChild.keys)+len(idxChild.keys) > maxKeys {
		return false
	}
	if idxChild.isLeaf() {
		idxChild.keys = append(idxChild.keys, idxRightChild.keys...)
		idxChild.valus = append(idxChild.valus, idxRightChild.valus...)
		if idxRightChild.next != nil {
			idxRightChild.next.prev = idxChild
		}
		idxChild.next = idxRightChild.next
	} else {
		idxChild.keys = append(idxChild.keys, parnet.keys[idx])
		idxChild.keys = append(idxChild.keys, idxRightChild.keys...)
		idxChild.childs = append(idxChild.childs, idxRightChild.childs...)
	}
	parnet.childs[idx+1] = idxChild
	parnet.keys = decreaseSpace(parnet.keys, idx)
	parnet.childs = decreaseSpace2(parnet.childs, idx)
	return true
}

func mergeLeft(parnet *Node, idx int, maxKeys int) bool {
	if idx == 0 {
		return false
	}
	idxLeftChild := parnet.childs[idx-1]
	idxChild := parnet.childs[idx]
	if len(idxLeftChild.keys)+len(idxChild.keys) > maxKeys {
		return false
	}
	if idxChild.isLeaf() {
		idxLeftChild.keys = append(idxLeftChild.keys, idxChild.keys...)
		idxLeftChild.valus = append(idxLeftChild.valus, idxChild.valus...)
		if idxChild.next != nil {
			idxChild.next.prev = idxLeftChild
		}
		idxLeftChild.next = idxChild.next
	} else {
		idxLeftChild.keys = append(idxLeftChild.keys, parnet.keys[idx-1])
		idxLeftChild.keys = append(idxLeftChild.keys, idxChild.keys...)
		idxLeftChild.childs = append(idxLeftChild.childs, idxChild.childs...)
	}
	parnet.childs[idx] = idxLeftChild
	parnet.keys = decreaseSpace(parnet.keys, idx-1)
	parnet.childs = decreaseSpace2(parnet.childs, idx-1)
	return true
}

func borrowFromRight(parent *Node, idx int, minKeys int) bool {
	if idx == len(parent.keys) {
		return false
	}
	iChild := parent.childs[idx]
	iRightChild := parent.childs[idx+1]

	if len(iRightChild.keys) <= minKeys {
		return false
	}
	if iChild.isLeaf() {
		parent.keys[idx] = iRightChild.keys[0]
		iChild.keys = append(iChild.keys, iRightChild.keys[0])
		iChild.valus = append(iChild.valus, iRightChild.valus[0])
		iRightChild.keys = decreaseSpace(iRightChild.keys, 0)
		iRightChild.valus = decreaseSpace(iRightChild.valus, 0)
	} else {
		iChild.keys = append(iChild.keys, parent.keys[idx])
		iChild.childs = append(iChild.childs, iRightChild.childs[0])
		parent.keys[idx] = iRightChild.keys[0]
		iRightChild.keys = decreaseSpace(iRightChild.keys, 0)
		iRightChild.childs = decreaseSpace2(iRightChild.childs, 0)
	}
	return true
}

func borrowFromLeft(parent *Node, idx int, minKeys int) bool {
	if idx == 0 {
		return false
	}
	iLeftChild := parent.childs[idx-1]
	iChild := parent.childs[idx]
	if len(iLeftChild.keys) <= minKeys {
		return false
	}
	lli := len(iLeftChild.keys) - 1
	if iChild.isLeaf() {
		iChild.keys = increaseSpace(iChild.keys, 0)
		iChild.keys[0] = iLeftChild.keys[lli]
		iChild.valus = increaseSpace(iChild.valus, 0)
		iChild.valus[0] = iLeftChild.valus[lli]
		iLeftChild.keys = decreaseSpace(iLeftChild.keys, lli)
		iLeftChild.valus = decreaseSpace(iLeftChild.valus, lli)
		parent.keys[idx-1] = iLeftChild.keys[lli-1]
	} else {
		iChild.keys = increaseSpace(iChild.keys, 0)
		iChild.keys[0] = parent.keys[idx-1]
		iChild.childs = increaseSpace2(iChild.childs, 0)
		iChild.childs[0] = iLeftChild.childs[lli+1]
		parent.keys[idx-1] = iLeftChild.keys[lli]
		iLeftChild.keys = decreaseSpace(iLeftChild.keys, lli)
		iLeftChild.childs = decreaseSpace2(iLeftChild.childs, lli+1)
	}
	return true
}

func (bp *bplusTree) maxKeys() int {
	return bp.order - 1
}

func (bp *bplusTree) minKeys() int {
	return (bp.order+1)/2 - 1
}

func increaseSpace(arr []interface{}, idx int) []interface{} {
	arr = append(arr, nil)
	if idx != len(arr)-1 {
		copy(arr[idx+1:], arr[idx:])
	}
	return arr
}

func increaseSpace2(arr []*Node, idx int) []*Node {
	arr = append(arr, nil)
	if idx != len(arr)-1 {
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
