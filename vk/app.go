package vk

// #define GLFW_INCLUDE_VULKAN
// #include <GLFW/glfw3.h>
import "C"

type App struct {
	// GLFW.
	win *C.GLFWwindow
	// Vulkan.
	instance       *C.VkInstance
	debugMessanger *C.VkDebugUtilsMessengerEXT
	physicalDevice *C.VkPhysicalDevice
	device         *C.VkDevice
}
