#ifndef __INVOKE_H__
#define __INVOKE_H__

#include <vulkan/vulkan.h>

extern VkResult invoke_CreateDebugUtilsMessengerEXT(
	PFN_vkCreateDebugUtilsMessengerEXT fn,
	VkInstance instance,
	const VkDebugUtilsMessengerCreateInfoEXT *pCreateInfo,
	const VkAllocationCallbacks *pAllocator,
	VkDebugUtilsMessengerEXT *pMessenger);

extern void invoke_DestroyDebugUtilsMessengerEXT(
	PFN_vkDestroyDebugUtilsMessengerEXT fn,
	VkInstance instance,
	VkDebugUtilsMessengerEXT messenger,
	const VkAllocationCallbacks *pAllocator);

#endif // #ifndef __INVOKE_H__
