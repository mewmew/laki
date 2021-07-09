#ifndef __VK_H__
#define __VK_H__

#define GLFW_INCLUDE_VULKAN // includes vulkan/vulkan.h
#include <GLFW/glfw3.h>

#include <stdbool.h>

#include "app.h"

// external

extern void init_vulkan(App *app);
extern void cleanup_vulkan(App *app);

// internal

extern VkInstance * create_instance();
extern const char ** get_extensions(uint32_t *pnenabled_extensions);
extern const char ** get_layers(uint32_t *pnenabled_layers);
extern bool has_extension(VkExtensionProperties *extensions, int nextensions, const char *extension_name);
extern bool has_layer(VkLayerProperties *layers, int nlayers, const char *layer_name);

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

#endif // #ifndef __VK_H__
