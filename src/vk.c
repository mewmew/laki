#include "defs.h"
#include "vk.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

const char *REQUIRED_EXTENSIONS[] = {
	"VK_EXT_debug_utils",
};
const int NREQUIRED_EXTENSIONS = sizeof(REQUIRED_EXTENSIONS)/sizeof(char *);

const char *REQUIRED_LAYERS[] = {
	"VK_LAYER_KHRONOS_validation",
};
const int NREQUIRED_LAYERS = sizeof(REQUIRED_LAYERS)/sizeof(char *);

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

	uint32_t nenabled_extensions = 0;
	const char **enabled_extensions = get_extensions(&nenabled_extensions);

	uint32_t nenabled_layers = 0;
	const char **enabled_layers = get_layers(&nenabled_layers);

	VkInstanceCreateInfo create_info = {};
	create_info.sType = VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO;
	create_info.flags = 0; // reserved.
	create_info.pApplicationInfo = &app_info;
	create_info.enabledExtensionCount = nenabled_extensions;
	create_info.ppEnabledExtensionNames = enabled_extensions;
	create_info.enabledLayerCount = nenabled_layers;
	create_info.ppEnabledLayerNames = enabled_layers;

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

const char ** get_extensions(uint32_t *pnenabled_extensions) {
	// Get supported extensions.
	uint32_t nextensions = 0;
	vkEnumerateInstanceExtensionProperties(NULL, &nextensions, NULL);
	VkExtensionProperties extensions[nextensions];
	vkEnumerateInstanceExtensionProperties(NULL, &nextensions, extensions);
	printf("nextensions: %d\n", nextensions);
	for (int i = 0; i < nextensions; i++) {
		printf("   extension: %-40s (0x%08X)\n", extensions[i].extensionName, extensions[i].specVersion);
	}

	// Get required extensions for GLFW.
	uint32_t nglfw_required_extensions = 0;
	const char **glfw_required_extensions = glfwGetRequiredInstanceExtensions(&nglfw_required_extensions);
	printf("nglfw_required_extensions: %d\n", nglfw_required_extensions);
	for (int i = 0; i < nglfw_required_extensions; i++) {
		printf("   glfw_required_extension: %s\n", glfw_required_extensions[i]);
	}
	printf("NREQUIRED_EXTENSIONS: %d\n", NREQUIRED_EXTENSIONS);
	for (int i = 0; i < NREQUIRED_EXTENSIONS; i++) {
		printf("   REQUIRED_EXTENSION: %s\n", REQUIRED_EXTENSIONS[i]);
	}

	// Check required extensions.
	uint32_t nenabled_extensions = 0;
	uint32_t ntotal_required_extensions = nglfw_required_extensions + NREQUIRED_EXTENSIONS;
	const char **enabled_extensions = calloc(ntotal_required_extensions, sizeof(char *));
	for (int j = 0; j < nglfw_required_extensions; j++) {
		if (!has_extension(extensions, nextensions, glfw_required_extensions[j])) {
			fprintf(stderr, "unable to locate extension \"%s\".\n", glfw_required_extensions[j]);
			continue; // ignore unsupported extension.
		}
		enabled_extensions[nenabled_extensions] = glfw_required_extensions[j];
		nenabled_extensions++;
	}
	for (int j = 0; j < NREQUIRED_EXTENSIONS; j++) {
		if (!has_extension(extensions, nextensions, REQUIRED_EXTENSIONS[j])) {
			fprintf(stderr, "unable to locate extension \"%s\".\n", REQUIRED_EXTENSIONS[j]);
			continue; // ignore unsupported extension.
		}
		enabled_extensions[nenabled_extensions] = REQUIRED_EXTENSIONS[j];
		nenabled_extensions++;
	}

	return enabled_extensions;
}

const char ** get_layers(uint32_t *pnenabled_layers) {
	// Get supported layers.
	uint32_t nlayers = 0;
	vkEnumerateInstanceLayerProperties(&nlayers, NULL);
	VkLayerProperties layers[nlayers];
	vkEnumerateInstanceLayerProperties(&nlayers, layers);
	printf("nlayers: %d\n", nlayers);
	for (int i = 0; i < nlayers; i++) {
		printf("   layer: %-35s (0x%08X) (0x%08X)\n", layers[i].layerName, layers[i].specVersion, layers[i].implementationVersion);
		printf("      desc: %s\n", layers[i].description);
	}
	printf("NREQUIRED_LAYERS: %d\n", NREQUIRED_LAYERS);
	for (int i = 0; i < NREQUIRED_LAYERS; i++) {
		printf("   REQUIRED_LAYER: %s\n", REQUIRED_LAYERS[i]);
	}

	// Check required layers.
	uint32_t nenabled_layers = 0;
	uint32_t ntotal_required_layers = NREQUIRED_LAYERS;
	const char **enabled_layers = calloc(ntotal_required_layers, sizeof(char *));
	for (int j = 0; j < NREQUIRED_LAYERS; j++) {
		if (!has_layer(layers, nlayers, REQUIRED_LAYERS[j])) {
			fprintf(stderr, "unable to locate layer \"%s\".\n", REQUIRED_LAYERS[j]);
			continue; // ignore unsupported layer.
		}
		enabled_layers[nenabled_layers] = REQUIRED_LAYERS[j];
		nenabled_layers++;
	}

	return enabled_layers;
}

bool has_extension(VkExtensionProperties *extensions, int nextensions, const char *extension_name) {
	for (int i = 0; i < nextensions; i++) {
		if (strcmp(extensions[i].extensionName, extension_name) == 0) {
			return true;
		}
	}
	return false;
}

bool has_layer(VkLayerProperties *layers, int nlayers, const char *layer_name) {
	for (int i = 0; i < nlayers; i++) {
		if (strcmp(layers[i].layerName, layer_name) == 0) {
			return true;
		}
	}
	return false;
}
