package vk

// #include "invoke.h"
//
// VkResult invoke_CreateDebugUtilsMessengerEXT(
// 	PFN_vkCreateDebugUtilsMessengerEXT fn,
// 	VkInstance instance,
// 	const VkDebugUtilsMessengerCreateInfoEXT *pCreateInfo,
// 	const VkAllocationCallbacks *pAllocator,
// 	VkDebugUtilsMessengerEXT *pMessenger) {
// 	return fn(instance, pCreateInfo, pAllocator, pMessenger);
// }
//
// void invoke_DestroyDebugUtilsMessengerEXT(
// 	PFN_vkDestroyDebugUtilsMessengerEXT fn,
// 	VkInstance instance,
// 	VkDebugUtilsMessengerEXT messenger,
// 	const VkAllocationCallbacks *pAllocator) {
// 	fn(instance, messenger, pAllocator);
// }
import "C"
