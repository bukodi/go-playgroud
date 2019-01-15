package main

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"unsafe"
)

func TestFileIO(t *testing.T) {
	var buff1 = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8}
	fmt.Printf("%#v\n", buff1)
	fmt.Printf("%#v\n", &buff1)
	fmt.Printf("%#v\n", buff1[3])
	fmt.Printf("%#v\n", &buff1[0])

	fmt.Printf("%#v\n", unsafe.Pointer(&buff1))
	fmt.Printf("%#v\n", unsafe.Pointer(&buff1[0]))
	fmt.Printf("%#v\n", unsafe.Pointer(&buff1[3]))

	var buffView []byte
	{
		sh := (*reflect.SliceHeader)(unsafe.Pointer(&buffView))
		sh.Data = uintptr(unsafe.Pointer(&buff1[0])) + 2
		sh.Len = 4
		sh.Cap = 4
	}
	fmt.Printf("%#v\n", buffView)
	runtime.KeepAlive(buff1)

}
