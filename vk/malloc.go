package vk

// #include "malloc.h"
//
// #include <stdlib.h>
//
// VkInstance * new_VkInstance() {
//    return calloc(1, sizeof(VkInstance));
// }
//
// VkInstanceCreateInfo * new_VkInstanceCreateInfo() {
//    return calloc(1, sizeof(VkInstanceCreateInfo));
// }
//
// VkPhysicalDevice * new_VkPhysicalDevice() {
//    return calloc(1, sizeof(VkPhysicalDevice));
// }
//
// VkDebugUtilsMessengerCreateInfoEXT * new_VkDebugUtilsMessengerCreateInfoEXT() {
//    return calloc(1, sizeof(VkDebugUtilsMessengerCreateInfoEXT));
// }
import "C"
