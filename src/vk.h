#ifndef __VK_H__
#define __VK_H__

#define GLFW_INCLUDE_VULKAN // includes vulkan/vulkan.h
#include <GLFW/glfw3.h>

extern void init_vulkan();

extern VkInstance create_instance();

extern void check_vk_extensions();

#endif // #ifndef __VK_H__
