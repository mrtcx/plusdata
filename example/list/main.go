// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/mrtcx/plusdata/list/singlelinkedlist"
)

func main() {
	list := singlelinkedlist.New()
	list.PushBack(2)  //2
	list.PushBack(3)  //2->3
	list.PushBack(4)  //2->3->4
	list.PushFront(0) //0->2->3->4

	list.Back()  // e4
	list.Front() //e1
	list.Size()  //4

	list.InsertAfter(list.Front(), 1) //0->1->2->3->4

	e := list.Front()
	for e != nil {
		fmt.Println(e.Value) //0 1 2 3 4
		e = e.Next()
	}

	anthor := singlelinkedlist.New()
	anthor.PushFront(-1) //-1
	anthor.Merge(list)   //-1->0->1->2->3->4

	list.Clean()
	list.Size() //0
}
