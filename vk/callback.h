#ifndef __CALLBACK_H__
#define __CALLBACK_H__

#include <vulkan/vulkan.h>

extern VkBool32 debugCallback(
	VkDebugUtilsMessageSeverityFlagBitsEXT messageSeverity,
	VkDebugUtilsMessageTypeFlagsEXT messageTypes,
	const VkDebugUtilsMessengerCallbackDataEXT *pCallbackData,
	void *pUserData);

#endif // #ifndef __CALLBACK_H__
