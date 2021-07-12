// TODO: continue at https://vulkan-tutorial.com/en/Drawing_a_triangle/Drawing/Framebuffers

// refs:
// * Graphics pipeline overview: https://vulkan-tutorial.com/en/Drawing_a_triangle/Graphics_pipeline_basics/Introduction

package vk

// #define GLFW_INCLUDE_VULKAN
// #include <GLFW/glfw3.h>
//
// #include "callback.h"
// #include "invoke.h"
// #include "malloc.h"
//
//#cgo pkg-config: glfw3
//#cgo pkg-config: vulkan
import "C"

import (
	"io/ioutil"
	"log"
	"os"
	"sort"
	"unsafe"

	"github.com/kr/pretty"
	"github.com/mewkiz/pkg/term"
	"github.com/pkg/errors"
)

var (
	// dbg is a logger with the "vk:" prefix which logs debug messages to
	// standard error.
	dbg = log.New(os.Stderr, term.MagentaBold("vk:")+" ", 0)
	// warn is a logger with the "vk:" prefix which logs warning messages to
	// standard error.
	warn = log.New(os.Stderr, term.RedBold("vk:")+" ", log.Lshortfile)
)

func init() {
	if !debug {
		dbg.SetOutput(ioutil.Discard)
	}
}

// Enable debug output.
const debug = true

// Extensions required by Vulkan instance.
var RequiredInstanceExtensions = []string{
	C.VK_EXT_DEBUG_UTILS_EXTENSION_NAME, // "VK_EXT_debug_utils"
}

// Validation layers required by Vulkan instance.
var RequiredLayers = []string{
	"VK_LAYER_KHRONOS_validation",
}

// Extensions required by device.
var RequiredDeviceExtensions = []string{
	C.VK_KHR_SWAPCHAIN_EXTENSION_NAME, // "VK_KHR_swapchain"
}

func Init() error {
	app := newApp()
	app.win = InitWindow()
	defer CleanupWindow(app.win)
	if err := InitVulkan(app); err != nil {
		return errors.WithStack(err)
	}
	defer CleanupVulkan(app)

	EventLoop(app.win)
	return nil
}

func InitVulkan(app *App) error {
	instance, err := createInstance()
	if err != nil {
		return errors.WithStack(err)
	}
	app.instance = instance
	debugMessanger, err := initDebugMessanger(app.instance)
	if err != nil {
		return errors.WithStack(err)
	}
	app.debugMessanger = debugMessanger
	surface, err := initSurface(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.surface = surface
	physicalDevice, err := initPhysicalDevice(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.physicalDevice = physicalDevice
	swapchainSupportInfo := getSwapchainSupportInfo(app, physicalDevice)
	app.swapchainSupportInfo = swapchainSupportInfo
	pretty.Println("   swapchainSupportInfo:", swapchainSupportInfo)
	device, err := initDevice(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.device = device
	initQueues(app)
	pretty.Println("app.QueueFamilyIndices:", app.QueueFamilyIndices)
	swapchain, err := initSwapchain(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.swapchain = swapchain
	pretty.Println("   swapchain:", swapchain)
	app.swapchainImgs = getSwapchainImgs(app)
	swapchainImgViews, err := initSwapchainImgViews(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.swapchainImgViews = swapchainImgViews
	renderPass, err := initRenderPass(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.renderPass = renderPass
	graphicsPipelines, err := initGraphicsPipeline(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.graphicsPipelines = graphicsPipelines
	pretty.Println("   graphicsPipelines:", graphicsPipelines)
	return nil
}

func CleanupVulkan(app *App) {
	for _, graphicsPipeline := range app.graphicsPipelines {
		C.vkDestroyPipeline(*app.device, graphicsPipeline, nil)
	}
	C.vkDestroyPipelineLayout(*app.device, *app.pipelineLayout, nil)
	C.vkDestroyRenderPass(*app.device, *app.renderPass, nil)
	C.vkDestroyShaderModule(*app.device, *app.fragmentShaderModule, nil)
	C.vkDestroyShaderModule(*app.device, *app.vertexShaderModule, nil)
	for i := range app.swapchainImgViews {
		C.vkDestroyImageView(*app.device, app.swapchainImgViews[i], nil)
	}
	C.vkDestroySwapchainKHR(*app.device, *app.swapchain, nil)
	C.vkDestroyDevice(*app.device, nil)
	app.physicalDevice = nil
	DestroyDebugUtilsMessengerEXT(*app.instance, *app.debugMessanger, nil)
	C.vkDestroySurfaceKHR(*app.instance, *app.surface, nil)
	C.vkDestroyInstance(*app.instance, nil)
}

func createInstance() (*C.VkInstance, error) {
	appInfo := C.VkApplicationInfo{
		sType:              C.VK_STRUCTURE_TYPE_APPLICATION_INFO,
		pApplicationName:   C.CString(AppTitle),
		applicationVersion: VK_MAKE_API_VERSION(0, 1, 0, 0),
		pEngineName:        C.CString("No Engine"),
		engineVersion:      VK_MAKE_API_VERSION(0, 1, 0, 0),
		apiVersion:         C.VK_API_VERSION_1_0,
	}

	enabledInstanceExtensions := getInstanceExtensions()
	dbg.Println("nenabledInstanceExtensions:", len(enabledInstanceExtensions))
	for _, enabledInstanceExtension := range enabledInstanceExtensions {
		dbg.Println("   enabledInstanceExtension:", enabledInstanceExtension)
	}

	enabledLayers := getLayers()
	dbg.Println("nenabledLayers:", len(enabledLayers))
	for _, enabledLayer := range enabledLayers {
		dbg.Println("   enabledLayer:", enabledLayer)
	}

	createInfo := C.new_VkInstanceCreateInfo()
	createInfo.sType = C.VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO
	createInfo.flags = 0 // reserved.
	createInfo.pApplicationInfo = &appInfo
	createInfo.enabledExtensionCount = C.uint32_t(len(enabledInstanceExtensions))
	createInfo.ppEnabledExtensionNames = getCStringSlice(enabledInstanceExtensions)
	createInfo.enabledLayerCount = C.uint32_t(len(enabledLayers))
	createInfo.ppEnabledLayerNames = getCStringSlice(enabledLayers)

	debugMessangerCreateInfo := C.new_VkDebugUtilsMessengerCreateInfoEXT()
	populateDebugMessangerCreateInfo(debugMessangerCreateInfo)
	createInfo.pNext = unsafe.Pointer(debugMessangerCreateInfo)

	instance := C.new_VkInstance()
	result := C.vkCreateInstance(createInfo, nil, instance)
	if result != C.VK_SUCCESS {
		return nil, errors.Errorf("unable to create Vulkan instance")
	}
	return instance, nil
}

func getInstanceExtensions() []string {
	// Get supported instance extensions.
	var ninstanceExtensions C.uint32_t
	C.vkEnumerateInstanceExtensionProperties(nil, &ninstanceExtensions, nil)
	instanceExtensions := make([]C.VkExtensionProperties, int(ninstanceExtensions))
	C.vkEnumerateInstanceExtensionProperties(nil, &ninstanceExtensions, &instanceExtensions[0])
	dbg.Println("ninstanceExtensions:", len(instanceExtensions))
	var instanceExtensionNames []string
	for _, instanceExtension := range instanceExtensions {
		instanceExtensionName := C.GoString(&instanceExtension.extensionName[0])
		dbg.Println("   instanceExtension:", instanceExtensionName)
		instanceExtensionNames = append(instanceExtensionNames, instanceExtensionName)
	}

	// Get required instance extensions for GLFW.
	var nglfwRequiredInstanceExtensions C.uint32_t
	_glfwRequiredInstanceExtensions := C.glfwGetRequiredInstanceExtensions(&nglfwRequiredInstanceExtensions)
	glfwRequiredInstanceExtensions := getStringSlice(unsafe.Pointer(_glfwRequiredInstanceExtensions), int(nglfwRequiredInstanceExtensions))
	dbg.Println("nglfwRequiredInstanceExtensions:", len(glfwRequiredInstanceExtensions))
	for _, glfwRequiredInstanceExtension := range glfwRequiredInstanceExtensions {
		dbg.Println("   glfwRequiredInstanceExtension:", glfwRequiredInstanceExtension)
	}

	// Get required instance extensions by user.
	dbg.Println("nrequiredInstanceExtensions:", len(RequiredInstanceExtensions))
	for _, requiredInstanceExtension := range RequiredInstanceExtensions {
		dbg.Println("   requiredInstanceExtension:", requiredInstanceExtension)
	}

	// Check required instance extensions.
	var enabledInstanceExtensions []string
	for _, glfwRequiredInstanceExtension := range glfwRequiredInstanceExtensions {
		if !contains(instanceExtensionNames, glfwRequiredInstanceExtension) {
			warn.Printf("unable to locate required extension %q", glfwRequiredInstanceExtension)
			continue
		}
		enabledInstanceExtensions = append(enabledInstanceExtensions, glfwRequiredInstanceExtension)
	}
	for _, requiredInstanceExtension := range RequiredInstanceExtensions {
		if !contains(instanceExtensionNames, requiredInstanceExtension) {
			warn.Printf("unable to locate required extension %q", requiredInstanceExtension)
			continue
		}
		enabledInstanceExtensions = append(enabledInstanceExtensions, requiredInstanceExtension)
	}

	return enabledInstanceExtensions
}

func getLayers() []string {
	// Get supported layers.
	var nlayers C.uint32_t
	C.vkEnumerateInstanceLayerProperties(&nlayers, nil)
	layers := make([]C.VkLayerProperties, int(nlayers))
	C.vkEnumerateInstanceLayerProperties(&nlayers, &layers[0])
	dbg.Println("nlayers:", len(layers))
	var layerNames []string
	for _, layer := range layers {
		layerName := C.GoString(&layer.layerName[0])
		layerDesc := C.GoString(&layer.description[0])
		dbg.Println("   layer:", layerName)
		dbg.Println("      desc:", layerDesc)
		layerNames = append(layerNames, layerName)
	}

	// Get required layers by user.
	dbg.Println("nrequiredLayers:", len(RequiredLayers))
	for _, requiredLayer := range RequiredLayers {
		dbg.Println("   requiredLayer:", requiredLayer)
	}

	// Check required layers.
	var enabledLayers []string
	for _, requiredLayer := range RequiredLayers {
		if !contains(layerNames, requiredLayer) {
			warn.Printf("unable to locate required layer %q", requiredLayer)
			continue
		}
		enabledLayers = append(enabledLayers, requiredLayer)
	}

	return enabledLayers
}

func getDeviceExtensions(physicalDevice *C.VkPhysicalDevice) []string {
	// Get supported device extensions.
	var ndeviceExtensions C.uint32_t
	C.vkEnumerateDeviceExtensionProperties(*physicalDevice, nil, &ndeviceExtensions, nil)
	deviceExtensions := make([]C.VkExtensionProperties, int(ndeviceExtensions))
	C.vkEnumerateDeviceExtensionProperties(*physicalDevice, nil, &ndeviceExtensions, &deviceExtensions[0])
	dbg.Println("ndeviceExtensions:", len(deviceExtensions))
	var deviceExtensionNames []string
	for _, deviceExtension := range deviceExtensions {
		deviceExtensionName := C.GoString(&deviceExtension.extensionName[0])
		dbg.Println("   deviceExtension:", deviceExtensionName)
		deviceExtensionNames = append(deviceExtensionNames, deviceExtensionName)
	}

	// Get required device extensions by user.
	dbg.Println("nrequiredDeviceExtensions:", len(RequiredDeviceExtensions))
	for _, requiredDeviceExtension := range RequiredDeviceExtensions {
		dbg.Println("   requiredDeviceExtension:", requiredDeviceExtension)
	}

	// Check required device extensions.
	var enabledDeviceExtensions []string
	for _, requiredDeviceExtension := range RequiredDeviceExtensions {
		if !contains(deviceExtensionNames, requiredDeviceExtension) {
			warn.Printf("unable to locate required extension %q", requiredDeviceExtension)
			continue
		}
		enabledDeviceExtensions = append(enabledDeviceExtensions, requiredDeviceExtension)
	}

	return enabledDeviceExtensions
}

func initPhysicalDevice(app *App) (*C.VkPhysicalDevice, error) {
	// TODO: rank physical devices by score if more than one is present. E.g.
	// prefer dedicated graphics card with capability for larger textures.
	//
	// ref: https://vulkan-tutorial.com/en/Drawing_a_triangle/Setup/Physical_devices_and_queue_families#page_Base-device-suitability-checks

	// Get physical devices.
	var nphysicalDevices C.uint32_t
	C.vkEnumeratePhysicalDevices(*app.instance, &nphysicalDevices, nil)
	if nphysicalDevices == 0 {
		return nil, errors.Errorf("unable to locate physical device (GPU)")
	}
	if nphysicalDevices > 1 {
		warn.Printf("multiple (%d) physical device (GPU) located; support for ranking physical devices not yet implemented", nphysicalDevices)
	}
	physicalDevices := make([]C.VkPhysicalDevice, int(nphysicalDevices))
	C.vkEnumeratePhysicalDevices(*app.instance, &nphysicalDevices, &physicalDevices[0])
	dbg.Println("nphysicalDevices:", len(physicalDevices))
	for _, physicalDevice := range physicalDevices {
		if !isSuitablePhysicalDevice(app, &physicalDevice) {
			continue
		}
		_physicalDevice := C.new_VkPhysicalDevice()
		*_physicalDevice = physicalDevice // allocate pointer on C heap.
		return _physicalDevice, nil
	}
	return nil, errors.Errorf("unable to locate suitable physical device (GPU)")
}

func isSuitablePhysicalDevice(app *App, physicalDevice *C.VkPhysicalDevice) bool {
	// Get device properties.
	var deviceProperties C.VkPhysicalDeviceProperties
	C.vkGetPhysicalDeviceProperties(*physicalDevice, &deviceProperties)
	deviceName := C.GoString(&deviceProperties.deviceName[0])
	dbg.Println("   deviceName:", deviceName)
	pretty.Println("   deviceProperties:", deviceProperties)

	// Get device features.
	var deviceFeatures C.VkPhysicalDeviceFeatures
	C.vkGetPhysicalDeviceFeatures(*physicalDevice, &deviceFeatures)
	pretty.Println("   deviceFeatures:", deviceFeatures)

	// Find queue which supports graphics operations.
	queueFamilies := getQueueFamilies(physicalDevice)
	dbg.Println("nqueueFamilies:", len(queueFamilies))
	for _, queueFamily := range queueFamilies {
		pretty.Println("   queueFamily:", queueFamily)
	}
	if _, ok := findQueueWithFlag(queueFamilies, C.VK_QUEUE_GRAPHICS_BIT); !ok {
		return false
	}
	if _, ok := findQueueWithPresentSupport(physicalDevice, app.surface, queueFamilies); !ok {
		return false
	}

	if !hasDeviceExtensionSupport(physicalDevice) {
		return false
	}

	swapchainSupportInfo := getSwapchainSupportInfo(app, physicalDevice)
	if len(swapchainSupportInfo.surfaceFormats) == 0 {
		return false
	}
	if len(swapchainSupportInfo.presentModes) == 0 {
		return false
	}

	return true
}

func hasDeviceExtensionSupport(physicalDevice *C.VkPhysicalDevice) bool {
	var ndeviceExtensions C.uint32_t
	C.vkEnumerateDeviceExtensionProperties(*physicalDevice, nil, &ndeviceExtensions, nil)
	deviceExtensions := make([]C.VkExtensionProperties, int(ndeviceExtensions))
	C.vkEnumerateDeviceExtensionProperties(*physicalDevice, nil, &ndeviceExtensions, &deviceExtensions[0])
	dbg.Println("ndeviceExtensions:", len(deviceExtensions))
	// Check that all required device extensions are present.
	m := make(map[string]bool)
	for _, requiredDeviceExtension := range RequiredDeviceExtensions {
		m[requiredDeviceExtension] = true
	}
	for _, deviceExtension := range deviceExtensions {
		deviceExtensionName := C.GoString(&deviceExtension.extensionName[0])
		dbg.Println("   deviceExtensionName:", deviceExtensionName)
		pretty.Println("   deviceExtension:", deviceExtension)
		delete(m, deviceExtensionName)
	}
	if len(m) > 1 {
		warn.Println("   missing required device extensions", m)
	}
	return len(m) == 0
}

func findQueueWithFlag(queueFamilies []C.VkQueueFamilyProperties, queueFlags C.VkQueueFlags) (int, bool) {
	for queueFamilyIndex, queueFamily := range queueFamilies {
		if queueFamily.queueFlags&queueFlags == queueFlags {
			return queueFamilyIndex, true
		}
	}
	return 0, false
}

func findQueueWithPresentSupport(physicalDevice *C.VkPhysicalDevice, surface *C.VkSurfaceKHR, queueFamilies []C.VkQueueFamilyProperties) (int, bool) {
	for queueFamilyIndex := range queueFamilies {
		var presentSupport C.VkBool32
		C.vkGetPhysicalDeviceSurfaceSupportKHR(*physicalDevice, C.uint(queueFamilyIndex), *surface, &presentSupport)
		if presentSupport == C.VK_TRUE {
			return queueFamilyIndex, true
		}
	}
	return 0, false
}

func initDebugMessanger(instance *C.VkInstance) (*C.VkDebugUtilsMessengerEXT, error) {
	debugMessangerCreateInfo := C.new_VkDebugUtilsMessengerCreateInfoEXT()
	populateDebugMessangerCreateInfo(debugMessangerCreateInfo)
	debugMessenger := C.new_VkDebugUtilsMessengerEXT()
	result := CreateDebugUtilsMessengerEXT(*instance, debugMessangerCreateInfo, nil, debugMessenger)
	if result != C.VK_SUCCESS {
		return nil, errors.Errorf("unable to register debug messanger (result=%d)", result)
	}
	return debugMessenger, nil
}

func CreateDebugUtilsMessengerEXT(instance C.VkInstance, pCreateInfo *C.VkDebugUtilsMessengerCreateInfoEXT, pAllocator *C.VkAllocationCallbacks, pMessenger *C.VkDebugUtilsMessengerEXT) C.VkResult {
	fn := C.vkGetInstanceProcAddr(instance, C.CString("vkCreateDebugUtilsMessengerEXT"))
	if fn == nil {
		return C.VK_ERROR_EXTENSION_NOT_PRESENT
	}
	return C.invoke_CreateDebugUtilsMessengerEXT(fn, instance, pCreateInfo, pAllocator, pMessenger)
}

func DestroyDebugUtilsMessengerEXT(instance C.VkInstance, messenger C.VkDebugUtilsMessengerEXT, pAllocator *C.VkAllocationCallbacks) {
	fn := C.vkGetInstanceProcAddr(instance, C.CString("vkDestroyDebugUtilsMessengerEXT"))
	if fn == nil {
		return
	}
	C.invoke_DestroyDebugUtilsMessengerEXT(fn, instance, messenger, pAllocator)
}

func populateDebugMessangerCreateInfo(createInfo *C.VkDebugUtilsMessengerCreateInfoEXT) {
	createInfo.sType = C.VK_STRUCTURE_TYPE_DEBUG_UTILS_MESSENGER_CREATE_INFO_EXT
	createInfo.messageSeverity = C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_VERBOSE_BIT_EXT | C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_INFO_BIT_EXT | C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_WARNING_BIT_EXT | C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_ERROR_BIT_EXT
	createInfo.messageType = C.VK_DEBUG_UTILS_MESSAGE_TYPE_GENERAL_BIT_EXT | C.VK_DEBUG_UTILS_MESSAGE_TYPE_VALIDATION_BIT_EXT | C.VK_DEBUG_UTILS_MESSAGE_TYPE_PERFORMANCE_BIT_EXT
	//createInfo.pfnUserCallback = (C.PFN_vkDebugUtilsMessengerCallbackEXT)(unsafe.Pointer(C.debug_callback))
	createInfo.pfnUserCallback = (C.PFN_vkDebugUtilsMessengerCallbackEXT)(unsafe.Pointer(C.debugCallback))
	createInfo.pUserData = nil // optional.
}

func getQueueFamilies(device *C.VkPhysicalDevice) []C.VkQueueFamilyProperties {
	var nqueueFamilies C.uint32_t
	C.vkGetPhysicalDeviceQueueFamilyProperties(*device, &nqueueFamilies, nil)
	queueFamilies := make([]C.VkQueueFamilyProperties, int(nqueueFamilies))
	C.vkGetPhysicalDeviceQueueFamilyProperties(*device, &nqueueFamilies, &queueFamilies[0])
	return queueFamilies
}

func initDevice(app *App) (*C.VkDevice, error) {
	queueFamilies := getQueueFamilies(app.physicalDevice)

	// Graphics queue.
	graphicsQueueFamilyIndex, ok := findQueueWithFlag(queueFamilies, C.VK_QUEUE_GRAPHICS_BIT)
	if !ok {
		return nil, errors.Errorf("unable to locate queue family with support for graphics operations")
	}
	app.graphicsQueueFamilyIndex = graphicsQueueFamilyIndex

	// Present queue.
	presentQueueFamilyIndex, ok := findQueueWithPresentSupport(app.physicalDevice, app.surface, queueFamilies)
	if !ok {
		return nil, errors.Errorf("unable to locate queue family with support for present operations")
	}
	app.presentQueueFamilyIndex = presentQueueFamilyIndex

	// Create queues.
	var queueCreateInfos []C.VkDeviceQueueCreateInfo
	// Find unique indices.
	queueFamilyIndices := unique(app.QueueFamilyIndices.Indices()...)
	for _, queueFamilyIndex := range queueFamilyIndices {
		const queueCount = 1
		queuePriorities := [queueCount]C.float{1.0}
		queueCreateInfo := C.new_VkDeviceQueueCreateInfo()
		queueCreateInfo.sType = C.VK_STRUCTURE_TYPE_DEVICE_QUEUE_CREATE_INFO
		queueCreateInfo.queueFamilyIndex = C.uint(queueFamilyIndex)
		queueCreateInfo.queueCount = queueCount
		queueCreateInfo.pQueuePriorities = &queuePriorities[0]
		queueCreateInfos = append(queueCreateInfos, *queueCreateInfo)
	}

	enabledFeatures := C.new_VkPhysicalDeviceFeatures()
	// TODO: enable device features here when needed.

	enabledDeviceExtensions := getDeviceExtensions(app.physicalDevice)
	dbg.Println("nenabledDeviceExtensions:", len(enabledDeviceExtensions))
	for _, enabledDeviceExtension := range enabledDeviceExtensions {
		dbg.Println("   enabledDeviceExtension:", enabledDeviceExtension)
	}

	createInfo := C.new_VkDeviceCreateInfo()
	createInfo.sType = C.VK_STRUCTURE_TYPE_DEVICE_CREATE_INFO
	createInfo.queueCreateInfoCount = C.uint(len(queueCreateInfos))
	createInfo.pQueueCreateInfos = &queueCreateInfos[0]
	createInfo.enabledLayerCount = 0 // ignored by recent version of Vulkan.
	createInfo.enabledExtensionCount = C.uint32_t(len(enabledDeviceExtensions))
	createInfo.ppEnabledExtensionNames = getCStringSlice(enabledDeviceExtensions)
	createInfo.pEnabledFeatures = enabledFeatures

	device := C.new_VkDevice()
	if result := C.vkCreateDevice(*app.physicalDevice, createInfo, nil, device); result != C.VK_SUCCESS {
		return nil, errors.Errorf("unable to create device (result=%d)", result)
	}
	return device, nil
}

func initQueues(app *App) {
	// Graphics queue.
	graphicsQueue := C.new_VkQueue()
	C.vkGetDeviceQueue(*app.device, C.uint(app.graphicsQueueFamilyIndex), 0, graphicsQueue)
	app.graphicsQueue = graphicsQueue
	// Present queue.
	presentQueue := C.new_VkQueue()
	C.vkGetDeviceQueue(*app.device, C.uint(app.presentQueueFamilyIndex), 0, presentQueue)
	app.presentQueue = presentQueue
}

func initSurface(app *App) (*C.VkSurfaceKHR, error) {
	surface := C.new_VkSurfaceKHR()
	if result := C.glfwCreateWindowSurface(*app.instance, app.win, nil, surface); result != C.VK_SUCCESS {
		return nil, errors.Errorf("unable to create window surface (result=%d)", result)
	}
	return surface, nil
}

func getSwapchainSupportInfo(app *App, physicalDevice *C.VkPhysicalDevice) *SwapchainSupportInfo {
	// Get surface capabilities.
	swapchainSupportInfo := &SwapchainSupportInfo{}
	var surfaceCapabilities C.VkSurfaceCapabilitiesKHR
	C.vkGetPhysicalDeviceSurfaceCapabilitiesKHR(*physicalDevice, *app.surface, &surfaceCapabilities)
	swapchainSupportInfo.surfaceCapabilities = &surfaceCapabilities

	// Get surface formats.
	var nsurfaceFormats C.uint32_t
	C.vkGetPhysicalDeviceSurfaceFormatsKHR(*physicalDevice, *app.surface, &nsurfaceFormats, nil)
	surfaceFormats := make([]C.VkSurfaceFormatKHR, int(nsurfaceFormats))
	C.vkGetPhysicalDeviceSurfaceFormatsKHR(*physicalDevice, *app.surface, &nsurfaceFormats, &surfaceFormats[0])
	swapchainSupportInfo.surfaceFormats = surfaceFormats

	// Get present modes.
	var npresentModes C.uint32_t
	C.vkGetPhysicalDeviceSurfacePresentModesKHR(*physicalDevice, *app.surface, &npresentModes, nil)
	presentModes := make([]C.VkPresentModeKHR, int(npresentModes))
	C.vkGetPhysicalDeviceSurfacePresentModesKHR(*physicalDevice, *app.surface, &npresentModes, &presentModes[0])
	swapchainSupportInfo.presentModes = presentModes

	return swapchainSupportInfo
}

func chooseSwapExtent(app *App, surfaceCapabilities *C.VkSurfaceCapabilitiesKHR) C.VkExtent2D {
	if surfaceCapabilities.currentExtent.width == C.UINT32_MAX || surfaceCapabilities.currentExtent.height == C.UINT32_MAX {
		var width, height C.int
		C.glfwGetFramebufferSize(app.win, &width, &height)
		actualExtent := C.VkExtent2D{
			width:  C.uint(clamp(int(width), int(surfaceCapabilities.minImageExtent.width), int(surfaceCapabilities.maxImageExtent.width))),
			height: C.uint(clamp(int(height), int(surfaceCapabilities.minImageExtent.height), int(surfaceCapabilities.maxImageExtent.height))),
		}
		return actualExtent
	}
	return surfaceCapabilities.currentExtent
}

func chooseSwapSurfaceFormat(surfaceFormats []C.VkSurfaceFormatKHR) C.VkSurfaceFormatKHR {
	for _, surfaceFormat := range surfaceFormats {
		if surfaceFormat.format == C.VK_FORMAT_B8G8R8A8_SRGB && surfaceFormat.colorSpace == C.VK_COLOR_SPACE_SRGB_NONLINEAR_KHR {
			return surfaceFormat
		}
	}
	return surfaceFormats[0]
}

func chooseSwapPresentMode(presentModes []C.VkPresentModeKHR) C.VkPresentModeKHR {
	for _, presentMode := range presentModes {
		if presentMode == C.VK_PRESENT_MODE_MAILBOX_KHR {
			return presentMode
		}
	}
	return C.VK_PRESENT_MODE_FIFO_KHR
}

func initSwapchain(app *App) (*C.VkSwapchainKHR, error) {
	extent := chooseSwapExtent(app, app.swapchainSupportInfo.surfaceCapabilities)
	surfaceFormat := chooseSwapSurfaceFormat(app.swapchainSupportInfo.surfaceFormats)
	presentMode := chooseSwapPresentMode(app.swapchainSupportInfo.presentModes)
	imageCount := app.swapchainSupportInfo.surfaceCapabilities.minImageCount
	// Max image count of zero means unlimited max image count.
	maxImageCount := app.swapchainSupportInfo.surfaceCapabilities.maxImageCount
	if maxImageCount == 0 || imageCount+1 <= maxImageCount {
		imageCount++ // use 1 image more than minimum to avoid having to wait for swap chain.
	}
	// Create swap chain.
	createInfo := C.VkSwapchainCreateInfoKHR{
		sType:            C.VK_STRUCTURE_TYPE_SWAPCHAIN_CREATE_INFO_KHR,
		surface:          *app.surface,
		minImageCount:    imageCount,
		imageFormat:      surfaceFormat.format,
		imageColorSpace:  surfaceFormat.colorSpace,
		imageExtent:      extent,
		imageArrayLayers: 1,
		imageUsage:       C.VK_IMAGE_USAGE_COLOR_ATTACHMENT_BIT, // NOTE: use VK_IMAGE_USAGE_TRANSFER_DST_BIT if rendering into separate image for post-processing, before copying result to swap chain.
		preTransform:     app.swapchainSupportInfo.surfaceCapabilities.currentTransform,
		compositeAlpha:   C.VK_COMPOSITE_ALPHA_OPAQUE_BIT_KHR,
		presentMode:      presentMode,
		clipped:          C.VK_TRUE, // NOTE: set to false if we need to be able to read pixels of areas obscured by other windows.
		oldSwapchain:     nil,
	}
	queueFamilyIndices := unique(app.QueueFamilyIndices.Indices()...)
	switch len(queueFamilyIndices) {
	case 1:
		// exclusive mode.
		createInfo.imageSharingMode = C.VK_SHARING_MODE_EXCLUSIVE
		createInfo.queueFamilyIndexCount = 0 // optional
		createInfo.pQueueFamilyIndices = nil // optional
	default:
		// concurrent mode.
		createInfo.imageSharingMode = C.VK_SHARING_MODE_CONCURRENT
		createInfo.queueFamilyIndexCount = C.uint(len(queueFamilyIndices))
		createInfo.pQueueFamilyIndices = getCUintSlice(queueFamilyIndices)
	}

	swapchain := C.new_VkSwapchainKHR()
	if result := C.vkCreateSwapchainKHR(*app.device, &createInfo, nil, swapchain); result != C.VK_SUCCESS {
		return nil, errors.Errorf("unable to create swap chain (result=%d)", result)
	}

	// Store swap chain image format and extent.
	app.swapchainImageFormat = surfaceFormat.format
	app.swapchainExtent = extent

	return swapchain, nil
}

func getSwapchainImgs(app *App) []C.VkImage {
	var nswapchainImgs C.uint32_t
	C.vkGetSwapchainImagesKHR(*app.device, *app.swapchain, &nswapchainImgs, nil)
	swapchainImgs := make([]C.VkImage, int(nswapchainImgs))
	C.vkGetSwapchainImagesKHR(*app.device, *app.swapchain, &nswapchainImgs, &swapchainImgs[0])
	return swapchainImgs
}

func initSwapchainImgViews(app *App) ([]C.VkImageView, error) {
	swapchainImgViews := make([]C.VkImageView, len(app.swapchainImgs))
	for i := range swapchainImgViews {
		createInfo := C.VkImageViewCreateInfo{
			sType:    C.VK_STRUCTURE_TYPE_IMAGE_VIEW_CREATE_INFO,
			image:    app.swapchainImgs[i],
			viewType: C.VK_IMAGE_VIEW_TYPE_2D,
			format:   app.swapchainImageFormat,
			components: C.VkComponentMapping{
				r: C.VK_COMPONENT_SWIZZLE_IDENTITY,
				g: C.VK_COMPONENT_SWIZZLE_IDENTITY,
				b: C.VK_COMPONENT_SWIZZLE_IDENTITY,
				a: C.VK_COMPONENT_SWIZZLE_IDENTITY,
			},
			subresourceRange: C.VkImageSubresourceRange{
				aspectMask:     C.VK_IMAGE_ASPECT_COLOR_BIT,
				baseMipLevel:   0,
				levelCount:     1,
				baseArrayLayer: 0,
				layerCount:     1,
			},
		}
		if result := C.vkCreateImageView(*app.device, &createInfo, nil, &swapchainImgViews[i]); result != C.VK_SUCCESS {
			return nil, errors.Errorf("unable to create image view of swap chain image (result=%d)", result)
		}
	}
	return swapchainImgViews, nil
}

func initRenderPass(app *App) (*C.VkRenderPass, error) {
	colorAttachment := C.VkAttachmentDescription{
		format:         app.swapchainImageFormat,
		samples:        C.VK_SAMPLE_COUNT_1_BIT,
		loadOp:         C.VK_ATTACHMENT_LOAD_OP_CLEAR,
		storeOp:        C.VK_ATTACHMENT_STORE_OP_STORE,
		stencilLoadOp:  C.VK_ATTACHMENT_LOAD_OP_DONT_CARE,  // NOTE: change if using stencils
		stencilStoreOp: C.VK_ATTACHMENT_STORE_OP_DONT_CARE, // NOTE: change if using stencils
		initialLayout:  C.VK_IMAGE_LAYOUT_UNDEFINED,
		finalLayout:    C.VK_IMAGE_LAYOUT_PRESENT_SRC_KHR,
	}
	colorAttachments := newVkAttachmentDescriptionSlice(colorAttachment)

	colorAttachmentRef := C.VkAttachmentReference{
		attachment: 0, // index of color attachment descriptor (we only have one).
		layout:     C.VK_IMAGE_LAYOUT_COLOR_ATTACHMENT_OPTIMAL,
	}
	colorAttachmentRefs := newVkAttachmentReferenceSlice(colorAttachmentRef)

	subpass := C.VkSubpassDescription{
		pipelineBindPoint:       C.VK_PIPELINE_BIND_POINT_GRAPHICS,
		inputAttachmentCount:    0,   // optional
		pInputAttachments:       nil, // optional
		colorAttachmentCount:    C.uint(len(colorAttachmentRefs)),
		pColorAttachments:       &colorAttachmentRefs[0],
		pResolveAttachments:     nil, // optional
		pDepthStencilAttachment: nil, // optional
		preserveAttachmentCount: 0,   // optional
		pPreserveAttachments:    nil, // optional
	}
	subpasses := newVkSubpassDescriptionSlice(subpass)

	renderPassCreateInfo := C.VkRenderPassCreateInfo{
		sType:           C.VK_STRUCTURE_TYPE_RENDER_PASS_CREATE_INFO,
		attachmentCount: C.uint(len(colorAttachments)),
		pAttachments:    &colorAttachments[0],
		subpassCount:    C.uint(len(subpasses)),
		pSubpasses:      &subpasses[0],
		dependencyCount: 0,   // optional
		pDependencies:   nil, // optional
	}
	renderPass := C.new_VkRenderPass()
	if result := C.vkCreateRenderPass(*app.device, &renderPassCreateInfo, nil, renderPass); result != C.VK_SUCCESS {
		return nil, errors.Errorf("unable to create render pass (result=%d)", result)
	}
	return renderPass, nil
}

func initGraphicsPipeline(app *App) ([]C.VkPipeline, error) {
	shaderStages, err := initShaderModules(app)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Vertex input.
	vertexInputState := C.VkPipelineVertexInputStateCreateInfo{
		sType:                           C.VK_STRUCTURE_TYPE_PIPELINE_VERTEX_INPUT_STATE_CREATE_INFO,
		vertexBindingDescriptionCount:   0,
		pVertexBindingDescriptions:      nil,
		vertexAttributeDescriptionCount: 0,
		pVertexAttributeDescriptions:    nil,
	}

	// Input assembler    (fixed-function stage)
	inputAssemblyState := C.VkPipelineInputAssemblyStateCreateInfo{
		sType:                  C.VK_STRUCTURE_TYPE_PIPELINE_INPUT_ASSEMBLY_STATE_CREATE_INFO,
		topology:               C.VK_PRIMITIVE_TOPOLOGY_TRIANGLE_LIST,
		primitiveRestartEnable: C.VK_FALSE,
	}

	// Vertex shader      (programmable)         // DONE
	//shaderStages[0]

	// Tessalation        (programmable)
	// Geometry shader    (programmable)

	// Viewports and scissors.
	viewport := C.VkViewport{
		x:        0.0,
		y:        0.0,
		width:    C.float(app.swapchainExtent.width),
		height:   C.float(app.swapchainExtent.height),
		minDepth: 0.0,
		maxDepth: 1.0,
	}
	viewports := newVkViewportSlice(viewport)
	scissor := C.VkRect2D{
		offset: C.VkOffset2D{x: 0, y: 0},
		extent: app.swapchainExtent,
	}
	scissors := newVkRect2DSlice(scissor)
	viewportState := C.VkPipelineViewportStateCreateInfo{
		sType:         C.VK_STRUCTURE_TYPE_PIPELINE_VIEWPORT_STATE_CREATE_INFO,
		viewportCount: C.uint(len(viewports)),
		pViewports:    &viewports[0],
		scissorCount:  C.uint(len(scissors)),
		pScissors:     &scissors[0],
	}

	// Rasterization      (fixed-function stage)
	rasterizationState := C.VkPipelineRasterizationStateCreateInfo{
		sType:                   C.VK_STRUCTURE_TYPE_PIPELINE_RASTERIZATION_STATE_CREATE_INFO,
		depthClampEnable:        C.VK_FALSE,
		rasterizerDiscardEnable: C.VK_FALSE,
		polygonMode:             C.VK_POLYGON_MODE_FILL,
		cullMode:                C.VK_CULL_MODE_BACK_BIT,
		frontFace:               C.VK_FRONT_FACE_CLOCKWISE,
		depthBiasEnable:         C.VK_FALSE,
		depthBiasConstantFactor: 0.0, // optional
		depthBiasClamp:          0.0, // optional
		depthBiasSlopeFactor:    0.0, // optional
		lineWidth:               1.0,
	}

	// Multisampling.
	multisampleState := C.VkPipelineMultisampleStateCreateInfo{
		sType:                 C.VK_STRUCTURE_TYPE_PIPELINE_MULTISAMPLE_STATE_CREATE_INFO,
		rasterizationSamples:  C.VK_SAMPLE_COUNT_1_BIT,
		sampleShadingEnable:   C.VK_FALSE,
		minSampleShading:      1.0,        // optional
		pSampleMask:           nil,        // optional
		alphaToCoverageEnable: C.VK_FALSE, // optional
		alphaToOneEnable:      C.VK_FALSE, // optional
	}

	// Depth and stencil testing.
	//depthStencilCreateInfo := C.VkPipelineDepthStencilStateCreateInfo

	// Fragment shader    (programmable)         // DONE
	//shaderStages[1]

	// Color blending     (fixed-function stage)
	colorBlendAttachment := C.VkPipelineColorBlendAttachmentState{
		blendEnable:         C.VK_FALSE,
		srcColorBlendFactor: C.VK_BLEND_FACTOR_ONE,  // optional
		dstColorBlendFactor: C.VK_BLEND_FACTOR_ZERO, // optional
		colorBlendOp:        C.VK_BLEND_OP_ADD,      // optional
		srcAlphaBlendFactor: C.VK_BLEND_FACTOR_ONE,  // optional
		dstAlphaBlendFactor: C.VK_BLEND_FACTOR_ZERO, // optional
		alphaBlendOp:        C.VK_BLEND_OP_ADD,      // optional
		colorWriteMask:      C.VK_COLOR_COMPONENT_R_BIT | C.VK_COLOR_COMPONENT_G_BIT | C.VK_COLOR_COMPONENT_B_BIT | C.VK_COLOR_COMPONENT_A_BIT,
	}
	colorBlendAttachments := newVkPipelineColorBlendAttachmentStateSlice(colorBlendAttachment)
	colorBlendState := C.VkPipelineColorBlendStateCreateInfo{
		sType:           C.VK_STRUCTURE_TYPE_PIPELINE_COLOR_BLEND_STATE_CREATE_INFO,
		logicOpEnable:   C.VK_FALSE,
		logicOp:         C.VK_LOGIC_OP_COPY, // optional
		attachmentCount: C.uint(len(colorBlendAttachments)),
		pAttachments:    &colorBlendAttachments[0],
		blendConstants:  [4]C.float{0.0, 0.0, 0.0, 0.0}, // optional
	}

	// Dynamic state.
	dynamicStates := []C.VkDynamicState{
		C.VK_DYNAMIC_STATE_VIEWPORT,
	}
	dynamicState := C.VkPipelineDynamicStateCreateInfo{
		sType:             C.VK_STRUCTURE_TYPE_PIPELINE_DYNAMIC_STATE_CREATE_INFO,
		dynamicStateCount: C.uint(len(dynamicStates)),
		pDynamicStates:    &dynamicStates[0],
	}

	// Uniform values.
	pipelineLayoutCreateInfo := C.VkPipelineLayoutCreateInfo{
		sType:                  C.VK_STRUCTURE_TYPE_PIPELINE_LAYOUT_CREATE_INFO,
		setLayoutCount:         0,   // optional
		pSetLayouts:            nil, // optional
		pushConstantRangeCount: 0,   // optional
		pPushConstantRanges:    nil, // optional
	}
	pipelineLayout := C.new_VkPipelineLayout()
	if result := C.vkCreatePipelineLayout(*app.device, &pipelineLayoutCreateInfo, nil, pipelineLayout); result != C.VK_SUCCESS {
		return nil, errors.Errorf("unable to create pipeline layout (result=%d)", result)
	}
	app.pipelineLayout = pipelineLayout

	graphicsPipelineCreateInfo := C.VkGraphicsPipelineCreateInfo{
		sType:               C.VK_STRUCTURE_TYPE_GRAPHICS_PIPELINE_CREATE_INFO,
		stageCount:          C.uint(len(shaderStages)),
		pStages:             &shaderStages[0],
		pVertexInputState:   &vertexInputState,
		pInputAssemblyState: &inputAssemblyState,
		pTessellationState:  nil, // optional
		pViewportState:      &viewportState,
		pRasterizationState: &rasterizationState,
		pMultisampleState:   &multisampleState,
		pDepthStencilState:  nil, // optional
		pColorBlendState:    &colorBlendState,
		pDynamicState:       &dynamicState,
		layout:              *pipelineLayout,
		renderPass:          *app.renderPass,
		subpass:             0,   // index of subpass in the render pass
		basePipelineHandle:  nil, // optional
		basePipelineIndex:   -1,  // optional
	}
	graphicsPipelineCreateInfos := newVkGraphicsPipelineCreateInfoSlice(graphicsPipelineCreateInfo)
	graphicsPipelines := newVkPipelineSlice(make([]C.VkPipeline, len(graphicsPipelineCreateInfos))...)
	if result := C.vkCreateGraphicsPipelines(*app.device, nil, C.uint(len(graphicsPipelineCreateInfos)), &graphicsPipelineCreateInfos[0], nil, &graphicsPipelines[0]); result != C.VK_SUCCESS {
		return nil, errors.Errorf("unable to create graphics pipeline (result=%d)", result)
	}
	return graphicsPipelines, nil
}

func initShaderModules(app *App) ([]C.VkPipelineShaderStageCreateInfo, error) {
	// Create vertex shader.
	vertexShaderModule, err := createShaderModule(app, "shaders/shader_vert.spv")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	app.vertexShaderModule = vertexShaderModule
	// Create fragment shader.
	fragmentShaderModule, err := createShaderModule(app, "shaders/shader_frag.spv")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	app.fragmentShaderModule = fragmentShaderModule
	// Create graphics pipeline.
	vertexShaderStageInfo := C.VkPipelineShaderStageCreateInfo{
		sType:  C.VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO,
		stage:  C.VK_SHADER_STAGE_VERTEX_BIT,
		module: *app.vertexShaderModule,
		pName:  C.CString("main"),
	}
	fragmentShaderStageInfo := C.VkPipelineShaderStageCreateInfo{
		sType:  C.VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO,
		stage:  C.VK_SHADER_STAGE_FRAGMENT_BIT,
		module: *app.fragmentShaderModule,
		pName:  C.CString("main"),
	}
	shaderStageCreateInfos := []C.VkPipelineShaderStageCreateInfo{
		vertexShaderStageInfo,
		fragmentShaderStageInfo,
	}
	return shaderStageCreateInfos, nil
}

func createShaderModule(app *App, shaderPath string) (*C.VkShaderModule, error) {
	dbg.Printf("loading shader %q", shaderPath)
	shaderData, err := ioutil.ReadFile(shaderPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	createInfo := C.VkShaderModuleCreateInfo{
		sType:    C.VK_STRUCTURE_TYPE_SHADER_MODULE_CREATE_INFO,
		codeSize: C.size_t(len(shaderData)),
		pCode:    getCUintSliceFromBytes(shaderData),
	}
	shaderModule := C.new_VkShaderModule()
	if result := C.vkCreateShaderModule(*app.device, &createInfo, nil, shaderModule); result != C.VK_SUCCESS {
		return nil, errors.Errorf("unable to create shader module %q (result=%d)", shaderPath, result)
	}
	return shaderModule, nil
}

// ### [ Helper functions ] ####################################################

func contains(ss []string, s string) bool {
	for i := range ss {
		if ss[i] == s {
			return true
		}
	}
	return false
}

func unique(xs ...int) []int {
	m := make(map[int]bool)
	for _, x := range xs {
		m[x] = true
	}
	var out []int
	for x := range m {
		out = append(out, x)
	}
	sort.Ints(out)
	return out
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	// min <= v && v <= max
	return v
}
