// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/mrtcx/plusdata/arrary"
	"github.com/mrtcx/plusdata/arrary/blockslices"
)

func main() {
	var arr arrary.Arrary = blockslices.New()
	arr.PushBack(1)       //[1]
	arr.PushBack(2)       //[1,2]
	arr.Size()            //2
	_ = arr.Front().(int) //1
	_ = arr.Back().(int)  //2

	for i := 0; i < arr.Size(); i++ {
		_ = arr.Get(i).(int) //i
	}
	for i := 0; i < arr.Size(); i++ {
		arr.Set(i, i*-1)
	}

	arr.PopBack() //[-1]
	arr.PopBack() //[]
	arr.Size()    //0
	arr.Front()   //nil
	arr.Back()    //nil

	arr.PushBack(1) //[1]
	arr.Clean()     //[]
}
