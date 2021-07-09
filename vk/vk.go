package vk

// void c_main();
//
//#cgo LDFLAGS: -llaki -L../
//#cgo pkg-config: glfw3
//#cgo pkg-config: vulkan
import "C"

func Init() {
	C.c_main()
}
