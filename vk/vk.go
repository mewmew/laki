// TODO: continue at https://vulkan-tutorial.com/en/Drawing_a_triangle/Setup/Physical_devices_and_queue_families

package vk

// #include "app.h"
// #include "vk.h"
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
	InitVulkan(app)
	defer CleanupVulkan(app)

	EventLoop(app.win)
}

func InitVulkan(app *C.App) {
	app.instance = C.create_instance()
	app.debug_messanger = C.init_debug_messanger(app.instance)
	app.device = C.init_device(app.instance)
}

func CleanupVulkan(app *C.App) {
	app.device = nil
	C.DestroyDebugUtilsMessengerEXT(*(app.instance), *(app.debug_messanger), nil)
	C.vkDestroyInstance(*(app.instance), nil)
}
