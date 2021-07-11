package vk

// #include <stdint.h>
// #include <stdlib.h>
import "C"

import (
	"reflect"
	"unsafe"
)

// data is of C type `char **`.
func getStringSlice(data unsafe.Pointer, n int) []string {
	sh := reflect.SliceHeader{
		Data: uintptr(data),
		Len:  n,
		Cap:  n,
	}
	slice := *(*[]*C.char)(unsafe.Pointer(&sh))
	var ss []string
	for i := range slice {
		s := C.GoString(slice[i])
		ss = append(ss, s)
	}
	return ss
}

func getCStringSlice(ss []string) **C.char {
	dst := make([]*C.char, len(ss))
	for i := range dst {
		dst[i] = C.CString(ss[i])
	}
	return &dst[0]
}

func getCUintSlice(xs []int) *C.uint {
	dst := make([]C.uint, len(xs))
	for i := range dst {
		dst[i] = C.uint(xs[i])
	}
	return &dst[0]
}

func getCUintSliceFromBytes(buf []byte) *C.uint32_t {
	// Allocate 4-byte aligned buffer in C.
	n := len(buf)
	if n%4 != 0 {
		n += 4 - n%4
	}
	_buf := C.calloc(C.size_t(n), 1)
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(_buf)),
		Len:  len(buf),
		Cap:  len(buf),
	}
	slice := *(*[]byte)(unsafe.Pointer(&sh))
	copy(slice, buf)
	return (*C.uint32_t)(unsafe.Pointer(&slice[0]))
}
