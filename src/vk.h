#ifndef __VK_H__
#define __VK_H__

#define GLFW_INCLUDE_VULKAN // includes vulkan/vulkan.h
#include <GLFW/glfw3.h>

#include <stdbool.h>

#include "app.h"

// internal

void populate_debug_messanger_create_info(VkDebugUtilsMessengerCreateInfoEXT *create_info);
extern VkDebugUtilsMessengerEXT * init_debug_messanger(VkInstance *instance);
extern VkBool32 debug_callback(
	VkDebugUtilsMessageSeverityFlagBitsEXT messageSeverity,
	VkDebugUtilsMessageTypeFlagsEXT messageTypes,
	const VkDebugUtilsMessengerCallbackDataEXT *pCallbackData,
	void *pUserData);

extern VkResult CreateDebugUtilsMessengerEXT(
	VkInstance instance,
	const VkDebugUtilsMessengerCreateInfoEXT *pCreateInfo,
	const VkAllocationCallbacks *pAllocator,
	VkDebugUtilsMessengerEXT *pMessenger);

extern void DestroyDebugUtilsMessengerEXT(
	VkInstance instance,
	VkDebugUtilsMessengerEXT messenger,
	const VkAllocationCallbacks *pAllocator);

extern VkPhysicalDevice * init_device(VkInstance *instance);
extern bool is_suitable_defice(VkPhysicalDevice *device);

#endif // #ifndef __VK_H__
