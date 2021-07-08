#include "defs.h"
#include "vk.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

const bool enable_validation_layers = true;
const char *enabled_validation_layers[] = {
	"VK_LAYER_KHRONOS_validation",
};
const int nenabled_validation_layers = sizeof(enabled_validation_layers)/sizeof(char *);

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
	if (enable_validation_layers) {
		create_info.enabledLayerCount = nenabled_validation_layers;
		create_info.ppEnabledLayerNames = enabled_validation_layers;
	} else {
		create_info.enabledLayerCount = 0;
	}

	VkInstance *instance = calloc(1, sizeof(VkInstance));
	VkResult result = vkCreateInstance(&create_info, NULL, instance);
	if (result != VK_SUCCESS) {
		fprintf(stderr, "unable to create instance.\n");
		exit(EXIT_FAILURE);
	}
	return instance;
}

void cleanup_vulkan(VkInstance *instance) {
	vkDestroyInstance(*instance, NULL);
}

void check_extensions() {
	uint32_t nextensions = 0;
	vkEnumerateInstanceExtensionProperties(NULL, &nextensions, NULL);
	VkExtensionProperties extensions[nextensions];
	vkEnumerateInstanceExtensionProperties(NULL, &nextensions, extensions);
	printf("nextensions: %d\n", nextensions);
	for (int i = 0; i < nextensions; i++) {
		printf("   extension: %-40s (0x%08X)\n", extensions[i].extensionName, extensions[i].specVersion);
	}
}

bool check_validation_layers() {
	uint32_t nlayers = 0;
	vkEnumerateInstanceLayerProperties(&nlayers, NULL);
	VkLayerProperties layers[nlayers];
	vkEnumerateInstanceLayerProperties(&nlayers, layers);
	printf("nlayers: %d\n", nlayers);
	for (int i = 0; i < nlayers; i++) {
		printf("   layer: %-35s (0x%08X) (0x%08X)\n", layers[i].layerName, layers[i].specVersion, layers[i].implementationVersion);
		printf("      desc: %s\n", layers[i].description);
	}
	for (int j = 0; j < nenabled_validation_layers; j++) {
		if (!has_layer(layers, nlayers, enabled_validation_layers[j])) {
			fprintf(stderr, "unable to locate validation layer \"%s\".\n", enabled_validation_layers[j]);
			return false;
		}
	}
	return true;
}

bool has_layer(VkLayerProperties *layers, int nlayers, const char *layer_name) {
	for (int i = 0; i < nlayers; i++) {
		if (strcmp(layers[i].layerName, layer_name) == 0) {
			return true;
		}
	}
	return false;
}
