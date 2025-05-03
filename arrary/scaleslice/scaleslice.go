// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package scaleslice

import "github.com/mrtcx/plusdata/arrary"

var _ arrary.Arrary = (*scaleslice)(nil)

type scaleslice struct {
	slice []interface{}
}

func New() *scaleslice {
	return &scaleslice{}
}

func (s *scaleslice) Clean() {
	s.slice = nil
}

func (s *scaleslice) Size() int {
	return len(s.slice)
}

func (s *scaleslice) Empty() bool {
	return len(s.slice) == 0
}

func (s *scaleslice) Front() interface{} {
	if len(s.slice) == 0 {
		return nil
	}
	return s.slice[0]
}

func (s *scaleslice) Back() interface{} {
	if len(s.slice) == 0 {
		return nil
	}
	return s.slice[len(s.slice)-1]
}

func (s *scaleslice) Get(index int) interface{} {
	return s.slice[index]
}

func (s *scaleslice) Set(index int, val interface{}) {
	s.slice[index] = val
}

func (s *scaleslice) PushBack(value interface{}) {
	s.slice = append(s.slice, value)
}

func (s *scaleslice) PopBack() interface{} {
	if len(s.slice) == 0 {
		return nil
	}
	pv := s.slice[len(s.slice)-1]
	s.slice = s.slice[0 : len(s.slice)-1]
	if cap(s.slice)/8 > len(s.slice) {
		newblock := make([]interface{}, len(s.slice))
		copy(newblock, s.slice)
		s.slice = newblock
	}
	return pv
}
