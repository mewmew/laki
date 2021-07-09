#include "defs.h"
#include "vk.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void populate_debug_messanger_create_info(VkDebugUtilsMessengerCreateInfoEXT *create_info) {
	create_info->sType = VK_STRUCTURE_TYPE_DEBUG_UTILS_MESSENGER_CREATE_INFO_EXT;
	create_info->messageSeverity = VK_DEBUG_UTILS_MESSAGE_SEVERITY_VERBOSE_BIT_EXT | VK_DEBUG_UTILS_MESSAGE_SEVERITY_INFO_BIT_EXT | VK_DEBUG_UTILS_MESSAGE_SEVERITY_WARNING_BIT_EXT | VK_DEBUG_UTILS_MESSAGE_SEVERITY_ERROR_BIT_EXT;
	create_info->messageType = VK_DEBUG_UTILS_MESSAGE_TYPE_GENERAL_BIT_EXT | VK_DEBUG_UTILS_MESSAGE_TYPE_VALIDATION_BIT_EXT | VK_DEBUG_UTILS_MESSAGE_TYPE_PERFORMANCE_BIT_EXT;
	create_info->pfnUserCallback = debug_callback;
	create_info->pUserData = NULL; // optional.
}

VkDebugUtilsMessengerEXT * init_debug_messanger(VkInstance *instance) {
	VkDebugUtilsMessengerCreateInfoEXT *debug_messanger_create_info = calloc(1, sizeof(VkDebugUtilsMessengerCreateInfoEXT));
	populate_debug_messanger_create_info(debug_messanger_create_info);
	VkDebugUtilsMessengerEXT *debug_messenger = calloc(1, sizeof(VkDebugUtilsMessengerEXT));
	VkResult result = CreateDebugUtilsMessengerEXT(*instance, debug_messanger_create_info, NULL, debug_messenger);
	if (result!= VK_SUCCESS) {
		fprintf(stderr, "unable to register debug messanger (result=%d).\n", result);
	}
	return debug_messenger;
}

VkBool32 debug_callback(
	VkDebugUtilsMessageSeverityFlagBitsEXT messageSeverity,
	VkDebugUtilsMessageTypeFlagsEXT messageTypes,
	const VkDebugUtilsMessengerCallbackDataEXT *pCallbackData,
	void *pUserData) {
	if ((messageSeverity&VK_DEBUG_UTILS_MESSAGE_SEVERITY_ERROR_BIT_EXT) != 0) {
		fprintf(stderr, "validation layer (error): %s\n", pCallbackData->pMessage);
	}
	if ((messageSeverity&VK_DEBUG_UTILS_MESSAGE_SEVERITY_WARNING_BIT_EXT) != 0) {
		fprintf(stderr, "validation layer (warning): %s\n", pCallbackData->pMessage);
	}
	if ((messageSeverity&VK_DEBUG_UTILS_MESSAGE_SEVERITY_INFO_BIT_EXT) != 0) {
		fprintf(stderr, "validation layer (info): %s\n", pCallbackData->pMessage);
	}
	if ((messageSeverity&VK_DEBUG_UTILS_MESSAGE_SEVERITY_VERBOSE_BIT_EXT) != 0) {
		fprintf(stderr, "validation layer (verbose): %s\n", pCallbackData->pMessage);
	}
	return VK_FALSE;
}

VkResult CreateDebugUtilsMessengerEXT(
	VkInstance instance,
	const VkDebugUtilsMessengerCreateInfoEXT *pCreateInfo,
	const VkAllocationCallbacks *pAllocator,
	VkDebugUtilsMessengerEXT *pMessenger) {
	PFN_vkCreateDebugUtilsMessengerEXT func = (PFN_vkCreateDebugUtilsMessengerEXT)vkGetInstanceProcAddr(instance, "vkCreateDebugUtilsMessengerEXT");
	if (func != NULL) {
		return func(instance, pCreateInfo, pAllocator, pMessenger);
	}
	return VK_ERROR_EXTENSION_NOT_PRESENT;
}

void DestroyDebugUtilsMessengerEXT(
	VkInstance instance,
	VkDebugUtilsMessengerEXT messenger,
	const VkAllocationCallbacks *pAllocator) {
	PFN_vkDestroyDebugUtilsMessengerEXT func = (PFN_vkDestroyDebugUtilsMessengerEXT)vkGetInstanceProcAddr(instance, "vkDestroyDebugUtilsMessengerEXT");
	if (func != NULL) {
		func(instance, messenger, pAllocator);
	}
}
