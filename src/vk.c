#include "defs.h"
#include "vk.h"

#include <stdio.h>
#include <stdlib.h>

VkInstance * init_vulkan() {
	VkInstance *instance = create_instance();
	return instance;
}

VkInstance * create_instance() {
	VkApplicationInfo app_info = {};
	app_info.sType = VK_STRUCTURE_TYPE_APPLICATION_INFO;
	app_info.pApplicationName = APP_TITLE;
	app_info.applicationVersion = VK_MAKE_VERSION(1, 0, 0);
	app_info.pEngineName = "No Engine";
	app_info.engineVersion = VK_MAKE_VERSION(1, 0, 0);
	app_info.apiVersion = VK_API_VERSION_1_0;

	uint32_t nglfw_extensions = 0;
	const char **glfw_extensions = glfwGetRequiredInstanceExtensions(&nglfw_extensions);
	printf("nglfw_extensions: %d\n", nglfw_extensions);
	for (int i = 0; i < nglfw_extensions; i++) {
		printf("   glfw_extension: %s\n", glfw_extensions[i]);
	}

	VkInstanceCreateInfo create_info = {};
	create_info.sType = VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO;
	create_info.flags = 0; // reserved.
	create_info.pApplicationInfo = &app_info;
	create_info.enabledExtensionCount = nglfw_extensions;
	create_info.ppEnabledExtensionNames = glfw_extensions;
	create_info.enabledLayerCount = 0;

	VkInstance *instance = calloc(1, sizeof(VkInstance));
	VkResult result = vkCreateInstance(&create_info, NULL, instance);
	if (result != VK_SUCCESS) {
		fprintf(stderr, "unable to create instance.\n");
		exit(EXIT_FAILURE);
	}
	return instance;
}

void check_vk_extensions() {
	// Check Vulkan extensions.
	uint32_t nextensions = 0;
	vkEnumerateInstanceExtensionProperties(NULL, &nextensions, NULL);
	VkExtensionProperties extensions[nextensions];
	//extensions = calloc(nextensions, sizeof(VkExtensionProperties));
	vkEnumerateInstanceExtensionProperties(NULL, &nextensions, extensions);
	printf("nextensions: %d\n", nextensions);
	for (int i = 0; i < nextensions; i++) {
		printf("   extension: %-40s (0x%08X)\n", extensions[i].extensionName, extensions[i].specVersion);
	}
}

void cleanup_vulkan(VkInstance *instance) {
	vkDestroyInstance(*instance, NULL);
}
