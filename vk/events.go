package vk

// #define GLFW_INCLUDE_VULKAN
// #include <GLFW/glfw3.h>
import "C"

import (
	"time"
)

func EventLoop(app *App) {
	// Poll events.
	dbg.Println("vk.EventLoop")
	currentFrame := 0
	now := time.Now()
	for C.glfwWindowShouldClose(app.win) == 0 {
		currentFrame++
		C.glfwPollEvents()

		// Render frame.
		if err := drawFrame(app); err != nil {
			warn.Println("%+v", err) // print warning and continue
		}
		if time.Since(now) >= time.Second {
			now = time.Now()
			dbg.Println("fps:", currentFrame)
			currentFrame = 0
		}
	}
	dbg.Println("waiting for device to become idle")
	if result := C.vkDeviceWaitIdle(*app.device); result != C.VK_SUCCESS {
		warn.Printf("unable to wait for device to become idle (result=%d)", result)
	}
}
