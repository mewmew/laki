package vk

// #define GLFW_INCLUDE_VULKAN
// #include <GLFW/glfw3.h>
import "C"

type App struct {
	// GLFW.
	win *C.GLFWwindow
	// Vulkan.
	instance        *C.VkInstance
	debug_messanger *C.VkDebugUtilsMessengerEXT
	device          *C.VkPhysicalDevice
}
