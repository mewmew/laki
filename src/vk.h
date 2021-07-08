#ifndef __VK_H__
#define __VK_H__

#define GLFW_INCLUDE_VULKAN // includes vulkan/vulkan.h
#include <GLFW/glfw3.h>

#include <stdbool.h>

// external

extern VkInstance * init_vulkan();
extern void cleanup_vulkan(VkInstance *instance);

// internal

extern VkInstance * create_instance();
extern const char ** get_extensions(uint32_t *pnenabled_extensions);
extern const char ** get_layers(uint32_t *pnenabled_layers);
extern bool has_extension(VkExtensionProperties *extensions, int nextensions, const char *extension_name);
extern bool has_layer(VkLayerProperties *layers, int nlayers, const char *layer_name);

#endif // #ifndef __VK_H__
