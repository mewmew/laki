package vk

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
