// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package circularblocks

import (
	"github.com/mrtcx/plusdata/deque"
	"github.com/mrtcx/plusdata/internal/circularbuffer"
)

var _ deque.Deque = (*circularBlocks)(nil)

const (
	_initBlockCap = 4
	_block        = 1024 * 32 / 8 / 2 //L1-cache-32kb, interface{} - 8 * 2
)

type circularBlocks struct {
	bbs  *circularbuffer.Buffer
	size int
}

func New() *circularBlocks {
	return &circularBlocks{
		bbs: circularbuffer.New(1),
	}
}

func (d *circularBlocks) Clean() {
	d.bbs, d.size = circularbuffer.New(1), 0
}

func (d *circularBlocks) Size() int {
	return d.size
}

func (d *circularBlocks) Empty() bool {
	return d.size == 0
}

func (d *circularBlocks) Front() interface{} {
	if d.size == 0 {
		return nil
	}
	return d.bbs.Front().(*circularbuffer.Buffer).Front()
}

func (d circularBlocks) Back() interface{} {
	if d.size == 0 {
		return nil
	}
	return d.bbs.Back().(*circularbuffer.Buffer).Back()
}

func (d circularBlocks) Get(index int) interface{} {
	frontblock := d.bbs.Front().(*circularbuffer.Buffer)
	if index < frontblock.Size() {
		return frontblock.Get(index)
	}
	index -= frontblock.Size()
	return d.bbs.Get(1 + index/_block).(*circularbuffer.Buffer).Get(index % _block)
}

func (d circularBlocks) Set(index int, val interface{}) {
	frontblock := d.bbs.Front().(*circularbuffer.Buffer)
	if index < frontblock.Size() {
		frontblock.Set(index, val)
		return
	}
	index -= frontblock.Size()
	d.bbs.Get(1+index/_block).(*circularbuffer.Buffer).Set(index%_block, val)
}

func (d *circularBlocks) PushFront(val interface{}) {
	var block *circularbuffer.Buffer
	if d.size == 0 {
		d.bbs.PushFront(circularbuffer.New(_initBlockCap))
	}
	block = d.bbs.Front().(*circularbuffer.Buffer)
	if block.IsFull() {
		if block.Size() == _block {
			block = circularbuffer.New(_initBlockCap)
			block.PushFront(val)
			if d.bbs.IsFull() {
				d.expandBlock(d.bbs)
			}
			d.bbs.PushFront(block)
		} else {
			d.expandBlock(block)
			block.PushFront(val)
		}
	} else {
		block.PushFront(val)
	}
	d.size++
}

func (d *circularBlocks) PushBack(val interface{}) {
	var cirbuf *circularbuffer.Buffer
	if d.size == 0 {
		d.bbs.PushBack(circularbuffer.New(_initBlockCap))
	}
	cirbuf = d.bbs.Back().(*circularbuffer.Buffer)
	if cirbuf.IsFull() {
		if cirbuf.Size() == _block {
			cirbuf = circularbuffer.New(_initBlockCap)
			cirbuf.PushBack(val)
			if d.bbs.IsFull() {
				d.expandBlock(d.bbs)
			}
			d.bbs.PushBack(cirbuf)
		} else {
			d.expandBlock(cirbuf)
			cirbuf.PushBack(val)
		}
	} else {
		cirbuf.PushBack(val)
	}
	d.size++
}

func (d *circularBlocks) PopFront() interface{} {
	if d.size == 0 {
		return nil
	}
	frontblock := d.bbs.Front().(*circularbuffer.Buffer)
	val := frontblock.PopFront()
	if frontblock.IsEmpty() {
		d.bbs.PopFront()
		d.shrinkBlock(d.bbs, 1)
	} else {
		d.shrinkBlock(frontblock, _initBlockCap)
	}
	d.size--
	return val
}

func (d *circularBlocks) PopBack() interface{} {
	if d.size == 0 {
		return nil
	}
	backblock := d.bbs.Back().(*circularbuffer.Buffer)
	val := backblock.PopBack()
	if backblock.IsEmpty() {
		d.bbs.PopBack()
		d.shrinkBlock(d.bbs, 1)
	} else {
		d.shrinkBlock(backblock, _initBlockCap)
	}
	d.size--
	return val
}

func (d *circularBlocks) expandBlock(block *circularbuffer.Buffer) {
	addSize := block.Size()
	if block.Size() >= 1024 {
		addSize = block.Size() / 2
	}
	if addSize+block.Capacity() > _block {
		block.ResetCapacity(_block)
		return
	}
	block.ResetCapacity(block.Capacity() + addSize)
}

func (d *circularBlocks) shrinkBlock(block *circularbuffer.Buffer, minCap int) {
	if block.Capacity() > minCap && block.Size() < block.Capacity()/4 {
		block.ResetCapacity(block.Capacity() / 4)
	}
}
