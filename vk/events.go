package vk

// #define GLFW_INCLUDE_VULKAN
// #include <GLFW/glfw3.h>
import "C"

import (
	"fmt"
)

func EventLoop(win *C.GLFWwindow) {
	// Poll events.
	fmt.Println("vk.EventLoop")
	for C.glfwWindowShouldClose(win) == 0 {
		C.glfwPollEvents()
	}
}
