#ifndef __MALLOC_H__
#define __MALLOC_H__

#include <vulkan/vulkan.h>

extern VkInstance * new_VkInstance();
extern VkInstanceCreateInfo * new_VkInstanceCreateInfo();
extern VkPhysicalDevice * new_VkPhysicalDevice();
extern VkDebugUtilsMessengerCreateInfoEXT * new_VkDebugUtilsMessengerCreateInfoEXT();
extern VkDebugUtilsMessengerEXT * new_VkDebugUtilsMessengerEXT();
extern VkDevice * new_VkDevice();
extern VkPhysicalDeviceFeatures * new_VkPhysicalDeviceFeatures();
extern VkDeviceCreateInfo * new_VkDeviceCreateInfo();
extern VkDeviceQueueCreateInfo * new_VkDeviceQueueCreateInfo();
extern VkQueue * new_VkQueue();
extern VkSurfaceKHR * new_VkSurfaceKHR();
extern VkSwapchainKHR * new_VkSwapchainKHR();
extern VkShaderModule * new_VkShaderModule();
extern VkPipelineLayout * new_VkPipelineLayout();
extern VkRenderPass * new_VkRenderPass();

#endif // #ifndef __MALLOC_H__
