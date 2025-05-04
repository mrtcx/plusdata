// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/mrtcx/plusdata/tree"
	"github.com/mrtcx/plusdata/tree/skiplist"
)

func main() {
	var tree tree.Tree = skiplist.New(tree.IntComparator)
	//添加
	tree.Insert(1, 1)
	tree.Insert(2, 2)
	tree.Insert(3, 3)
	tree.Insert(1, -1)
	//查找
	e1 := tree.Find(1)
	v1, exist1 := tree.Get(1)
	fmt.Println(v1, exist1)           //-1, true
	fmt.Println(e1.Key(), e1.Value()) //1, -1
	//查找前序和后继（即使key不在树中）
	tree.Prev(2) //1
	tree.Next(2) //3
	tree.Prev(3) //2
	tree.Next(3) //nil
	tree.Prev(1) //nil
	//正序遍历
	e := tree.Left()
	for e != nil {
		key, value := e.Key(), e.Value()
		fmt.Println(key, value)
		e = e.Next()
	}
	//倒序遍历
	e = tree.Right()
	for e != nil {
		key, value := e.Key(), e.Value()
		fmt.Println(key, value)
		e = e.Prev()
	}
	//区间
	e = tree.Find(2)
	for l, r := e.Prev(), e.Next(); l != nil && r != nil; l, r = l.Prev(), r.Next() {
		lv := l.Value()
		rv := r.Value()
		l.SetValue(rv)
		r.SetValue(lv)
	}
	//删除
	tree.Remove(1)
	tree.Remove(2)
	tree.Remove(3)
	tree.Find(1) //nil
	tree.Size()  //0
}
