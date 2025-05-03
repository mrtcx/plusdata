// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package btree

import (
	"github.com/mrtcx/plusdata/tree"
)

type bTree struct {
	root  *Node
	cmp   tree.Comparator
	size  int
	order int
}

type Node struct {
	keys   []interface{}
	values []interface{}
	childs []*Node
}

func (node *Node) isLeaf() bool {
	return len(node.childs) == 0
}

func New(cmp tree.Comparator, order int) *bTree {
	return &bTree{
		cmp:   cmp,
		order: order,
	}
}

func (bp *bTree) Clean() {
	bp.root, bp.size = nil, 0
}

func (bp *bTree) Size() int {
	return bp.size
}

func (bp *bTree) Empty() bool {
	return bp.size == 0
}

func (bp *bTree) Insert(key, val interface{}) {
	root := bp.insert(bp.root, key, val)
	if len(root.keys) > int(bp.maxKeys()) {
		skey, sval, sleft, sright := split(root)
		root = &Node{
			keys:   []interface{}{skey},
			values: []interface{}{sval},
			childs: []*Node{sleft, sright},
		}
	}
	bp.root = root
}

func (bp *bTree) insert(root *Node, key, val interface{}) *Node {
	if root == nil {
		bp.size++
		return &Node{
			keys:   []interface{}{key},
			values: []interface{}{val},
		}
	}
	idx := bp.binarySearchIdx(key, root)
	if idx < len(root.keys) && bp.cmp(root.keys[idx], key) == 0 {
		root.values[idx] = val
		return root
	}
	if root.isLeaf() {
		bp.size++
		root.keys = increaseSpace(root.keys, idx)
		root.values = increaseSpace(root.values, idx)
		root.keys[idx], root.values[idx] = key, val
		return root
	} else {
		retNode := bp.insert(root.childs[idx], key, val)
		if len(retNode.keys) > bp.maxKeys() {
			skey, sval, sleft, sright := split(retNode)
			root.keys = increaseSpace(root.keys, idx)
			root.values = increaseSpace(root.values, idx)
			root.childs = increaseSpace2(root.childs, idx)
			root.keys[idx] = skey
			root.values[idx] = sval
			root.childs[idx], root.childs[idx+1] = sleft, sright
		}
		return root
	}
}

func (bp *bTree) Remove(key interface{}) {
	if bp.root == nil {
		return
	}
	root := bp.remove(bp.root, key)
	if len(root.keys) == 0 {
		if root.isLeaf() {
			bp.root = nil
		} else {
			bp.root = bp.root.childs[0]
		}
	}
}

func (bp *bTree) remove(root *Node, key interface{}) *Node {
	if root == nil {
		return nil
	}
	idx := bp.binarySearchIdx(key, root)
	if idx < len(root.keys) && bp.cmp(root.keys[idx], key) == 0 {
		if root.isLeaf() {
			bp.size--
			root.keys = decreaseSpace(root.keys, idx)
			root.values = decreaseSpace(root.values, idx)
			return root
		} else {
			mright := mostRight(root.childs[idx])
			rkey, rval := mright.keys[len(mright.keys)-1], mright.values[len(mright.values)-1]
			root.keys[idx], root.values[idx] = rkey, rval
			key = rkey
		}
	}
	if !root.isLeaf() {
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

func (bp *bTree) Get(key interface{}) (interface{}, bool) {
	node, idx := bp.findNodeIdx(key)
	if node == nil {
		return nil, false
	}
	return node.values[idx], true
}

func (bp *bTree) findNodeIdx(key interface{}) (*Node, int) {
	root := bp.root
	if root == nil {
		return nil, 0
	}
	for {
		idx := bp.binarySearchIdx(key, root)
		if idx < len(root.keys) && bp.cmp(root.keys[idx], key) == 0 {
			return root, idx
		}
		if root.isLeaf() {
			return nil, 0
		}
		root = root.childs[idx]
	}
}

func (bp *bTree) findPrevNodeIdx(key interface{}) (*Node, int) {
	root := bp.root
	if root == nil {
		return nil, 0
	}
	var lparent *Node
	var lparentIdx int
	for {
		idx := bp.binarySearchIdx(key, root)
		if idx < len(root.keys) && bp.cmp(root.keys[idx], key) == 0 && !root.isLeaf() {
			mright := mostRight(root.childs[idx])
			return mright, len(mright.keys) - 1
		}
		if root.isLeaf() {
			if idx-1 >= 0 {
				return root, idx - 1
			}
			if lparent == nil {
				return nil, 0
			}
			return lparent, lparentIdx
		}
		if idx > 0 {
			lparent, lparentIdx = root, idx-1
		}
		root = root.childs[idx]
	}
}

func (bp *bTree) findNextNodeIdx(key interface{}) (*Node, int) {
	root := bp.root
	if root == nil {
		return nil, 0
	}
	var rparent *Node
	var rparentIdx int
	for {
		idx := bp.binarySearchIdx(key, root)
		if idx < len(root.keys) && bp.cmp(key, root.keys[idx]) == 0 {
			if !root.isLeaf() {
				return mostLeft(root.childs[idx+1]), 0
			}
			idx++
		}
		if root.isLeaf() {
			if idx < len(root.keys) {
				return root, idx
			}
			if rparent == nil {
				return nil, 0
			}
			return rparent, rparentIdx
		}
		if idx < len(root.keys) {
			rparent, rparentIdx = root, idx
		}
		root = root.childs[idx]
	}
}

func (bp *bTree) binarySearchIdx(key interface{}, root *Node) int {
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

func split(node *Node) (key interface{}, val interface{}, left *Node, right *Node) {
	pivot := len(node.keys) / 2
	key = node.keys[pivot]
	val = node.values[pivot]
	right = &Node{}
	right.keys = append(right.keys, node.keys[pivot+1:]...)
	right.values = append(right.values, node.values[pivot+1:]...)
	node.keys = node.keys[:pivot]
	node.values = node.values[:pivot]
	if !node.isLeaf() {
		right.childs = append(right.childs, node.childs[pivot+1:]...)
		node.childs = node.childs[:pivot+1]
	}
	left = node
	return key, val, left, right
}

func mergeRight(parnet *Node, idx int, maxKeys int) bool {
	if idx == len(parnet.keys) {
		return false
	}
	idxChild := parnet.childs[idx]
	idxRightChild := parnet.childs[idx+1]
	if len(idxRightChild.keys)+len(idxChild.keys)+1 > maxKeys {
		return false
	}
	idxChild.keys = append(idxChild.keys, parnet.keys[idx])
	idxChild.values = append(idxChild.values, parnet.values[idx])
	idxChild.keys = append(idxChild.keys, idxRightChild.keys...)
	idxChild.values = append(idxChild.values, idxRightChild.values...)
	if !idxChild.isLeaf() {
		idxChild.childs = append(idxChild.childs, idxRightChild.childs...)
	}
	parnet.childs[idx+1] = idxChild
	parnet.keys = decreaseSpace(parnet.keys, idx)
	parnet.values = decreaseSpace(parnet.values, idx)
	parnet.childs = decreaseSpace2(parnet.childs, idx)
	return true
}

func mergeLeft(parnet *Node, idx int, maxKeys int) bool {
	if idx == 0 {
		return false
	}
	idxLeftChild := parnet.childs[idx-1]
	idxChild := parnet.childs[idx]
	if len(idxLeftChild.keys)+len(idxChild.keys)+1 > maxKeys {
		return false
	}
	idxLeftChild.keys = append(idxLeftChild.keys, parnet.keys[idx-1])
	idxLeftChild.values = append(idxLeftChild.values, parnet.values[idx-1])
	idxLeftChild.keys = append(idxLeftChild.keys, idxChild.keys...)
	idxLeftChild.values = append(idxLeftChild.values, idxChild.values...)
	if !idxChild.isLeaf() {
		idxLeftChild.childs = append(idxLeftChild.childs, idxChild.childs...)
	}
	parnet.childs[idx] = idxLeftChild
	parnet.keys = decreaseSpace(parnet.keys, idx-1)
	parnet.values = decreaseSpace(parnet.values, idx-1)
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
	iChild.keys = append(iChild.keys, parent.keys[idx])
	iChild.values = append(iChild.values, parent.values[idx])
	parent.keys[idx] = iRightChild.keys[0]
	parent.values[idx] = iRightChild.values[0]
	iRightChild.keys = decreaseSpace(iRightChild.keys, 0)
	iRightChild.values = decreaseSpace(iRightChild.values, 0)
	if !iChild.isLeaf() {
		iChild.childs = append(iChild.childs, iRightChild.childs[0])
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
	iChild.keys = increaseSpace(iChild.keys, 0)
	iChild.values = increaseSpace(iChild.values, 0)
	iChild.keys[0] = parent.keys[idx-1]
	iChild.values[0] = parent.values[idx-1]
	parent.keys[idx-1] = iLeftChild.keys[lli]
	parent.values[idx-1] = iLeftChild.values[lli]
	iLeftChild.keys = decreaseSpace(iLeftChild.keys, lli)
	iLeftChild.values = decreaseSpace(iLeftChild.values, lli)
	if !iChild.isLeaf() {
		iChild.childs = increaseSpace2(iChild.childs, 0)
		iChild.childs[0] = iLeftChild.childs[lli+1]
		iLeftChild.childs = decreaseSpace2(iLeftChild.childs, lli+1)
	}
	return true
}

func (bp *bTree) maxKeys() int {
	return bp.order
}

func (bp *bTree) minKeys() int {
	return (bp.order+1)/2 - 1
}

func mostRight(root *Node) *Node {
	for !root.isLeaf() {
		root = root.childs[len(root.childs)-1]
	}
	return root
}

func mostLeft(root *Node) *Node {
	for !root.isLeaf() {
		root = root.childs[0]
	}
	return root
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
