#include "events.h"

void event_loop(GLFWwindow *win) {
	// Poll events.
	while (!glfwWindowShouldClose(win)) {
		glfwPollEvents();
	}
}
