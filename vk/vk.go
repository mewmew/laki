// TODO: continue at https://vulkan-tutorial.com/en/Drawing_a_triangle/Presentation/Swap_chain

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
	device, err := initDevice(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.device = device
	initQueues(app)
	pretty.Println("app.QueueFamilyIndices:", app.QueueFamilyIndices)
	return nil
}

func CleanupVulkan(app *App) {
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
		if !isSuitableDevice(app, &physicalDevice) {
			continue
		}
		_physicalDevice := C.new_VkPhysicalDevice()
		*_physicalDevice = physicalDevice // allocate pointer on C heap.
		return _physicalDevice, nil
	}
	return nil, errors.Errorf("unable to locate suitable physical device (GPU)")
}

func isSuitableDevice(app *App, physicalDevice *C.VkPhysicalDevice) bool {
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

	swapChainSupportInfo := getSwapChainSupportInfo(app, physicalDevice)
	if len(swapChainSupportInfo.surfaceFormats) == 0 {
		return false
	}
	if len(swapChainSupportInfo.presentModes) == 0 {
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
	queueFamilyIndices := unique(graphicsQueueFamilyIndex, presentQueueFamilyIndex)
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

func getSwapChainSupportInfo(app *App, physicalDevice *C.VkPhysicalDevice) *SwapChainSupportInfo {
	// Get surface capabilities.
	swapChainSupportInfo := &SwapChainSupportInfo{}
	var surfaceCapabilities C.VkSurfaceCapabilitiesKHR
	C.vkGetPhysicalDeviceSurfaceCapabilitiesKHR(*physicalDevice, *app.surface, &surfaceCapabilities)
	swapChainSupportInfo.surfaceCapabilities = &surfaceCapabilities

	// Get surface formats.
	var nsurfaceFormats C.uint32_t
	C.vkGetPhysicalDeviceSurfaceFormatsKHR(*physicalDevice, *app.surface, &nsurfaceFormats, nil)
	surfaceFormats := make([]C.VkSurfaceFormatKHR, int(nsurfaceFormats))
	C.vkGetPhysicalDeviceSurfaceFormatsKHR(*physicalDevice, *app.surface, &nsurfaceFormats, &surfaceFormats[0])
	swapChainSupportInfo.surfaceFormats = surfaceFormats

	// Get present modes.
	var npresentModes C.uint32_t
	C.vkGetPhysicalDeviceSurfacePresentModesKHR(*physicalDevice, *app.surface, &npresentModes, nil)
	presentModes := make([]C.VkPresentModeKHR, int(npresentModes))
	C.vkGetPhysicalDeviceSurfacePresentModesKHR(*physicalDevice, *app.surface, &npresentModes, &presentModes[0])
	swapChainSupportInfo.presentModes = presentModes

	return swapChainSupportInfo
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
