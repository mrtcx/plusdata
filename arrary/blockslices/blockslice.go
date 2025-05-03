// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package blockslices

import "github.com/mrtcx/plusdata/arrary"

var _ arrary.Arrary = (*blockSlice)(nil)

const (
	_initBlockCap = 4
	_block        = 1024 * 32 / 8 / 2 //L1-cache-32kb, interface{} - 8 * 2
)

type blockSlice struct {
	size   int
	blocks [][]interface{}
}

func New() *blockSlice {
	return &blockSlice{}
}

func (b *blockSlice) Clean() {
	b.size, b.blocks = 0, nil
}

func (b *blockSlice) Size() int {
	return b.size
}

func (b *blockSlice) Empty() bool {
	return b.size == 0
}

func (b *blockSlice) Front() interface{} {
	if b.size == 0 {
		return nil
	}
	return b.blocks[0][0]
}

func (b *blockSlice) Back() interface{} {
	if b.size == 0 {
		return nil
	}
	idx := b.size - 1
	return b.blocks[idx/_block][idx%_block]
}

func (b *blockSlice) Get(index int) interface{} {
	return b.blocks[index/_block][index%_block]
}

func (b *blockSlice) Set(index int, val interface{}) {
	b.blocks[index/_block][index%_block] = val
}

func (b *blockSlice) PushBack(value interface{}) {
	blockIdx := b.size / _block
	if blockIdx == len(b.blocks) {
		b.blocks = append(b.blocks, make([]interface{}, 0, _initBlockCap))
	}
	b.blocks[blockIdx] = append(b.blocks[blockIdx], value)
	b.size++
}

func (b *blockSlice) PopBack() interface{} {
	if b.size == 0 {
		return nil
	}
	index := b.size - 1
	blockIdx := index / _block
	posIdx := index % _block
	pv := b.blocks[blockIdx][posIdx]
	b.blocks[blockIdx] = b.blocks[blockIdx][:posIdx]
	if posIdx == 0 {
		b.blocks = b.blocks[:blockIdx]
		b.shrink1()
	} else {
		b.shrink2()
	}
	b.size--
	return pv
}

func (b *blockSlice) shrink1() {
	if cap(b.blocks)/4 > len(b.blocks) {
		newblock := make([][]interface{}, len(b.blocks))
		copy(newblock, b.blocks)
		b.blocks = newblock
	}
}

func (b *blockSlice) shrink2() {
	block := b.blocks[len(b.blocks)-1]
	if cap(block)/4 > len(block) {
		newblock := make([]interface{}, len(block))
		copy(newblock, block)
		b.blocks[len(b.blocks)-1] = newblock
	}
}
