package vk

// #define GLFW_INCLUDE_VULKAN
// #include <GLFW/glfw3.h>
//
// #include "callback.h"
import "C"

func InitWindow(app *App) *C.GLFWwindow {
	dbg.Println("vk.InitWindow")
	// Initialize GLFW.
	C.glfwInit()
	C.glfwWindowHint(C.GLFW_CLIENT_API, C.GLFW_NO_API) // skip OpenGL context.
	//C.glfwWindowHint(C.GLFW_RESIZABLE, C.GLFW_FALSE)
	// Create window.
	win := C.glfwCreateWindow(WindowWidth, WindowHeight, C.CString(AppTitle), nil, nil)
	_framebufferResizeCallback = func(win *C.GLFWwindow, width, height int) {
		dbg.Println("framebufferResizeCallback")
		dbg.Println("   width:", width)
		dbg.Println("   height:", height)
		app.framebufferResized = true
	}
	C.glfwSetFramebufferSizeCallback(win, (*[0]byte)(C.framebufferResizeCallback))
	return win
}

func CleanupWindow(win *C.GLFWwindow) {
	dbg.Println("vk.CleanupWindow")
	// Terminate window.
	C.glfwDestroyWindow(win)
	// Terminate GLFW.
	C.glfwTerminate()
}
