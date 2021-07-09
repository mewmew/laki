package vk

// #include "events.h"
//
//#cgo CFLAGS: -I../src
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
