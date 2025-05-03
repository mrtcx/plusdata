// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/mrtcx/plusdata/deque"
	"github.com/mrtcx/plusdata/deque/circularblocks"
)

func main() {
	var q deque.Deque = circularblocks.New()
	q.PushBack(3)       //[3]
	q.PushBack(4)       //[3,4]
	q.PushFront(2)      //[2,3,4]
	q.PushFront(1)      //[1,2,3,4]
	q.Size()            //4
	_ = q.Front().(int) //1
	_ = q.Back().(int)  //4

	for i := 0; i < q.Size(); i++ {
		_ = q.Get(i).(int) //i
	}
	for i := 0; i < q.Size(); i++ {
		q.Set(i, i*-1)
	}

	q.PopBack()  //[-1,-2,-3]
	q.PopFront() //[-2,-3]
	q.PopBack()  //[-2]
	q.PopFront() //[]
	q.Size()     //0
	q.Front()    //nil
	q.Back()     //nil

	q.PushBack(1) //[1]
	q.Clean()     //[]
}
