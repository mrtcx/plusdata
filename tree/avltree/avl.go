// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package avltree

import (
	"github.com/mrtcx/plusdata/tree"
)

var _ tree.Tree = (*avlTree)(nil)

type avlTree struct {
	root *Node
	cmp  tree.Comparator
	size int
}

type Node struct {
	key    interface{}
	value  interface{}
	lchild *Node
	rchild *Node
	h      int8
}

// 加一个nil的虚拟节点，左旋和右旋可以少很多特判
var _nil *Node = &Node{h: 1}

func newNode(key, val interface{}) *Node {
	p := &Node{key: key, value: val, h: 1, rchild: _nil, lchild: _nil}
	return p
}

func New(cmp tree.Comparator) *avlTree {
	return &avlTree{
		cmp:  cmp,
		root: _nil,
	}
}

func (avl *avlTree) Clean() {
	avl.size, avl.root = 0, _nil
}

func (avl *avlTree) Size() int {
	return avl.size
}

func (avl *avlTree) Empty() bool {
	return avl.size == 0
}

func (avl *avlTree) Insert(key, val interface{}) {
	avl.root = avl.insert(avl.root, key, val)
}

func (avl *avlTree) insert(root *Node, key interface{}, val interface{}) *Node {
	if root == _nil {
		avl.size++
		return newNode(key, val)
	}
	less := avl.cmp(key, root.key)
	if less == 0 {
		root.value = val
		return root
	} else if less < 0 {
		root.lchild = avl.insert(root.lchild, key, val)
	} else {
		root.rchild = avl.insert(root.rchild, key, val)
	}
	newroot := maintain(root)
	updateHeight(root)
	return newroot
}

func (avl *avlTree) Remove(key interface{}) {
	avl.root = avl.remove(avl.root, key)
}

func (avl *avlTree) remove(root *Node, key interface{}) *Node {
	if root == _nil {
		return _nil
	}
	less := avl.cmp(key, root.key)
	if less < 0 {
		root.lchild = avl.remove(root.lchild, key)
	} else if less > 0 {
		root.rchild = avl.remove(root.rchild, key)
	} else {
		if root.lchild == _nil || root.rchild == _nil {
			temp := root.lchild
			if temp == _nil {
				temp = root.rchild
			}
			avl.size--
			return temp
		} else {
			temp := precusor(root)
			root.key = temp.key
			root.value = temp.value
			root.lchild = avl.remove(root.lchild, temp.key)
		}
	}
	newroot := maintain(root)
	updateHeight(root)
	return newroot
}

func (avl *avlTree) Get(key interface{}) (interface{}, bool) {
	node := avl.findNode(key)
	if node == nil {
		return nil, false
	}
	return node.value, true
}

func (avl *avlTree) findNode(key interface{}) *Node {
	root := avl.root
	for {
		if root == _nil {
			return nil
		}
		less := avl.cmp(key, root.key)
		if less == 0 {
			return root
		} else if less < 0 {
			root = root.lchild
		} else {
			root = root.rchild
		}
	}
}

func (avl *avlTree) leftNode() *Node {
	root := avl.root
	if root == _nil {
		return nil
	}
	for root.lchild != _nil {
		root = root.lchild
	}
	return root
}

func (avl *avlTree) rightNode() *Node {
	root := avl.root
	if root == _nil {
		return nil
	}
	for root.rchild != _nil {
		root = root.rchild
	}
	return root
}

func (avl *avlTree) findPrevNode(key interface{}) *Node {
	root := avl.root
	var lparent *Node = nil
	for {
		if root == _nil {
			return nil
		}
		less := avl.cmp(key, root.key)
		if less > 0 {
			if root.rchild == _nil {
				return root
			}
			root, lparent = root.rchild, root
		} else {
			if root.lchild == _nil {
				return lparent
			}
			root = root.lchild
		}
	}
}

func (avl *avlTree) findNextNode(key interface{}) *Node {
	root := avl.root
	var rparent *Node = nil
	for {
		if root == _nil {
			return nil
		}
		less := avl.cmp(key, root.key)
		if less >= 0 {
			if root.rchild == _nil {
				return rparent
			}
			root = root.rchild
		} else {
			if root.lchild == _nil {
				return root
			}
			root, rparent = root.lchild, root
		}
	}
}

func maintain(root *Node) *Node {
	diff := root.lchild.h - root.rchild.h
	if diff <= 1 && diff >= -1 {
		return root
	}
	if root.lchild.h > root.rchild.h {
		if root.lchild.lchild.h < root.lchild.rchild.h {
			root.lchild = leftRotate(root.lchild)
		}
		root = rightRotate(root)
	} else {
		if root.rchild.rchild.h < root.rchild.lchild.h {
			root.rchild = rightRotate(root.rchild)
		}
		root = leftRotate(root)
	}
	return root
}

func precusor(root *Node) *Node {
	temp := root.lchild
	for temp.rchild != _nil {
		temp = temp.rchild
	}
	return temp
}

func updateHeight(root *Node) {
	root.h = root.lchild.h + 1
	if root.rchild.h >= root.h {
		root.h = root.rchild.h + 1
	}
}

func leftRotate(root *Node) *Node {
	temp := root.rchild
	root.rchild = temp.lchild
	temp.lchild = root
	updateHeight(root)
	updateHeight(temp)
	return temp
}

func rightRotate(root *Node) *Node {
	temp := root.lchild
	root.lchild = temp.rchild
	temp.rchild = root
	updateHeight(root)
	updateHeight(temp)
	return temp
}
