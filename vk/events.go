package vk

// #define GLFW_INCLUDE_VULKAN
// #include <GLFW/glfw3.h>
import "C"

func EventLoop(app *App) {
	// Poll events.
	dbg.Println("vk.EventLoop")
	for C.glfwWindowShouldClose(app.win) == 0 {
		C.glfwPollEvents()

		// Render frame.
		if err := drawFrame(app); err != nil {
			warn.Println("%+v", err) // print warning and continue
		}
	}
}
