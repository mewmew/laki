package vk

// #include "app.h"
//
//#cgo CFLAGS: -I../src
import "C"

type App struct {
	// GLFW.
	win *C.GLFWwindow
	// Vulkan.
	instance        *C.VkInstance
	debug_messanger *C.VkDebugUtilsMessengerEXT
	device          *C.VkPhysicalDevice
}
