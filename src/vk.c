#include "defs.h"
#include "vk.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void populate_debug_messanger_create_info(VkDebugUtilsMessengerCreateInfoEXT *create_info) {
	create_info->sType = VK_STRUCTURE_TYPE_DEBUG_UTILS_MESSENGER_CREATE_INFO_EXT;
	create_info->messageSeverity = VK_DEBUG_UTILS_MESSAGE_SEVERITY_VERBOSE_BIT_EXT | VK_DEBUG_UTILS_MESSAGE_SEVERITY_INFO_BIT_EXT | VK_DEBUG_UTILS_MESSAGE_SEVERITY_WARNING_BIT_EXT | VK_DEBUG_UTILS_MESSAGE_SEVERITY_ERROR_BIT_EXT;
	create_info->messageType = VK_DEBUG_UTILS_MESSAGE_TYPE_GENERAL_BIT_EXT | VK_DEBUG_UTILS_MESSAGE_TYPE_VALIDATION_BIT_EXT | VK_DEBUG_UTILS_MESSAGE_TYPE_PERFORMANCE_BIT_EXT;
	create_info->pfnUserCallback = debug_callback;
	create_info->pUserData = NULL; // optional.
}

VkDebugUtilsMessengerEXT * init_debug_messanger(VkInstance *instance) {
	VkDebugUtilsMessengerCreateInfoEXT *debug_messanger_create_info = calloc(1, sizeof(VkDebugUtilsMessengerCreateInfoEXT));
	populate_debug_messanger_create_info(debug_messanger_create_info);
	VkDebugUtilsMessengerEXT *debug_messenger = calloc(1, sizeof(VkDebugUtilsMessengerEXT));
	VkResult result = CreateDebugUtilsMessengerEXT(*instance, debug_messanger_create_info, NULL, debug_messenger);
	if (result!= VK_SUCCESS) {
		fprintf(stderr, "unable to register debug messanger (result=%d).\n", result);
	}
	return debug_messenger;
}

VkBool32 debug_callback(
	VkDebugUtilsMessageSeverityFlagBitsEXT messageSeverity,
	VkDebugUtilsMessageTypeFlagsEXT messageTypes,
	const VkDebugUtilsMessengerCallbackDataEXT *pCallbackData,
	void *pUserData) {
	if ((messageSeverity&VK_DEBUG_UTILS_MESSAGE_SEVERITY_ERROR_BIT_EXT) != 0) {
		fprintf(stderr, "validation layer (error): %s\n", pCallbackData->pMessage);
	}
	if ((messageSeverity&VK_DEBUG_UTILS_MESSAGE_SEVERITY_WARNING_BIT_EXT) != 0) {
		fprintf(stderr, "validation layer (warning): %s\n", pCallbackData->pMessage);
	}
	if ((messageSeverity&VK_DEBUG_UTILS_MESSAGE_SEVERITY_INFO_BIT_EXT) != 0) {
		fprintf(stderr, "validation layer (info): %s\n", pCallbackData->pMessage);
	}
	if ((messageSeverity&VK_DEBUG_UTILS_MESSAGE_SEVERITY_VERBOSE_BIT_EXT) != 0) {
		fprintf(stderr, "validation layer (verbose): %s\n", pCallbackData->pMessage);
	}
	return VK_FALSE;
}

VkResult CreateDebugUtilsMessengerEXT(
	VkInstance instance,
	const VkDebugUtilsMessengerCreateInfoEXT *pCreateInfo,
	const VkAllocationCallbacks *pAllocator,
	VkDebugUtilsMessengerEXT *pMessenger) {
	PFN_vkCreateDebugUtilsMessengerEXT func = (PFN_vkCreateDebugUtilsMessengerEXT)vkGetInstanceProcAddr(instance, "vkCreateDebugUtilsMessengerEXT");
	if (func != NULL) {
		return func(instance, pCreateInfo, pAllocator, pMessenger);
	}
	return VK_ERROR_EXTENSION_NOT_PRESENT;
}

void DestroyDebugUtilsMessengerEXT(
	VkInstance instance,
	VkDebugUtilsMessengerEXT messenger,
	const VkAllocationCallbacks *pAllocator) {
	PFN_vkDestroyDebugUtilsMessengerEXT func = (PFN_vkDestroyDebugUtilsMessengerEXT)vkGetInstanceProcAddr(instance, "vkDestroyDebugUtilsMessengerEXT");
	if (func != NULL) {
		func(instance, messenger, pAllocator);
	}
}

VkPhysicalDevice * init_device(VkInstance *instance) {
	// TODO: rank physical devices by score if more than one is present. E.g.
	// prefer dedicated graphics card with capability for larger textures.
	//
	// ref: https://vulkan-tutorial.com/en/Drawing_a_triangle/Setup/Physical_devices_and_queue_families#page_Base-device-suitability-checks

	// Get physical devices.
	uint32_t ndevices = 0;
	vkEnumeratePhysicalDevices(*instance, &ndevices, NULL);
	if (ndevices == 0) {
		fprintf(stderr, "unable to locate GPU physical device.\n");
		exit(EXIT_FAILURE);
	}
	VkPhysicalDevice *devices = calloc(ndevices, sizeof(VkPhysicalDevice));
	vkEnumeratePhysicalDevices(*instance, &ndevices, devices);
	printf("ndevices: %d\n", ndevices);
	for (int i = 0; i < ndevices; i++) {
		if (!is_suitable_defice(&devices[i])) {
			continue;
		}
		//printf("   device: %-40s (0x%08X)\n", devices[i].deviceName, devices[i].specVersion);
		return &devices[i];
	}
	fprintf(stderr, "unable to locate suitable GPU physical device.\n");
	exit(EXIT_FAILURE);
}

bool is_suitable_defice(VkPhysicalDevice *device) {
	// Get device properties.
	VkPhysicalDeviceProperties device_properties;
	vkGetPhysicalDeviceProperties(*device, &device_properties);
	printf("device_properties\n");
	printf("   apiVersion: 0x%08X\n", device_properties.apiVersion);
	printf("   driverVersion: 0x%08X\n", device_properties.driverVersion);
	printf("   vendorID: 0x%08X\n", device_properties.vendorID);
	printf("   deviceID: 0x%08X\n", device_properties.deviceID);
	printf("   deviceType: %d\n", device_properties.deviceType);
	printf("   deviceType: %s\n", device_properties.deviceName);
	printf("   pipelineCacheUUID: ");
	for (int i = 0; i < VK_UUID_SIZE; i++) {
		if (i > 0) {
			printf(" ");
		}
		printf("%02X", device_properties.pipelineCacheUUID[i]);
	}
	printf("\n");
	//printf("   deviceType: %v\n", device_properties.limits);
	//printf("   sparseProperties: %v\n", device_properties.sparseProperties;

	// Get device features.
	VkPhysicalDeviceFeatures device_features;
	vkGetPhysicalDeviceFeatures(*device, &device_features);

	printf("device_features\n");
	printf("   robustBufferAccess: %d\n", device_features.robustBufferAccess);
	printf("   fullDrawIndexUint32: %d\n", device_features.fullDrawIndexUint32);
	printf("   imageCubeArray: %d\n", device_features.imageCubeArray);
	printf("   independentBlend: %d\n", device_features.independentBlend);
	printf("   geometryShader: %d\n", device_features.geometryShader);
	printf("   tessellationShader: %d\n", device_features.tessellationShader);
	printf("   sampleRateShading: %d\n", device_features.sampleRateShading);
	printf("   dualSrcBlend: %d\n", device_features.dualSrcBlend);
	printf("   logicOp: %d\n", device_features.logicOp);
	printf("   multiDrawIndirect: %d\n", device_features.multiDrawIndirect);
	printf("   drawIndirectFirstInstance: %d\n", device_features.drawIndirectFirstInstance);
	printf("   depthClamp: %d\n", device_features.depthClamp);
	printf("   depthBiasClamp: %d\n", device_features.depthBiasClamp);
	printf("   fillModeNonSolid: %d\n", device_features.fillModeNonSolid);
	printf("   depthBounds: %d\n", device_features.depthBounds);
	printf("   wideLines: %d\n", device_features.wideLines);
	printf("   largePoints: %d\n", device_features.largePoints);
	printf("   alphaToOne: %d\n", device_features.alphaToOne);
	printf("   multiViewport: %d\n", device_features.multiViewport);
	printf("   samplerAnisotropy: %d\n", device_features.samplerAnisotropy);
	printf("   textureCompressionETC2: %d\n", device_features.textureCompressionETC2);
	printf("   textureCompressionASTC_LDR: %d\n", device_features.textureCompressionASTC_LDR);
	printf("   textureCompressionBC: %d\n", device_features.textureCompressionBC);
	printf("   occlusionQueryPrecise: %d\n", device_features.occlusionQueryPrecise);
	printf("   pipelineStatisticsQuery: %d\n", device_features.pipelineStatisticsQuery);
	printf("   vertexPipelineStoresAndAtomics: %d\n", device_features.vertexPipelineStoresAndAtomics);
	printf("   fragmentStoresAndAtomics: %d\n", device_features.fragmentStoresAndAtomics);
	printf("   shaderTessellationAndGeometryPointSize: %d\n", device_features.shaderTessellationAndGeometryPointSize);
	printf("   shaderImageGatherExtended: %d\n", device_features.shaderImageGatherExtended);
	printf("   shaderStorageImageExtendedFormats: %d\n", device_features.shaderStorageImageExtendedFormats);
	printf("   shaderStorageImageMultisample: %d\n", device_features.shaderStorageImageMultisample);
	printf("   shaderStorageImageReadWithoutFormat: %d\n", device_features.shaderStorageImageReadWithoutFormat);
	printf("   shaderStorageImageWriteWithoutFormat: %d\n", device_features.shaderStorageImageWriteWithoutFormat);
	printf("   shaderUniformBufferArrayDynamicIndexing: %d\n", device_features.shaderUniformBufferArrayDynamicIndexing);
	printf("   shaderSampledImageArrayDynamicIndexing: %d\n", device_features.shaderSampledImageArrayDynamicIndexing);
	printf("   shaderStorageBufferArrayDynamicIndexing: %d\n", device_features.shaderStorageBufferArrayDynamicIndexing);
	printf("   shaderStorageImageArrayDynamicIndexing: %d\n", device_features.shaderStorageImageArrayDynamicIndexing);
	printf("   shaderClipDistance: %d\n", device_features.shaderClipDistance);
	printf("   shaderCullDistance: %d\n", device_features.shaderCullDistance);
	printf("   shaderFloat64: %d\n", device_features.shaderFloat64);
	printf("   shaderInt64: %d\n", device_features.shaderInt64);
	printf("   shaderInt16: %d\n", device_features.shaderInt16);
	printf("   shaderResourceResidency: %d\n", device_features.shaderResourceResidency);
	printf("   shaderResourceMinLod: %d\n", device_features.shaderResourceMinLod);
	printf("   sparseBinding: %d\n", device_features.sparseBinding);
	printf("   sparseResidencyBuffer: %d\n", device_features.sparseResidencyBuffer);
	printf("   sparseResidencyImage2D: %d\n", device_features.sparseResidencyImage2D);
	printf("   sparseResidencyImage3D: %d\n", device_features.sparseResidencyImage3D);
	printf("   sparseResidency2Samples: %d\n", device_features.sparseResidency2Samples);
	printf("   sparseResidency4Samples: %d\n", device_features.sparseResidency4Samples);
	printf("   sparseResidency8Samples: %d\n", device_features.sparseResidency8Samples);
	printf("   sparseResidency16Samples: %d\n", device_features.sparseResidency16Samples);
	printf("   sparseResidencyAliased: %d\n", device_features.sparseResidencyAliased);
	printf("   variableMultisampleRate: %d\n", device_features.variableMultisampleRate);
	printf("   inheritedQueries: %d\n", device_features.inheritedQueries);

	return true;
}
