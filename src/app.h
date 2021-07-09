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
	VkPhysicalDevice *device;
} App;

App * new_app();

#endif // #ifndef __APP_H__
