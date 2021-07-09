package vk

// #include "window.h"
//
//#cgo CFLAGS: -I../src
import "C"

import (
	"fmt"
)

func InitWindow() *C.GLFWwindow {
	fmt.Println("vk.InitWindow")
	// Initialize GLFW.
	C.glfwInit()
	C.glfwWindowHint(C.GLFW_CLIENT_API, C.GLFW_NO_API) // skip OpenGL context.
	C.glfwWindowHint(C.GLFW_RESIZABLE, C.GLFW_FALSE)
	// Create window.
	win := C.glfwCreateWindow(WindowWidth, WindowHeight, C.CString(AppTitle), nil, nil)
	return win
}

func CleanupWindow(win *C.GLFWwindow) {
	fmt.Println("vk.CleanupWindow")
	// Terminate window.
	C.glfwDestroyWindow(win)
	// Terminate GLFW.
	C.glfwTerminate()
}
