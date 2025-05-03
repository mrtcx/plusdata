# plusdata

plusdata补充golang常用的数据结构, 包括如下:
- [数组](#数组)
    - 分块切片
- [双端队列](#双端队列)
    - 分块循环buffer
- [堆](#堆)
    - 数组堆
    - 树形堆
- [链表](#链表)
    - 单链表
- [树](#树)
    - [跳表](#跳表) (这里把跳表归为一种树)
	- [avl树](#avl树)
    - [红黑树](#红黑树)
    - [b树](#红黑树)
    - [b+树](#b+树)

plusdata提供简洁的操作，每个数据结构只提供适合他的操作。

### 数组

分块切片，外层对应一个块级别的切片，块内对应一个元素的切片，块数和块内元素在容量使用率低于1/4时触发缩容，容量使用满后触发扩容。大容量扩缩容时涉及到拷贝消耗，因此使用了分块来避免。
```golang
┌───────────────────┐┌───────────────────┐┌───────────────────┐
│ slice0            ││ slice1            ││ slice2            │
│ ┌──┬──┬──┬──┐     ││ ┌──┬──┬──┬──┐     ││ ┌──┬──┬──┬──┐     │
│ │0 │1 │2 │3 │     ││ │4 │5 │6 │7 │     ││ │8 │9 │10│11│     │ ....
│ └──┴──┴──┴──┘     ││ └──┴──┴──┴──┘     ││ └──┴──┴──┴──┘     │
└────────┬──────────┘└────────┬──────────┘└────────┬──────────┘
        块索引0                 块索引1                  块索引2   
```

|对比 |按索引访问 |扩缩容 |对存储要求|
|:----|:------------|:--------|-------:|
| 不分块 | O(1) | 有大容量拷贝消耗 | 对线性地址有要求，尽管是虚拟地址，但也不是免费的|
| 分块 | O(1), 略低于不分块，涉及到取于运算 | 没有大容量拷贝消耗 | 对存储地址要求低|

对比原生slice测试, 添加元素:
```golang
添加元素:
BenchmarkPushBack/slice[10w]                 448           3059057 ns/op         9266362 B/op      42885 allocs/op
BenchmarkPushBack/blockslice[10w]            400           2595552 ns/op         5734315 B/op     100290 allocs/op

BenchmarkPushBack/slice[50w]                  60          18978696 ns/op        44714352 B/op         35 allocs/op
BenchmarkPushBack/blockslice[50w]             75          15277259 ns/op        28586926 B/op     502444 allocs/op

BenchmarkPushBack/slice[100w]                 24          50157507 ns/op        88017264 B/op         38 allocs/op
BenchmarkPushBack/blockslice[100w]            34          39477684 ns/op        57176750 B/op    1005130 allocs/op

BenchmarkPushBack/slice[200w]                 10         108629888 ns/op        172673392 B/op        41 allocs/op
BenchmarkPushBack/blockslice[200w]            18          64358171 ns/op        114372910 B/op   2010501 allocs/op
```
对比原生slice测试，按索引访问:
```golang
BenchmarkGet/slice[10w]                    35071             36641 ns/op               0 B/op          0 allocs/op
BenchmarkGet/blockslice[10w]               10000            112515 ns/op               0 B/op          0 allocs/op

BenchmarkGet/slice[50w]                     7168            172188 ns/op               0 B/op          0 allocs/op
BenchmarkGet/blockslice[50w]                2216            546468 ns/op               0 B/op          0 allocs/op

BenchmarkGet/slice[100w]                    3297            324049 ns/op               0 B/op          0 allocs/op
BenchmarkGet/blockslice[100w]               1059           1063612 ns/op               0 B/op          0 allocs/op

BenchmarkGet/slice[200w]                    1904            644378 ns/op               0 B/op          0 allocs/op
BenchmarkGet/blockslice[200w]                577           2074720 ns/op               0 B/op          0 allocs/op
```

**相关操作：**
```golang
type Arrary interface {
	Size() int
	Empty() bool
	Clean()
	Get(index int) interface{} //按索引访问
	Set(index int, value interface{}) //按索引赋值
	Front() interface{}
	Back() interface{}
	PushBack(interface{})
	PopBack() interface{}
}
```

**复杂度:**
> 除非另外说明，复杂度均为O(1)

|操作 |复杂度 |
|:-------|---------:|
|PushBack() | O(1), 扩容时触发拷贝，分块切片只会做小容量的拷贝|
|PopBack() |O(1), 缩容时触发拷贝，分块切片只会做小容量的拷贝|

**使用示例:**
```golang
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

	for i := 0; i < arr.Size(); i++ { //按索引访问
		_ = arr.Get(i).(int) //i
	}
	for i := 0; i < arr.Size(); i++ { //按索引赋值
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
```

### 双端队列

双端队列使用分块循环buffer实现，外层对应一个块级别的循环buffer，块内对应一个元素的循环buff，块数和块内元素在容量使用率低于1/4时触发缩容，容量使用满后扩容一倍，使用分块为了避免单一循环buff在大容量拷贝时带来的损耗。
```golang
          块头                                                              块尾
            ▼          →                  →                     →            ▼
            ╔═════════════════════╦═════════════════════╦═════════════════════
            ║ buffer1            ║║    buffer2         ║║    buffer3         ║         
            ║↑                 ↓ ║║↑                 ↓ ║║↑                  ↓║
↑      .....╠════╦═══╦═══╦═══╦═══╣╠═══╦═══╦═══╦═══╦═══╦╣╠═══╦═══╦═══╦═══╦═══╦╣.....  ↓
            ║   →   →   →   →   →│║   →   →   →   →   →│║   →   →   →   →   →│
            ║↑                  ║║↑                 ↓ ║║↑                  ↓ ║
            ╚════════════════════╝╚════════════════════╝╚════════════════════╝
																	       
```

**相关操作：**
```golang
type Arrary interface {
	Size() int
	Empty() bool
	Clean()
	Get(index int) interface{}
	Set(index int, value interface{})
	Front() interface{}
	Back() interface{}
	PushBack(interface{})
	PopBack() interface{}
    PushFront(value interface{})
	PopFront() interface{}
}
```

**复杂度:**
> 除非另外说明，复杂度均为O(1)

|操作 |复杂度 |
|:-------|---------:|
|PushBack() | O(1), 扩容时触发拷贝，分块只会做小容量的拷贝|
|PopBack() |O(1), 缩容时触发拷贝，分块只会做小容量的拷贝|
|PushFront() | O(1), 扩容时触发拷贝，分块只会做小容量的拷贝|
|PopFront() |O(1), 缩容时触发拷贝，分块只会做小容量的拷贝|

**使用示例:**
```golang
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
```
### 堆

提供了两种堆，一种底层存储是[数组](#数组)， 另一种底层存储是[树](#树)。 数组堆内元素没有去重，排序树堆内元素是去重的, 两者的复杂度都是log(N)，需要去重选择排序树堆，不需要去重选择数组堆， 数组堆存储使用更少、性能更高。
```golang
                              +----+ 
                                0 
                            /         \
                    +---+             +---+
                      2                    1 
                   /      \             /      \
               +--+       +--+       +--+       +--+
                6           7         3           4 
             /    \       /   \     /    \      /    \
           +-+     +-+  +-+    +-+  +-+   +-+   +-+   +-+
            9        8   10     12   5     11    13    14
```

**相关操作：**
```golang
type Heap interface {
	Size() int
	Empty() bool
	Clean()
	Top() interface{}
	Push(interface{})
	Pop() interface{}
}
```

**复杂度:**
> 除非另外说明，复杂度均为O(1)

|操作 |复杂度 |
|:-------|---------:|
|Push() | O(logN) |
|Pop() |O(logN)|

**使用示例:**
```golang
package main

import (
	"github.com/mrtcx/plusdata/heap"
	"github.com/mrtcx/plusdata/heap/arraryheap"
	"github.com/mrtcx/plusdata/heap/treeheap"
	"github.com/mrtcx/plusdata/tree"
	"github.com/mrtcx/plusdata/tree/rbtree"
)

func main() {
	//数组堆，堆内元素不会去重
	var heap heap.Heap = arraryheap.New(heap.IntLess)
	heap.Push(4)
	heap.Push(2)
	heap.Push(2)
	heap.Push(3)
	heap.Push(1)
	heap.Size() //5
	for !heap.Empty() {
		_ = heap.Top().(int) //1, 2, 2, 3, 4, 数组堆内元素没有去重
		heap.Pop()
	}

	//排序树堆
	heap = treeheap.New(rbtree.New(tree.IntComparator)) //这里使用红黑树
	heap.Push(4)
	heap.Push(2)
	heap.Push(2)
	heap.Push(3)
	heap.Push(1)
	heap.Size() //4
	for !heap.Empty() {
		heap.Top() //1, 2, 3, 4, 排序树堆内元素会去重
		heap.Pop()
	}

	heap.Push(1) //[1]
	heap.Clean() //[]
	heap.Size()  //0
}
```

### 链表

单链表
```golang
Front                                                                 Back
│                                                                     │
▼                                                                     ▼
┌─────────┬───────┐   ┌─────────┬───────┐   ┌─────────┬───────┐  ┌─────────┬───────┐     
│ 数据: 1 │ Next:•├─→ │ 数据: 2 │ Next:•├─→ │ 数据: 3 │ Next:•├─→ │ 数据: 4 │ Next:•├─→ NULL
└─────────┴───────┘   └─────────┴───────┘   └─────────┴───────┘  └─────────┴───────┘
```

**相关操作：**
```golang
func (*List) Size() int
func (*List) Empty() bool
func (*List) Clean()
func (*List) Front() *Element
func (*List) Back() *Element
func (*List) PushBack(interface{}) *Element
func (*List) PopBack() interface{}
func (*List) PushFront(value interface{}) *Element
func (*List) PopFront() interface{}
func (*List) Merge(delete *List)
func (*Element) Next() *Element
```

**复杂度**
>均为O(1)

**使用示例**
```golang
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
```

### 树

均为排序平衡树，把跳表归为一种树，因为他符合两个特点：排序、平衡。跳表根据概率分布实现了一种平衡，他的用法也契合排序树的用法。

**树相关操作：**
```golang
type Tree interface {
	Size() int
	Empty() bool
	Clean()
	Insert(key interface{}, value interface{}) //添加元素，key重复添加，value为最新值
	Remove(key interface{})  //删除元素
	Get(key interface{}) (interface{}, bool) //获取元素值
	Find(key interface{}) Element //获取元素
	Left() Element   //顺序遍历中，最左端的元素
	Right() Element  //顺序遍历中，最右端的元素
	Prev(key interface{}) Element  //查询key的前驱元素，即使key不存在
	Next(key interface{}) Element  //查询key的后继元素，即使key不存在
}

type Element interface { //元素访问，前序和后序遍历，(注意：不要在遍历的过程中，对树做添加和删除操作)
	Key() interface{}
	Value() interface{}
	SetValue(value interface{})
	Prev() Element //前驱
	Next() Element //后继
}
```

**复杂度：**

|操作     | 描述                            | 跳表    | avl树 | 红黑树  | b树   | b+树 |
|:-------|---------------------------------|---------|-------|--------|-------|------:|
|Insert  |添加元素，key重复添加，value为最新值 |  log(N) | log(N)| log(N)| log(N)| log(N)|
|Remove  |删除元素                          |  log(N) | log(N)| log(N)| log(N)| log(N)|
|Get     |获取元素值                        |  log(N) | log(N)| log(N)| log(N)| log(N)|
|Find     |获取元素                        |  log(N) | log(N)| log(N)| log(N)| log(N)|
|Left     |顺序遍历中，最左端的元素           |  log(1) | log(N)| log(N)| log(N)| log(N)|
|Right     |顺序遍历中，最右端的元素           |  log(1) | log(N)| log(N)| log(N)| log(N)|
|Prev     |查询key的前驱元素，即使key不存在树中  |  log(N) | log(N)| log(N)| log(N)| log(N)|
|Next     |查询key的后继元素，即使key不存在书中  |  log(N) | log(N)| log(N)| log(N)| log(N)|
|Element.Prev| 向前遍历                  |  log(N) | log(N)| log(N)| log(N)| log(1)|
|Element.Next| 向后遍历                  |  log(1) | log(N)| log(N)| log(N)| log(1)|

**使用示例：**
```golang
package main

import (
	"fmt"
	"github.com/mrtcx/plusdata/tree"
	"github.com/mrtcx/plusdata/tree/rbtree"
)

func main() {
	var tree tree.Tree = rbtree.New(tree.IntComparator) //这里使用红黑树，可以换成avl、跳表、b、b+树
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
	//查找前序和后继
	tree.Prev(2) //1
	tree.Next(2) //3
	tree.Prev(3) //2
	tree.Next(3) //nil
	tree.Next(1) //nil
	//遍历 (注意：不要在遍历的过程中，对树做添加和删除操作)
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
```

#### 跳表

跳表给每一个节点一个概率化高度，让每个节点的高度达到一种分布均衡，在节点数量够多的时候呈现:
高度大于等于1的节点数: n
高度大于等于2的节点数: n / 2
高度大于等于3的节点数: n / 2 / 2
...
在这种分布下，增、删、查操作中可以近似表现出log(N)的性能，代码量也极其简单和少量。概率近似平衡，比其他平衡树性能会差一点。

```golang
L32		0 ───────────────────────────────────────→      ....  → NUL
.		. 
.		. 
.		. 
L4      0 ───────────────────────────────────→ 53  →    ....  → NULL
        │                                      │
L3      0 ───────────────→ 28 ──────────────→  53  →    ....  →  NULL
        │                  │                   │
L2      0 ───────────────→ 28 ────→ 42 ──────→ 53  →    ....  →  NULL
        │                  │         │         │
L1      0  ─────→ 19 ───→  28  ───→ 42 ──→ 47 ─→53 →    ....  →  NULL
    	│         │        │         │     │   │   
L0      0 → 12 → 19 → 23 → 28 → 35 → 42 → 47 → 53 → 60 → .... →  NULL
		▲													 ▲
		│													 │
		front		    								     back
```
**使用示例：**
```golang
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
	//查找前序和后继
	tree.Prev(2) //1
	tree.Next(2) //3
	tree.Prev(3) //2
	tree.Next(3) //nil
	tree.Next(1) //nil
	//遍历 (注意：不要在遍历的过程中，做添加和删除操作)
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
```

**复杂度:**

|树      |Insert |  Remove | Get、Find | Left | Right | Prev | Next | Element.Next |  Element.Prev |
|:-------|-------|---------|-----------|------|-------|------|------|--------------|---------|
|跳表     | O(logN)|O(logN)|  O(logN) | O(1)| O(1)| O(logN)| O(logN)| O(log1)| O(logN)|
> 因为节点并没有采用双向列表，所有前序遍历Element.Prev为log(N)

#### avl树

plusdata中，avl为一颗严格的平衡树，每一个节点的左子树和右子树高度差不超过2， 查找性能略优于红黑树，但是因为极度的平衡，添加/删除过程中容易触发失衡导致旋转操作，比起红黑树适合读多写少的场景。
```golang
                ┌──────────────────────[13(h=3)]──────────────────────────┐
                │                              				   			  │
        ┌─────[7(h=3)]──────┐       							┌────[17(h=2)]────────┐
        │                   │       							│                     │
    ┌─[4(h=2)]─┐      	┌─[10(h=2)]─┐ 					┌─[15(h=1)]─┐   		┌─[19(h=1)]
    │          │      	│           │ 					│           │   		│            
 ┌[1(h=1)]┐ ┌[6(h=1)] [8(h=1)]┐ ┌[12(h=1)] 		 [14(h=0)]  [16(h=0)]  		[18(h=0)]   
 │        │ │                 │ │           		                               
[2]      [3][5]             [9][11]        	                              
```

**使用示例：**
```golang
package main

import (
	"fmt"
	"github.com/mrtcx/plusdata/tree"
	"github.com/mrtcx/plusdata/tree/avltree"
)

func main() {
	var tree tree.Tree = avltree.New(tree.IntComparator)
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
	//查找前序和后继
	tree.Prev(2) //1
	tree.Next(2) //3
	tree.Prev(3) //2
	tree.Next(3) //nil
	tree.Next(1) //nil
	//遍历 (注意：不要在遍历的过程中，去对树做添加和删除操作)
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
```

**复杂度:**

|树      |Insert |  Remove | Get、Find | Left | Right | Prev | Next | Element.Next |  Element.Prev |
|:-------|-------|---------|-----------|------|-------|------|------|--------------|---------|
|avl树    | O(logN)|O(logN)|  O(logN) |O(logN)| O(logN)| O(logN)| O(logN)| O(logN)| O(logN)|

#### 红黑树

红黑树根据黑色节点和红色节点规则来平衡，平衡条件比avl树相对松散，更适合频繁插入/删除。
```golang
                        ┌────────────────────100(B)────────────────────┐
                        │                                              │
                ┌──────50(R)──────┐                           ┌──────150(R)─────┐
                │                 │                           │                 │
        ┌─────30(B)────┐     ┌───70(B)───┐             ┌───130(B)───┐     ┌───170(B)───┐
        │              │     │           │             │            │     │            │
    ┌─20(R)─┐      ┌─40(R)─┐60(R)   ┌─80(R)        ┌─120(R)─┐  ┌─140(R)─┐160(R)     180(B)
    │       │      │       │        │              │        │  │        │   │     
 ┌10(B)  ┌25(B) ┌35(B)  ┌45(B)  ┌65(B)          110(B) 125(B)135(B) 145(B) 165(B) 
 │       │      │       │       │                          
5(R)   15(R)  28(R)   42(R)   55(R)
```

**使用示例：**
```golang
package main

import (
	"fmt"
	"github.com/mrtcx/plusdata/tree"
	"github.com/mrtcx/plusdata/tree/rbtree"
)

func main() {
	var tree tree.Tree = rbtree.New(tree.IntComparator)
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
	//查找前序和后继
	tree.Prev(2) //1
	tree.Next(2) //3
	tree.Prev(3) //2
	tree.Next(3) //nil
	tree.Next(1) //nil
	//遍历 (注意：不要在遍历的过程中，去对树做添加和删除操作)
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
```

**复杂度:**

|树      |Insert |  Remove | Get、Find | Left | Right | Prev | Next | Element.Next |  Element.Prev |
|:-------|-------|---------|-----------|------|-------|------|------|--------------|---------|
|红黑树树  | O(logN)|O(logN)|  O(logN) |O(logN)| O(logN)| O(logN)| O(logN)| O(logN)| O(logN)|

#### b树

在plusdata的b树中，阶数为order的b树，节点数最大为order, 节点数最小为(order + 1) / 2 - 1（root节点除外）， 子节点数比节点数多+1(叶子节点没有子节点)。
多叉平衡树，一个节点对应一块元素集，块内元素集中，局部性好，二分性能高，与avl树、红黑树相比对局部性操作有更好的性能，如遍历、区间操作。删除/插入导致节点会引起分裂和组成，导致块内元素搬迁，适合读多写少的场景。
```golang
              ┌──────────[20]───────────┐
              │                         │               
      ┌─────[10]─────┐          ┌──────[30,60]────────┐            
      │              │          │          │          │          
    [5,7]           [15]       [25]       [35]      [70,80,100]
    /  |  \         /   \      /  \      /   \     /   |   |   \
[1,2] [6] [8,9][12,13] [17][22,23][27][33,34][38] [65][75][85,90][102]
										（3阶b树）
```
**使用示例：**
```golang
package main

import (
	"fmt"
	"github.com/mrtcx/plusdata/tree"
	"github.com/mrtcx/plusdata/tree/btree"
)

func main() {
	var tree tree.Tree = btree.New(tree.IntComparator, 3) //3阶b树
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
	//查找前序和后继
	tree.Prev(2) //1
	tree.Next(2) //3
	tree.Prev(3) //2
	tree.Next(3) //nil
	tree.Next(1) //nil
	//遍历 (注意：不要在遍历的过程中，去对树做添加和删除操作)
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
```

**复杂度：**

|树      |Insert |  Remove | Get、Find | Left | Right | Prev | Next | Element.Next |  Element.Prev |
|:-------|-------|---------|-----------|------|-------|------|------|--------------|---------|
|b树      | O(logN)|O(logN)|  O(logN) |O(logN)| O(logN)| O(logN)| O(logN)| O(logN)| O(logN)|

#### b+树

在plusdata的b+树中，阶数为order的b+树，节点数最大为order - 1, 节点数最小为(order + 1) / 2 - 1（root节点除外）， 子节点数比节点数多+1(叶子节点没有子节点)。
b+树对b树的一种改良，非叶子结点不存储值，叶子结点相连，更容易遍历（最主要的是在磁盘场景中，可以将叶子结点统一放在特定的磁盘空间中，达到顺序io的目的）。 但是缺点是增加了一层索引节点，存储空间使用更大.
```golang
               ┌────────────[17]──────────────────┐ 
               │                                  │ 
        ┌─────[9]────────┐                ┌──────[30]─────┐    
        │                │                │               │    
     [2, 6]             [13]             [23]            [34]   
   /    |  \           /    \          /     \          /   \   
[1,2]←→[6]←→[8,9]←→[12,13]←→[17] ←→ [22,23]←→[30]←→[33,34]←→[38]
                           （3阶b+树）
```

**使用示例：**
```golang
package main

import (
	"fmt"
	"github.com/mrtcx/plusdata/tree"
	"github.com/mrtcx/plusdata/tree/bplustree"
)

func main() {
	var tree tree.Tree = bplustree.New(tree.IntComparator, 3) //3阶b+树
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
	//查找前序和后继
	tree.Prev(2) //1
	tree.Next(2) //3
	tree.Prev(3) //2
	tree.Next(3) //nil
	tree.Next(1) //nil
	//遍历 (注意：不要在遍历的过程中，去对树做添加和删除操作)
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
```
**复杂度：**

|树      |Insert |  Remove | Get、Find | Left | Right | Prev | Next | Element.Next |  Element.Prev |
|:-------|-------|---------|-----------|------|-------|------|------|--------------|---------|
|b+树      | O(logN)|O(logN)|  O(logN) |O(logN)| O(logN)| O(logN)| O(logN)| O(1)| O(1)|


#### 对比和选择

通常很多人认为只有在磁盘存储中会有b树、b+树的身影，因为磁盘io太过缓慢，通过读取一整块磁盘数据到内存，减少磁盘io，avl和红黑树没办法做到这一点。
而在内存中则认为没有b树和b+树的必要，这并不正确。因为cpu的L1、L2、L3-cache，通过组相连的方式，同样是按块来访问内存的。内存可以看做是磁盘的缓存，cpu的L1、L2、L3-cache可以看作是内存的缓存，在整个存储器山中，相对作用是一致的，但我们常常会忽略cpu-cache和内存的差距，b树、b+树在内存和磁盘之间发挥着作用，同样在cpu-cahce和内存之间也可以发挥着作用。（这个作用常常被我们忽略：**局部性。**） 在磁盘和内存设置b、b+树的节点大小为16kb，在cpu和内存设置节点大小为一个cpu-cache组相连的大小就能类似一样的效果。

**存储器山:**
|cpu cache | 内存 | 磁盘|
|:-------|---------|-----:|
|很快，大概是内存的100倍 |快、大概是磁盘的10万倍|很慢|

```
                            ▲
                           / \ 
                          / ▲ \ 
                         / CPU \ 
                        /___▼___\ 
                       /         \ 
                      /   L1 Cache\ 
                     / (~1ns, 32KB)\ 
                    /______▼________\ 
                   /                 \ 
                  /     L2 Cache      \ 
                 /   (~4ns, 256KB)     \ 
                /________▼______________\ 
               /                         \ 
              /        L3 Cache           \ 
             /        (~10ns, 8MB)         \ 
            /__________▼____________________\ 
           /                                 \ 
          /                                   \ 
         /             内存                     \
        /           (~100ns, 16GB)              \ 
       /______________▼__________________________\ 
      /                                           \
     /                                             \ 
    /                                               \ 
   /             	 磁盘                             \  
  / 	           (~1ms, 1TB)                        \ 
 /______________▼______________________________________\ 
```

**选择和对比:**

|树      |  性能优劣势 |
|:-------|----------:|
|跳表|实现简单，性能比其他平衡树差一点|
|b、b+树|一个块一个节点，块内元素集中在一起，二分性能强， 适合局部强的操作，比如遍历、区间读写，失衡重组和分裂时需要移动块内元素，消耗高，适合读多写少|
|avl、红黑树| 一个元素一个节点，存储上分散，局部性比b、b+树差，应对频繁插入/删除消耗少一些｜


|树      |平衡条件  | 读写 | 存储空间| 遍历能力|
|:-------|--------|------|-----|-----:|
| avl树  | 平衡条件苛刻，极易失衡|适合读多写少|一样|向前、向后便利需要log(N)查找|
| 红黑树 | 平衡条件松散|适合频繁插入/删除|一样|向前、向后便利需要log(N)查找|

|树      |平衡条件     | 读写  | 存储空间|遍历能力|
|:-------|--------|------|-----|-------:|
|b树    | 大体一致 |  大体一致   | 少了一半节点 |向前、向后便利需要log(N)查找|
|b+树 |  大体一致 |   大体一致  | 多了一层索引节点 |根据叶子结点直接向前、向后遍历|

操作符合离散、局部？读多写少、读多写多？是一件很难预判的事情，这也是为什么更多人采用红黑树，它更适合频繁的插入/删除。对各种树有一个清晰的认知，能做对更好选择。直接的方式可以拿着操作对着复杂度表看。

## 后续计划
1. 范型版本， 用interface{}来对数据类型进行一层转换，性能上有损失，会补充范型版本。

