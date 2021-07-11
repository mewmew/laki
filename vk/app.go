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
	graphicsQueue  *C.VkQueue
	presentQueue   *C.VkQueue
	surface        *C.VkSurfaceKHR
	*QueueFamilyIndices
	swapchainSupportInfo *SwapchainSupportInfo
	swapchain            *C.VkSwapchainKHR
	swapchainImageFormat C.VkFormat
	swapchainExtent      C.VkExtent2D
	swapchainImgs        []C.VkImage
	swapchainImgViews    []C.VkImageView
}

func newApp() *App {
	return &App{
		QueueFamilyIndices: newQueueFamilyIndices(),
	}
}

type QueueFamilyIndices struct {
	graphicsQueueFamilyIndex int
	presentQueueFamilyIndex  int
}

func newQueueFamilyIndices() *QueueFamilyIndices {
	return &QueueFamilyIndices{
		graphicsQueueFamilyIndex: -1,
		presentQueueFamilyIndex:  -1,
	}
}

func (queueFamilyIndices *QueueFamilyIndices) Indices() []int {
	return []int{
		queueFamilyIndices.graphicsQueueFamilyIndex,
		queueFamilyIndices.presentQueueFamilyIndex,
	}
}

type SwapchainSupportInfo struct {
	surfaceCapabilities *C.VkSurfaceCapabilitiesKHR
	surfaceFormats      []C.VkSurfaceFormatKHR
	presentModes        []C.VkPresentModeKHR
}
