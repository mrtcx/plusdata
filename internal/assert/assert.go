// Copyright (c) 2025 Tian ChunXing. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package assert

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"testing"
)

func Equal(t *testing.T, got interface{}, wanted interface{}) {
	if wanted == nil && got != nil && reflect.ValueOf(got).IsNil() {
		got = nil
	}
	if got != wanted {
		logf("wanted %v got %+v", wanted, got)
		t.Fail()
	}
}

func NotEqual(t *testing.T, got interface{}, unwanted interface{}) {
	if unwanted == nil && got != nil && reflect.ValueOf(got).IsNil() {
		got = nil
	}
	if got == unwanted {
		logf("unwanted %+v", unwanted)
		t.Fail()
	}
}

func Greater(t *testing.T, maxer, miner int) {
	if maxer <= miner {
		logf("wanted greater %d, got %d", miner, maxer)
		t.Fail()
	}
}

func logf(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	file = path.Base(file)
	fmt.Printf("    "+fmt.Sprintf("%s:%d: ", file, line)+format+"\n", args...)
}
