// Copyright 2009 The Go Authors. All rights reserved.

// Package pooll provides basic pool implementation such as Memory
package pool

// //#include <sys/syscall.h>
// //#define gettid() syscall(__NR_gettid)
// import "C"

// func Gettid() int {
// 	cc := C.gettid(1)
// }
var BP *BytePool = nil

func init() {
	BP = NewBytePool(8, 102400)
}

func SetBytePoolMax(max int) {
	if BP.End < max {
		BP = NewBytePool(8, max)
	}
}
