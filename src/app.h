#ifndef __APP_H__
#define __APP_H__

#define GLFW_INCLUDE_VULKAN // includes vulkan/vulkan.h
#include <GLFW/glfw3.h>

typedef struct {
	// GLFW.
	GLFWwindow *win;
	// Vulkan.
	VkInstance *instance;
	VkDebugUtilsMessengerEXT *debug_messanger;
} App;

#endif // #ifndef __APP_H__
