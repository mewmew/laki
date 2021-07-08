#include "window.h"

const int WINDOW_WIDTH = 1024;
const int WINDOW_HEIGHT = 768;
const char *WINDOW_TITLE = "laki";

GLFWwindow * init_window() {
	// Initialize GLFW.
	glfwInit();
	glfwWindowHint(GLFW_CLIENT_API, GLFW_NO_API); // skip OpenGL context.
	glfwWindowHint(GLFW_RESIZABLE, GLFW_FALSE);
	// Create window.
	GLFWwindow *win = glfwCreateWindow(WINDOW_WIDTH, WINDOW_HEIGHT, WINDOW_TITLE, NULL, NULL);
	return win;
}

void cleanup_window(GLFWwindow *win) {
	// Terminate window.
	glfwDestroyWindow(win);
	// Terminate GLFW.
	glfwTerminate();
}
