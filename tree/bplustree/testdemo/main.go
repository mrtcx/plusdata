// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import "fmt"

func main() {
	var ca []int = make([]int, 0, 4)
	fmt.Println(cap(ca), len(ca))
	var copya []int = make([]int, 0)
	fmt.Println(cap(copya), len(copya))
	copy(copya, ca)
	fmt.Println(cap(copya), len(copya))

	var a []int = []int{1}
	a = append(a, []int{}...)
	fmt.Println(len(a), a)
	a = append(a, 2)
	a = append(a, 3)
	copy(a[0:], a[1:])
	a[0] = 2
	fmt.Println(a)
	fmt.Println(len(a), a)

	a = []int{1, 2, 3, 4, 5}
	fmt.Println(cap(a), len(a))
	copy(a[2:], a[5:])
	fmt.Println(a)
}
