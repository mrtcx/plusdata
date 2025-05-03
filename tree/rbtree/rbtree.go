// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package rbtree

import (
	"github.com/mrtcx/plusdata/tree"
)

var _ tree.Tree = (*rbTree)(nil)

type rbTree struct {
	root *Node
	cmp  tree.Comparator
	size int
}

type Node struct {
	key            interface{}
	value          interface{}
	color          rbColor //0- 红色， 1-黑色 2-双重黑
	lchild, rchild *Node
}

type rbColor uint8

const (
	red         rbColor = 0
	black       rbColor = 1
	doubleBlack rbColor = 2
)

// 虚拟的NIL叶子节点(黑色),所有nil都用它来替代，代码会少很多特判
var _nil *Node = &Node{
	color: black,
}

func newNode(k, v interface{}) *Node {
	t := Node{
		key:    k,
		value:  v,
		color:  red,
		lchild: _nil,
		rchild: _nil,
	}
	return &t
}

func New(cmp tree.Comparator) *rbTree {
	return &rbTree{
		root: _nil,
		cmp:  cmp,
	}
}

func (rb *rbTree) Clean() {
	rb.root, rb.size = _nil, 0
}

func (rb *rbTree) Size() int {
	return rb.size
}

func (rb *rbTree) Empty() bool {
	return rb.size == 0
}

func (rb *rbTree) Insert(key, value interface{}) {
	rb.root = rb.insert(rb.root, key, value)
	rb.root.color = black
}

func (rb *rbTree) insert(root *Node, key, value interface{}) *Node {
	if root == _nil {
		rb.size++
		return newNode(key, value)
	}
	less := rb.cmp(key, root.key)
	if less == 0 {
		root.value = value
		return root
	} else if less < 0 {
		root.lchild = rb.insert(root.lchild, key, value)
	} else {
		root.rchild = rb.insert(root.rchild, key, value)
	}
	return insertMaintain(root)
}

// 修复双红冲突
func insertMaintain(root *Node) *Node {
	if root.lchild.color == black && root.rchild.color == black {
		return root
	}
	hasRedChild := func(root *Node) bool {
		return root.lchild.color == red || root.rchild.color == 0
	}
	if root.lchild.color == red && root.rchild.color == red {
		root.color, root.lchild.color, root.rchild.color = red, black, black
	} else if root.lchild.color == red && hasRedChild(root.lchild) {
		if root.lchild.rchild.color == red {
			root.lchild = leftRota(root.lchild)
		}
		root = rightRota(root)
		root.lchild.color = black
	} else if root.rchild.color == red && hasRedChild(root.rchild) {
		if root.rchild.lchild.color == red {
			root.rchild = rightRota(root.rchild)
		}
		root = leftRota(root)
		root.rchild.color = black
	}
	return root
}

func (rb *rbTree) Remove(key interface{}) {
	rb.root = rb.remove(rb.root, key)
	rb.root.color = black
}

func (rb *rbTree) remove(root *Node, key interface{}) *Node {
	if root == _nil {
		return _nil
	}
	less := rb.cmp(key, root.key)
	if less == 0 {
		if root.lchild == _nil || root.rchild == _nil {
			tmp := root.lchild
			if tmp == _nil {
				tmp = root.rchild
			}
			tmp.color += root.color //此时会有"双黑"冲突
			rb.size--
			return tmp
		} else {
			tmp := precursor(root)
			root.key = tmp.key
			root.value = tmp.value
			root.lchild = rb.remove(root.lchild, tmp.key)
		}
	} else if less < 0 {
		root.lchild = rb.remove(root.lchild, key)
	} else {
		root.rchild = rb.remove(root.rchild, key)
	}
	return removeMaintain(root)
}

func (rb *rbTree) Get(key interface{}) (interface{}, bool) {
	node := rb.findNode(key)
	if node == nil {
		return nil, false
	}
	return node.value, true
}

func (rb *rbTree) leftNode() *Node {
	root := rb.root
	if root == _nil {
		return nil
	}
	for root.lchild != _nil {
		root = root.lchild
	}
	return root
}

func (rb *rbTree) rightNode() *Node {
	root := rb.root
	if root == _nil {
		return nil
	}
	for root.rchild != _nil {
		root = root.rchild
	}
	return root
}

func (rb *rbTree) findNode(key interface{}) *Node {
	root := rb.root
	for {
		if root == _nil {
			return nil
		}
		less := rb.cmp(key, root.key)
		if less == 0 {
			return root
		} else if less < 0 {
			root = root.lchild
		} else {
			root = root.rchild
		}
	}
}

func (rb *rbTree) findPrevNode(key interface{}) *Node {
	root := rb.root
	var lparent *Node = nil
	for {
		if root == _nil {
			return nil
		}
		less := rb.cmp(key, root.key)
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

func (rb *rbTree) findNextNode(key interface{}) *Node {
	root := rb.root
	var rparent *Node = nil
	for {
		if root == _nil {
			return nil
		}
		less := rb.cmp(key, root.key)
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

// 修复双黑冲突
func removeMaintain(root *Node) *Node {
	//无双黑冲突
	if root.lchild.color != doubleBlack && root.rchild.color != doubleBlack {
		return root
	}
	hasRedChild := func(inode *Node) bool {
		return inode.lchild.color == red || inode.rchild.color == red
	}

	//无红色孩子,也无红色孙子
	if root.lchild.color == doubleBlack && root.rchild.color != red && !hasRedChild(root.rchild) ||
		root.rchild.color == doubleBlack && root.lchild.color != red && !hasRedChild(root.lchild) {
		root.lchild.color -= black
		root.rchild.color -= black
		root.color += black
		return root
	}

	//有红色的孩子
	if hasRedChild(root) {
		if root.lchild.color == red {
			root = rightRota(root)
			root.color = black
			root.rchild.color = red
			root.rchild = removeMaintain(root.rchild)
		} else {
			root = leftRota(root)
			root.color = black
			root.lchild.color = red
			root.lchild = removeMaintain(root.lchild)
		}
		return root
	}

	//有红色的孙子
	if root.lchild.color == doubleBlack {
		root.lchild.color = black
		if root.rchild.rchild.color != red {
			root.rchild.color = red
			root.rchild.lchild.color = black
			root.rchild = rightRota(root.rchild)
		}
		root = leftRota(root)
		root.color = root.lchild.color
		root.lchild.color, root.rchild.color = black, black
		return root
	} else {
		root.rchild.color = black
		if root.lchild.lchild.color != red {
			root.lchild.color = red
			root.lchild.rchild.color = black
			root.lchild = leftRota(root.lchild)
		}
		root = rightRota(root)
		root.color = root.rchild.color
		root.lchild.color, root.rchild.color = black, black
		return root
	}
}

func precursor(root *Node) *Node {
	root = root.lchild
	for root.rchild != _nil {
		root = root.rchild
	}
	return root
}

func leftRota(root *Node) *Node {
	temp := root.rchild
	root.rchild = temp.lchild
	temp.lchild = root
	return temp
}

func rightRota(root *Node) *Node {
	temp := root.lchild
	root.lchild = temp.rchild
	temp.rchild = root
	return temp
}
