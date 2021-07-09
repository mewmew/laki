#ifndef __MALLOC_H__
#define __MALLOC_H__

#include <vulkan/vulkan.h>

extern VkInstance * new_VkInstance();
extern VkInstanceCreateInfo * new_VkInstanceCreateInfo();
extern VkPhysicalDevice * new_VkPhysicalDevice();
extern VkDebugUtilsMessengerCreateInfoEXT * new_VkDebugUtilsMessengerCreateInfoEXT();
extern VkDebugUtilsMessengerEXT * new_VkDebugUtilsMessengerEXT();

#endif // #ifndef __MALLOC_H__
