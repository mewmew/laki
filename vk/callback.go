package vk

// #define GLFW_INCLUDE_VULKAN
// #include <GLFW/glfw3.h>
import "C"

import (
	"log"
	"os"
	"unsafe"

	"github.com/mewkiz/pkg/term"
)

var (
	// dbg is a logger with the "vk:" prefix which logs debug messages to
	// standard error.
	dbgValidationLayer = log.New(os.Stderr, term.CyanBold("vk (validation layer):")+" ", 0)
	// warn is a logger with the "vk:" prefix which logs warning messages to
	// standard error.
	warnValidationLayer = log.New(os.Stderr, term.RedBold("vk (validation layer):")+" ", log.Lshortfile)
)

//export debugCallback
func debugCallback(messageSeverity C.VkDebugUtilsMessageSeverityFlagBitsEXT, messageTypes C.VkDebugUtilsMessageTypeFlagsEXT, pCallbackData *C.VkDebugUtilsMessengerCallbackDataEXT, pUserData unsafe.Pointer) C.VkBool32 {
	message := C.GoString(pCallbackData.pMessage)
	if (messageSeverity & C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_ERROR_BIT_EXT) != 0 {
		warnValidationLayer.Println("(error):", message)
	}
	if (messageSeverity & C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_WARNING_BIT_EXT) != 0 {
		warnValidationLayer.Println("(warning):", message)
	}
	if (messageSeverity & C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_INFO_BIT_EXT) != 0 {
		dbgValidationLayer.Println("(info):", message)
	}
	if (messageSeverity & C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_VERBOSE_BIT_EXT) != 0 {
		dbgValidationLayer.Println("(verbose):", message)
	}
	return C.VK_FALSE
}

//export framebufferResizeCallback
func framebufferResizeCallback(win *C.GLFWwindow, width, height C.int) {
	if _framebufferResizeCallback != nil {
		_framebufferResizeCallback(win, int(width), int(height))
	}
}

var _framebufferResizeCallback func(win *C.GLFWwindow, width, height int)
