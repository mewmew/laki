package vk

// #include "app.h"
// #include "events.h"
// #include "vk.h"
// #include "window.h"
//
//#cgo CFLAGS: -I../src
//#cgo LDFLAGS: -llaki -L../
//#cgo pkg-config: glfw3
//#cgo pkg-config: vulkan
import "C"

func Init() {
	app := C.new_app()
	app.win = InitWindow()
	defer CleanupWindow(app.win)
	C.init_vulkan(app)
	defer C.cleanup_vulkan(app)

	EventLoop(app.win)
}
