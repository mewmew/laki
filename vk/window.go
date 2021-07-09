package vk

// #define GLFW_INCLUDE_VULKAN
// #include <GLFW/glfw3.h>
import "C"

func InitWindow() *C.GLFWwindow {
	dbg.Println("vk.InitWindow")
	// Initialize GLFW.
	C.glfwInit()
	C.glfwWindowHint(C.GLFW_CLIENT_API, C.GLFW_NO_API) // skip OpenGL context.
	C.glfwWindowHint(C.GLFW_RESIZABLE, C.GLFW_FALSE)
	// Create window.
	win := C.glfwCreateWindow(WindowWidth, WindowHeight, C.CString(AppTitle), nil, nil)
	return win
}

func CleanupWindow(win *C.GLFWwindow) {
	dbg.Println("vk.CleanupWindow")
	// Terminate window.
	C.glfwDestroyWindow(win)
	// Terminate GLFW.
	C.glfwTerminate()
}
