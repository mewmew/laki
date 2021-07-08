#ifndef __VK_H__
#define __VK_H__

#define GLFW_INCLUDE_VULKAN // includes vulkan/vulkan.h
#include <GLFW/glfw3.h>

#include <stdbool.h>

extern const bool enable_validation_layers;
extern const char *enabled_validation_layers[];
extern const int nenabled_validation_layers;

extern VkInstance * init_vulkan();
extern VkInstance * create_instance();
void cleanup_vulkan(VkInstance *instance);

extern void check_extensions();
extern bool check_validation_layers();

extern bool has_layer(VkLayerProperties *layers, int nlayers, const char *layer_name);

#endif // #ifndef __VK_H__
