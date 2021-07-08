#include <stdio.h>

#define GLFW_INCLUDE_VULKAN
#include <GLFW/glfw3.h>

int main(int argc, char **argv) {
	// Initialize GLFW window.
	glfwInit();
	glfwWindowHint(GLFW_CLIENT_API, GLFW_NO_API);
	GLFWwindow *win = glfwCreateWindow(1024, 768, "laki", NULL, NULL);

	// Check Vulkan extensions.
	uint32_t nextensions = 0;
	vkEnumerateInstanceExtensionProperties(NULL, &nextensions, NULL);
	printf("nextensions: %d\n", nextensions);

	// Poll events.
	while (!glfwWindowShouldClose(win)) {
		glfwPollEvents();
	}

	// Terminate GLFW window.
	glfwDestroyWindow(win);
	glfwTerminate();

	return 42;
}
