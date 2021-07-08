#include <vulkan/vulkan.h>

#include <stdio.h>

void check_vk_extensions() {
	// Check Vulkan extensions.
	uint32_t nextensions = 0;
	vkEnumerateInstanceExtensionProperties(NULL, &nextensions, NULL);
	printf("nextensions: %d\n", nextensions);
}
