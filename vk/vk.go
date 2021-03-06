// TODO: continue at https://vulkan-tutorial.com/en/Uniform_buffers/Descriptor_layout_and_buffer

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
//#cgo linux pkg-config: glfw3
//#cgo linux pkg-config: vulkan
//
//#cgo windows LDFLAGS: -lglfw3dll
//#cgo windows LDFLAGS: -lvulkan-1
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

// Maximum number of frames processed concurrently by GPU.
const MaxFramesInFlight = 2

func Init() error {
	app := newApp()
	app.win = InitWindow(app)
	defer CleanupWindow(app.win)
	if err := InitVulkan(app); err != nil {
		return errors.WithStack(err)
	}
	defer CleanupVulkan(app)

	EventLoop(app)
	return nil
}

func InitVulkan(app *App) error {
	// Create Vulkan instance.
	instance, err := initInstance()
	if err != nil {
		return errors.WithStack(err)
	}
	app.instance = instance
	// Create debug messanger.
	debugMessanger, err := initDebugMessanger(app.instance)
	if err != nil {
		return errors.WithStack(err)
	}
	app.debugMessanger = debugMessanger
	// Create Vulkan surface.
	surface, err := initSurface(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.surface = surface
	// Create Vulkan physical device.
	physicalDevice, err := initPhysicalDevice(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.physicalDevice = physicalDevice
	// Create Vulkan logical device.
	device, err := initDevice(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.device = device
	// Init queue indices.
	initQueues(app)

	// Create swapchain.
	swapchain, err := initSwapchain(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.swapchain = swapchain
	// Create swapchain images.
	app.swapchainImgs = getSwapchainImgs(app)
	// Create swapchain image views.
	swapchainImgViews, err := initSwapchainImgViews(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.swapchainImgViews = swapchainImgViews
	// Create render pass.
	renderPass, err := initRenderPass(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.renderPass = renderPass
	// Create graphics pipeline.
	graphicsPipelines, err := initGraphicsPipeline(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.graphicsPipelines = graphicsPipelines
	// Create framebuffers.
	framebuffers, err := initFramebuffers(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.swapchainFramebuffers = framebuffers
	// Create command pool.
	//
	// NOTE: command pool does not need to be re-initialized during
	// recreateSwapchain.
	commandPool, err := initCommandPool(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.commandPool = commandPool
	// Create vertex buffer.
	// top-left
	topLeft := Vertex{
		pos:   vec2(-0.5, -0.5),    // x, y
		color: vec3(1.0, 0.0, 0.0), // red
	}
	// top-right
	topRight := Vertex{
		pos:   vec2(0.5, -0.5),     // x, y
		color: vec3(0.0, 1.0, 0.0), // green
	}
	// bottom-right
	bottomRight := Vertex{
		pos:   vec2(0.5, 0.5),      // x, y
		color: vec3(0.0, 0.0, 1.0), // blue
	}
	// bottom-left
	bottomLeft := Vertex{
		pos:   vec2(-0.5, 0.5),     // x, y
		color: vec3(1.0, 1.0, 1.0), // white
	}
	vertices := []Vertex{
		// first triangle.
		topLeft,
		topRight,
		bottomRight,
		// second triangle.
		bottomRight,
		bottomLeft,
		topLeft,
	}
	indices, uniqueVertices := uniqueIndexList(vertices)
	app.indices = indices
	app.uniqueVertices = uniqueVertices
	if err := createVertexBuffer(app, uniqueVertices); err != nil {
		return errors.WithStack(err)
	}
	// Create index buffer in GPU memory.
	if err := createIndexBuffer(app, indices); err != nil {
		return errors.WithStack(err)
	}
	// Create command buffers.
	commandBuffers, err := initCommandBuffers(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.swapchainCommandBuffers = commandBuffers
	if err := recordRenderCommands(app); err != nil {
		return errors.WithStack(err)
	}
	// Sync objects.
	if err := initSyncObjects(app); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func CleanupVulkan(app *App) {
	for i := range app.imageAvailableSemaphores {
		C.vkDestroyFence(*app.device, *app.imagesInFlightFences[i], nil)
		C.vkDestroyFence(*app.device, *app.framesInFlightFences[i], nil)
		C.vkDestroySemaphore(*app.device, *app.imageAvailableSemaphores[i], nil)
		C.vkDestroySemaphore(*app.device, *app.renderFinishedSemaphores[i], nil)
	}
	C.vkFreeMemory(*app.device, *app.indexBufferMem, nil)
	C.vkDestroyBuffer(*app.device, *app.indexBuffer, nil)
	C.vkFreeMemory(*app.device, *app.vertexBufferMem, nil)
	C.vkDestroyBuffer(*app.device, *app.vertexBuffer, nil)
	cleanupSwapchain(app)
	C.vkDestroyCommandPool(*app.device, *app.commandPool, nil)
	C.vkDestroyDevice(*app.device, nil) // free command pool after command buffers allocated in pool.
	app.physicalDevice = nil
	DestroyDebugUtilsMessengerEXT(*app.instance, *app.debugMessanger, nil)
	C.vkDestroySurfaceKHR(*app.instance, *app.surface, nil)
	C.vkDestroyInstance(*app.instance, nil)
}

func cleanupSwapchain(app *App) {
	for i := range app.swapchainFramebuffers {
		if app.swapchainFramebuffers[i] != nil {
			C.vkDestroyFramebuffer(*app.device, app.swapchainFramebuffers[i], nil)
			app.swapchainFramebuffers[i] = nil
		}
	}
	if len(app.swapchainCommandBuffers) > 0 {
		C.vkFreeCommandBuffers(*app.device, *app.commandPool, C.uint(len(app.swapchainCommandBuffers)), &app.swapchainCommandBuffers[0])
		app.swapchainCommandBuffers = nil
	}
	if len(app.graphicsPipelines) > 0 {
		for _, graphicsPipeline := range app.graphicsPipelines {
			C.vkDestroyPipeline(*app.device, graphicsPipeline, nil)
		}
		app.graphicsPipelines = nil
	}
	if app.pipelineLayout != nil {
		C.vkDestroyPipelineLayout(*app.device, *app.pipelineLayout, nil)
		app.pipelineLayout = nil
	}
	if app.renderPass != nil {
		C.vkDestroyRenderPass(*app.device, *app.renderPass, nil)
		app.renderPass = nil
	}
	if len(app.swapchainImgViews) > 0 {
		for i := range app.swapchainImgViews {
			C.vkDestroyImageView(*app.device, app.swapchainImgViews[i], nil)
		}
		app.swapchainImgViews = nil
	}
	if app.swapchain != nil {
		C.vkDestroySwapchainKHR(*app.device, *app.swapchain, nil)
		app.swapchain = nil
	}
}

func initInstance() (*C.VkInstance, error) {
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
		warn.Printf("   missing required device extensions: %#v", m)
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
	dbg.Println("vk.chooseSwapExtent")
	if surfaceCapabilities.currentExtent.width == C.UINT32_MAX || surfaceCapabilities.currentExtent.height == C.UINT32_MAX {
		var width, height C.int
		C.glfwGetFramebufferSize(app.win, &width, &height)
		dbg.Printf("   framebuffer size (%dx%d)", width, height)
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

func recreateSwapchain(app *App) error {
	var width, height C.int
	for {
		C.glfwGetFramebufferSize(app.win, &width, &height)
		minimized := width == 0 || height == 0
		if !minimized {
			break
		}
		C.glfwWaitEvents() // wait until window is not minimized.
	}

	if result := C.vkDeviceWaitIdle(*app.device); result != C.VK_SUCCESS {
		return errors.Errorf("unable to wait for device to become idle (result=%d)", result)
	}

	cleanupSwapchain(app)

	// Create swapchain.
	swapchain, err := initSwapchain(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.swapchain = swapchain
	// Create swapchain images.
	app.swapchainImgs = getSwapchainImgs(app)
	// Create swapchain image views.
	swapchainImgViews, err := initSwapchainImgViews(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.swapchainImgViews = swapchainImgViews
	// Create render pass.
	renderPass, err := initRenderPass(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.renderPass = renderPass
	// Create graphics pipeline.
	graphicsPipelines, err := initGraphicsPipeline(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.graphicsPipelines = graphicsPipelines
	// Create framebuffers.
	framebuffers, err := initFramebuffers(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.swapchainFramebuffers = framebuffers
	// Create command buffers.
	commandBuffers, err := initCommandBuffers(app)
	if err != nil {
		return errors.WithStack(err)
	}
	app.swapchainCommandBuffers = commandBuffers
	if err := recordRenderCommands(app); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func initSwapchain(app *App) (*C.VkSwapchainKHR, error) {
	dbg.Println("vk.initSwapchain")
	swapchainSupportInfo := getSwapchainSupportInfo(app, app.physicalDevice)
	app.swapchainSupportInfo = swapchainSupportInfo
	pretty.Println("   swapchainSupportInfo:", swapchainSupportInfo)
	extent := chooseSwapExtent(app, app.swapchainSupportInfo.surfaceCapabilities)
	dbg.Println("   extent:", extent)
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
	dependency := C.VkSubpassDependency{
		srcSubpass:      C.VK_SUBPASS_EXTERNAL,
		dstSubpass:      0, // index of first and only subpass.
		srcStageMask:    C.VK_PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT,
		dstStageMask:    C.VK_PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT,
		srcAccessMask:   0,
		dstAccessMask:   C.VK_ACCESS_COLOR_ATTACHMENT_WRITE_BIT,
		dependencyFlags: 0, // optional
	}
	dependencies := newVkSubpassDependencySlice(dependency)
	renderPassCreateInfo := C.VkRenderPassCreateInfo{
		sType:           C.VK_STRUCTURE_TYPE_RENDER_PASS_CREATE_INFO,
		attachmentCount: C.uint(len(colorAttachments)),
		pAttachments:    &colorAttachments[0],
		subpassCount:    C.uint(len(subpasses)),
		pSubpasses:      &subpasses[0],
		dependencyCount: C.uint(len(dependencies)),
		pDependencies:   &dependencies[0],
	}
	renderPass := C.new_VkRenderPass()
	if result := C.vkCreateRenderPass(*app.device, &renderPassCreateInfo, nil, renderPass); result != C.VK_SUCCESS {
		return nil, errors.Errorf("unable to create render pass (result=%d)", result)
	}
	return renderPass, nil
}

func initGraphicsPipeline(app *App) ([]C.VkPipeline, error) {
	shaderStages, cleanupShaderModules, err := initShaderModules(app)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer cleanupShaderModules()

	// Vertex input.
	bindingDescs, attrDescs := getBindingDescs()
	vertexInputState := C.VkPipelineVertexInputStateCreateInfo{
		sType:                           C.VK_STRUCTURE_TYPE_PIPELINE_VERTEX_INPUT_STATE_CREATE_INFO,
		vertexBindingDescriptionCount:   C.uint(len(bindingDescs)),
		pVertexBindingDescriptions:      &bindingDescs[0],
		vertexAttributeDescriptionCount: C.uint(len(attrDescs)),
		pVertexAttributeDescriptions:    &attrDescs[0],
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
	//dynamicStates := []C.VkDynamicState{
	//	C.VK_DYNAMIC_STATE_VIEWPORT,
	//}
	//dynamicState := C.VkPipelineDynamicStateCreateInfo{
	//	sType:             C.VK_STRUCTURE_TYPE_PIPELINE_DYNAMIC_STATE_CREATE_INFO,
	//	dynamicStateCount: C.uint(len(dynamicStates)),
	//	pDynamicStates:    &dynamicStates[0],
	//}

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
		//pDynamicState:       &dynamicState,
		layout:             *pipelineLayout,
		renderPass:         *app.renderPass,
		subpass:            0,   // index of subpass in the render pass
		basePipelineHandle: nil, // optional
		basePipelineIndex:  -1,  // optional
	}
	graphicsPipelineCreateInfos := newVkGraphicsPipelineCreateInfoSlice(graphicsPipelineCreateInfo)
	graphicsPipelines := newVkPipelineSlice(make([]C.VkPipeline, len(graphicsPipelineCreateInfos))...)
	if result := C.vkCreateGraphicsPipelines(*app.device, nil, C.uint(len(graphicsPipelineCreateInfos)), &graphicsPipelineCreateInfos[0], nil, &graphicsPipelines[0]); result != C.VK_SUCCESS {
		return nil, errors.Errorf("unable to create graphics pipeline (result=%d)", result)
	}
	return graphicsPipelines, nil
}

func initShaderModules(app *App) (shaderStageCreateInfos []C.VkPipelineShaderStageCreateInfo, cleanup func(), err error) {
	// Create vertex shader.
	vertexShaderModule, err := createShaderModule(app, "shaders/shader_vert.spv")
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	// Create fragment shader.
	fragmentShaderModule, err := createShaderModule(app, "shaders/shader_frag.spv")
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	// Create graphics pipeline.
	vertexShaderStageInfo := C.VkPipelineShaderStageCreateInfo{
		sType:  C.VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO,
		stage:  C.VK_SHADER_STAGE_VERTEX_BIT,
		module: *vertexShaderModule,
		pName:  C.CString("main"),
	}
	fragmentShaderStageInfo := C.VkPipelineShaderStageCreateInfo{
		sType:  C.VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO,
		stage:  C.VK_SHADER_STAGE_FRAGMENT_BIT,
		module: *fragmentShaderModule,
		pName:  C.CString("main"),
	}
	shaderStageCreateInfos = []C.VkPipelineShaderStageCreateInfo{
		vertexShaderStageInfo,
		fragmentShaderStageInfo,
	}
	cleanup = func() {
		C.vkDestroyShaderModule(*app.device, *fragmentShaderModule, nil)
		C.vkDestroyShaderModule(*app.device, *vertexShaderModule, nil)
	}
	return shaderStageCreateInfos, cleanup, nil
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

func initFramebuffers(app *App) ([]C.VkFramebuffer, error) {
	framebuffers := newVkFramebufferSlice(make([]C.VkFramebuffer, len(app.swapchainImgViews))...)
	for i := range app.swapchainImgViews {
		attachments := newVkImageViewSlice(app.swapchainImgViews[i])
		framebufferCreateInfo := C.VkFramebufferCreateInfo{
			sType:           C.VK_STRUCTURE_TYPE_FRAMEBUFFER_CREATE_INFO,
			renderPass:      *app.renderPass,
			attachmentCount: C.uint(len(attachments)),
			pAttachments:    &attachments[0],
			width:           app.swapchainExtent.width,
			height:          app.swapchainExtent.height,
			layers:          1,
		}
		if result := C.vkCreateFramebuffer(*app.device, &framebufferCreateInfo, nil, &framebuffers[i]); result != C.VK_SUCCESS {
			return nil, errors.Errorf("unable to create framebuffer (result=%d)", result)
		}
	}
	return framebuffers, nil
}

func initCommandPool(app *App) (*C.VkCommandPool, error) {
	commandPoolCreateInfo := C.VkCommandPoolCreateInfo{
		sType:            C.VK_STRUCTURE_TYPE_COMMAND_POOL_CREATE_INFO,
		queueFamilyIndex: C.uint(app.graphicsQueueFamilyIndex),
	}
	commandPool := C.new_VkCommandPool()
	if result := C.vkCreateCommandPool(*app.device, &commandPoolCreateInfo, nil, commandPool); result != C.VK_SUCCESS {
		return nil, errors.Errorf("unable to create command pool (result=%d)", result)
	}
	return commandPool, nil
}

func initCommandBuffers(app *App) ([]C.VkCommandBuffer, error) {
	commandBuffers := newVkCommandBufferSlice(make([]C.VkCommandBuffer, len(app.swapchainFramebuffers))...)
	commandBufferAllocateInfo := C.VkCommandBufferAllocateInfo{
		sType:              C.VK_STRUCTURE_TYPE_COMMAND_BUFFER_ALLOCATE_INFO,
		commandPool:        *app.commandPool,
		level:              C.VK_COMMAND_BUFFER_LEVEL_PRIMARY,
		commandBufferCount: C.uint(len(commandBuffers)),
	}
	if result := C.vkAllocateCommandBuffers(*app.device, &commandBufferAllocateInfo, &commandBuffers[0]); result != C.VK_SUCCESS {
		return nil, errors.Errorf("unable to create command buffers (result=%d)", result)
	}
	return commandBuffers, nil
}

func recordRenderCommands(app *App) error {
	for i := range app.swapchainCommandBuffers {
		commandBufferBeginInfo := C.VkCommandBufferBeginInfo{
			sType:            C.VK_STRUCTURE_TYPE_COMMAND_BUFFER_BEGIN_INFO,
			flags:            0,   // optional
			pInheritanceInfo: nil, // optional
		}
		if result := C.vkBeginCommandBuffer(app.swapchainCommandBuffers[i], &commandBufferBeginInfo); result != C.VK_SUCCESS {
			return errors.Errorf("unable to begin recording command buffer (result=%d)", result)
		}

		clearColor := C.VkClearValue{
			0.0, 0.0, 0.0, 1.0, // r, g, b, a
		}
		clearColors := newVkClearValueSlice(clearColor)

		renderPassBeginInfo := C.VkRenderPassBeginInfo{
			sType:       C.VK_STRUCTURE_TYPE_RENDER_PASS_BEGIN_INFO,
			renderPass:  *app.renderPass,
			framebuffer: app.swapchainFramebuffers[i],
			renderArea: C.VkRect2D{
				offset: C.VkOffset2D{x: 0, y: 0},
				extent: app.swapchainExtent,
			},
			clearValueCount: C.uint(len(clearColors)),
			pClearValues:    &clearColors[0],
		}

		C.vkCmdBeginRenderPass(app.swapchainCommandBuffers[i], &renderPassBeginInfo, C.VK_SUBPASS_CONTENTS_INLINE)

		C.vkCmdBindPipeline(app.swapchainCommandBuffers[i], C.VK_PIPELINE_BIND_POINT_GRAPHICS, app.graphicsPipelines[0]) // NOTE: we only use one graphics pipeline.

		vertexBuffers := []C.VkBuffer{
			*app.vertexBuffer,
		}
		offsets := []C.VkDeviceSize{
			0,
		}
		const firstVertexBufferBinding = 0
		C.vkCmdBindVertexBuffers(app.swapchainCommandBuffers[i], firstVertexBufferBinding, C.uint(len(vertexBuffers)), &vertexBuffers[0], &offsets[0])
		const indexBufferOffset = 0
		C.vkCmdBindIndexBuffer(app.swapchainCommandBuffers[i], *app.indexBuffer, indexBufferOffset, C.VK_INDEX_TYPE_UINT32)

		const (
			instanceCount = 1
			firstIndex    = 0
			vertexOffset  = 0
			firstInstance = 0
		)
		C.vkCmdDrawIndexed(app.swapchainCommandBuffers[i], C.uint(len(app.indices)), instanceCount, firstIndex, vertexOffset, firstInstance)

		C.vkCmdEndRenderPass(app.swapchainCommandBuffers[i])

		if result := C.vkEndCommandBuffer(app.swapchainCommandBuffers[i]); result != C.VK_SUCCESS {
			return errors.Errorf("unable to record command buffer (result=%d)", result)
		}
	}
	return nil
}

func initSyncObjects(app *App) error {
	semaphoreCreateInfo := C.VkSemaphoreCreateInfo{
		sType: C.VK_STRUCTURE_TYPE_SEMAPHORE_CREATE_INFO,
	}
	fenceCreateInfo := C.VkFenceCreateInfo{
		sType: C.VK_STRUCTURE_TYPE_FENCE_CREATE_INFO,
		flags: C.VK_FENCE_CREATE_SIGNALED_BIT, // start fence in signalled state.
	}
	app.imagesInFlightFences = make([]*C.VkFence, len(app.swapchainImgs))
	for i := range app.imageAvailableSemaphores {
		// Image available semaphore.
		imageAvailableSemaphore := C.new_VkSemaphore()
		if result := C.vkCreateSemaphore(*app.device, &semaphoreCreateInfo, nil, imageAvailableSemaphore); result != C.VK_SUCCESS {
			return errors.Errorf("unable to create semaphore (result=%d)", result)
		}
		app.imageAvailableSemaphores[i] = imageAvailableSemaphore
		// Rendering finished semaphore.
		renderFinishedSemaphore := C.new_VkSemaphore()
		if result := C.vkCreateSemaphore(*app.device, &semaphoreCreateInfo, nil, renderFinishedSemaphore); result != C.VK_SUCCESS {
			return errors.Errorf("unable to create semaphore (result=%d)", result)
		}
		app.renderFinishedSemaphores[i] = renderFinishedSemaphore
		// In-flight fence.
		framesInFlightFence := C.new_VkFence()
		if result := C.vkCreateFence(*app.device, &fenceCreateInfo, nil, framesInFlightFence); result != C.VK_SUCCESS {
			return errors.Errorf("unable to create fence (result=%d)", result)
		}
		app.framesInFlightFences[i] = framesInFlightFence
		// Images in-flight fence.
		imagesInFlightFence := C.new_VkFence()
		if result := C.vkCreateFence(*app.device, &fenceCreateInfo, nil, imagesInFlightFence); result != C.VK_SUCCESS {
			return errors.Errorf("unable to create fence (result=%d)", result)
		}
		app.imagesInFlightFences[i] = imagesInFlightFence
	}
	return nil
}

func drawFrame(app *App) error {
	const (
		nfences = 1
		timeout = C.UINT64_MAX // disable timeout
	)
	C.vkWaitForFences(*app.device, nfences, app.framesInFlightFences[app.curFrame], C.VK_TRUE, timeout)

	//dbg.Println("vk.drawFrame")
	var imageIndex C.uint32_t // swapchainImgs array index
	if result := C.vkAcquireNextImageKHR(*app.device, *app.swapchain, timeout, *app.imageAvailableSemaphores[app.curFrame], nil, &imageIndex); result != C.VK_SUCCESS {
		switch result {
		case C.VK_ERROR_OUT_OF_DATE_KHR:
			// Recreate swapchain; window resolution has most likely been changed.
			if err := recreateSwapchain(app); err != nil {
				return errors.WithStack(err)
			}
			return nil // early return, try again on next call to drawFrame.
		case C.VK_SUBOPTIMAL_KHR:
			// nothing to do; present aquired image even if suboptimal.
		default:
			return errors.Errorf("unable to aquire next image (result=%d)", result)
		}
	}
	// check if frame is used by previous frame.
	if app.imagesInFlightFences[imageIndex] != nil {
		C.vkWaitForFences(*app.device, nfences, app.imagesInFlightFences[imageIndex], C.VK_TRUE, timeout)
	}

	waitSemaphores := newVkSemaphoreSlice(*app.imageAvailableSemaphores[app.curFrame])
	waitStages := []C.VkPipelineStageFlags{
		C.VK_PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT,
	}
	signalSemaphores := newVkSemaphoreSlice(*app.renderFinishedSemaphores[app.curFrame])
	submitInfo := C.VkSubmitInfo{
		sType:                C.VK_STRUCTURE_TYPE_SUBMIT_INFO,
		waitSemaphoreCount:   C.uint(len(waitSemaphores)),
		pWaitSemaphores:      &waitSemaphores[0],
		pWaitDstStageMask:    &waitStages[0],
		commandBufferCount:   1,
		pCommandBuffers:      &app.swapchainCommandBuffers[imageIndex],
		signalSemaphoreCount: C.uint(len(signalSemaphores)),
		pSignalSemaphores:    &signalSemaphores[0],
	}
	submits := newVkSubmitInfoSlice(submitInfo)
	C.vkResetFences(*app.device, nfences, app.framesInFlightFences[app.curFrame])
	if result := C.vkQueueSubmit(*app.graphicsQueue, C.uint(len(submits)), &submits[0], *app.framesInFlightFences[app.curFrame]); result != C.VK_SUCCESS {
		return errors.Errorf("unable to submit command buffers to graphics queue (result=%d)", result)
	}
	// Present frame.
	swapchains := newVkSwapchainKHRSlice(*app.swapchain)
	imageIndices := newCUint32Slice(imageIndex)
	presentInfo := C.VkPresentInfoKHR{
		sType:              C.VK_STRUCTURE_TYPE_PRESENT_INFO_KHR,
		waitSemaphoreCount: C.uint(len(signalSemaphores)),
		pWaitSemaphores:    &signalSemaphores[0],
		swapchainCount:     C.uint(len(swapchains)),
		pSwapchains:        &swapchains[0],
		pImageIndices:      &imageIndices[0],
		pResults:           nil, // optional
	}
	result := C.vkQueuePresentKHR(*app.presentQueue, &presentInfo)
	switch {
	case result == C.VK_ERROR_OUT_OF_DATE_KHR, result == C.VK_SUBOPTIMAL_KHR, app.framebufferResized:
		// Recreate swapchain; window resolution has most likely been changed.
		if err := recreateSwapchain(app); err != nil {
			return errors.WithStack(err)
		}
		app.framebufferResized = false
	default:
		if result != C.VK_SUCCESS {
			return errors.Errorf("unable to queue image for presentation (result=%d)", result)
		}
	}

	app.curFrame = (app.curFrame + 1) % MaxFramesInFlight

	return nil
}

func getVerticesSize(vertices []Vertex) C.VkDeviceSize {
	return C.VkDeviceSize(int(unsafe.Sizeof(vertices[0])) * len(vertices))
}

func getIndicesSize(indices []uint32) C.VkDeviceSize {
	return C.VkDeviceSize(int(unsafe.Sizeof(indices[0])) * len(indices))
}

func fillVertexBuffer(app *App, vertices []Vertex, bufferMem *C.VkDeviceMemory) error {
	// Fill vertex buffer with data.
	const offset = 0
	var data unsafe.Pointer
	size := getVerticesSize(vertices)
	if result := C.vkMapMemory(*app.device, *bufferMem, offset, size, 0, &data); result != C.VK_SUCCESS {
		return errors.Errorf("unable to map memory of vertex buffer with size=%d (result=%d)", size, result)
	}
	dst := unsafe.Slice((*byte)(data), size)
	src := unsafe.Slice((*byte)(unsafe.Pointer(&vertices[0])), size)
	copy(dst, src)
	C.vkUnmapMemory(*app.device, *bufferMem)
	return nil
}

func fillIndexBuffer(app *App, indices []uint32, bufferMem *C.VkDeviceMemory) error {
	// Fill index buffer with data.
	const offset = 0
	var data unsafe.Pointer
	size := getIndicesSize(indices)
	if result := C.vkMapMemory(*app.device, *bufferMem, offset, size, 0, &data); result != C.VK_SUCCESS {
		return errors.Errorf("unable to map memory of index buffer with size=%d (result=%d)", size, result)
	}
	dst := unsafe.Slice((*byte)(data), size)
	src := unsafe.Slice((*byte)(unsafe.Pointer(&indices[0])), size)
	copy(dst, src)
	C.vkUnmapMemory(*app.device, *bufferMem)
	return nil
}

func findMemoryType(app *App, typeFilter C.uint, properties C.VkMemoryPropertyFlags) (uint32, error) {
	var memProperties C.VkPhysicalDeviceMemoryProperties
	C.vkGetPhysicalDeviceMemoryProperties(*app.physicalDevice, &memProperties)
	for i := 0; i < int(memProperties.memoryTypeCount); i++ {
		if typeFilter&C.uint(1<<i) != 0 && memProperties.memoryTypes[i].propertyFlags&properties == properties {
			return uint32(i), nil
		}
	}
	return 0, errors.Errorf("unable to find suitable memory type for filter 0x%08X", uint32(typeFilter))
}

func createBuffer(app *App, size C.VkDeviceSize, usage C.VkBufferUsageFlags, properties C.VkMemoryPropertyFlags) (*C.VkBuffer, *C.VkDeviceMemory, error) {
	bufferCreateInfo := C.VkBufferCreateInfo{
		sType:                 C.VK_STRUCTURE_TYPE_BUFFER_CREATE_INFO,
		size:                  size,
		usage:                 usage,
		sharingMode:           C.VK_SHARING_MODE_EXCLUSIVE,
		queueFamilyIndexCount: 0,   // optional
		pQueueFamilyIndices:   nil, // optional
	}
	buffer := C.new_VkBuffer()
	if result := C.vkCreateBuffer(*app.device, &bufferCreateInfo, nil, buffer); result != C.VK_SUCCESS {
		return nil, nil, errors.Errorf("unable to create vertex buffer (result=%d)", result)
	}
	// Get memory requirements.
	var memRequirements C.VkMemoryRequirements
	C.vkGetBufferMemoryRequirements(*app.device, *buffer, &memRequirements)
	// Allocate memory.
	memoryTypeIndex, err := findMemoryType(app, memRequirements.memoryTypeBits, properties)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	memAllocInfo := C.VkMemoryAllocateInfo{
		sType:           C.VK_STRUCTURE_TYPE_MEMORY_ALLOCATE_INFO,
		allocationSize:  memRequirements.size,
		memoryTypeIndex: C.uint(memoryTypeIndex),
	}
	bufferMem := C.new_VkDeviceMemory()
	if result := C.vkAllocateMemory(*app.device, &memAllocInfo, nil, bufferMem); result != C.VK_SUCCESS {
		return nil, nil, errors.Errorf("unable to allocate memory of size=%d (result=%d)", memRequirements.size, result)
	}
	const memoryOffset = 0
	if result := C.vkBindBufferMemory(*app.device, *buffer, *bufferMem, memoryOffset); result != C.VK_SUCCESS {
		return nil, nil, errors.Errorf("unable to bind memory of vertex buffer (result=%d)", result)
	}
	return buffer, bufferMem, nil
}

func copyBuffer(app *App, dstBuffer, srcBuffer *C.VkBuffer, size C.VkDeviceSize) error {
	tmpCommandBuffers := newVkCommandBufferSlice(make([]C.VkCommandBuffer, 1)...)
	commandBufferAllocateInfo := C.VkCommandBufferAllocateInfo{
		sType:              C.VK_STRUCTURE_TYPE_COMMAND_BUFFER_ALLOCATE_INFO,
		commandPool:        *app.commandPool,
		level:              C.VK_COMMAND_BUFFER_LEVEL_PRIMARY,
		commandBufferCount: C.uint(len(tmpCommandBuffers)),
	}
	if result := C.vkAllocateCommandBuffers(*app.device, &commandBufferAllocateInfo, &tmpCommandBuffers[0]); result != C.VK_SUCCESS {
		return errors.Errorf("unable to create command buffers (result=%d)", result)
	}
	defer C.vkFreeCommandBuffers(*app.device, *app.commandPool, C.uint(len(tmpCommandBuffers)), &tmpCommandBuffers[0])
	commandBufferBeginInfo := C.VkCommandBufferBeginInfo{
		sType:            C.VK_STRUCTURE_TYPE_COMMAND_BUFFER_BEGIN_INFO,
		flags:            C.VK_COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT,
		pInheritanceInfo: nil, // optional
	}
	if result := C.vkBeginCommandBuffer(tmpCommandBuffers[0], &commandBufferBeginInfo); result != C.VK_SUCCESS {
		return errors.Errorf("unable to begin recording command buffer (result=%d)", result)
	}
	copyRegions := []C.VkBufferCopy{
		{
			srcOffset: 0,
			dstOffset: 0,
			size:      size,
		},
	}
	C.vkCmdCopyBuffer(tmpCommandBuffers[0], *srcBuffer, *dstBuffer, C.uint(len(copyRegions)), &copyRegions[0])
	if result := C.vkEndCommandBuffer(tmpCommandBuffers[0]); result != C.VK_SUCCESS {
		return errors.Errorf("unable to record command buffer (result=%d)", result)
	}
	submitInfo := C.VkSubmitInfo{
		sType:              C.VK_STRUCTURE_TYPE_SUBMIT_INFO,
		commandBufferCount: C.uint(len(tmpCommandBuffers)),
		pCommandBuffers:    &tmpCommandBuffers[0],
	}
	submits := newVkSubmitInfoSlice(submitInfo)
	if result := C.vkQueueSubmit(*app.graphicsQueue, C.uint(len(submits)), &submits[0], nil); result != C.VK_SUCCESS {
		return errors.Errorf("unable to submit command buffers to graphics queue (result=%d)", result)
	}
	if result := C.vkQueueWaitIdle(*app.graphicsQueue); result != C.VK_SUCCESS {
		return errors.Errorf("unable to wait for graphics queue to become idle (result=%d)", result)
	}
	return nil
}

func createVertexBuffer(app *App, uniqueVertices []Vertex) error {
	// Create vertex staging buffer in CPU memory.
	vertexBufferSize := getVerticesSize(uniqueVertices)
	stagingBufferUsage := C.VkBufferUsageFlags(C.VK_BUFFER_USAGE_TRANSFER_SRC_BIT)
	stagingBufferProperties := C.VkMemoryPropertyFlags(C.VK_MEMORY_PROPERTY_HOST_VISIBLE_BIT | C.VK_MEMORY_PROPERTY_HOST_COHERENT_BIT)
	stagingBuffer, stagingBufferMem, err := createBuffer(app, vertexBufferSize, stagingBufferUsage, stagingBufferProperties)
	if err != nil {
		return errors.WithStack(err)
	}
	defer C.vkDestroyBuffer(*app.device, *stagingBuffer, nil)
	defer C.vkFreeMemory(*app.device, *stagingBufferMem, nil)
	if err := fillVertexBuffer(app, uniqueVertices, stagingBufferMem); err != nil {
		return errors.WithStack(err)
	}
	// Create vertex buffer in GPU memory.
	vertexBufferUsage := C.VkBufferUsageFlags(C.VK_BUFFER_USAGE_TRANSFER_DST_BIT | C.VK_BUFFER_USAGE_VERTEX_BUFFER_BIT)
	vertexBufferProperties := C.VkMemoryPropertyFlags(C.VK_MEMORY_PROPERTY_DEVICE_LOCAL_BIT)
	vertexBuffer, vertexBufferMem, err := createBuffer(app, vertexBufferSize, vertexBufferUsage, vertexBufferProperties)
	if err != nil {
		return errors.WithStack(err)
	}
	app.vertexBuffer = vertexBuffer
	app.vertexBufferMem = vertexBufferMem
	// Copy staging buffer to vertex buffer.
	if err := copyBuffer(app, vertexBuffer, stagingBuffer, vertexBufferSize); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func createIndexBuffer(app *App, indices []uint32) error {
	// Create index staging buffer in CPU memory.
	indexBufferSize := getIndicesSize(indices)
	stagingBufferUsage := C.VkBufferUsageFlags(C.VK_BUFFER_USAGE_TRANSFER_SRC_BIT)
	stagingBufferProperties := C.VkMemoryPropertyFlags(C.VK_MEMORY_PROPERTY_HOST_VISIBLE_BIT | C.VK_MEMORY_PROPERTY_HOST_COHERENT_BIT)
	stagingBuffer, stagingBufferMem, err := createBuffer(app, indexBufferSize, stagingBufferUsage, stagingBufferProperties)
	if err != nil {
		return errors.WithStack(err)
	}
	defer C.vkDestroyBuffer(*app.device, *stagingBuffer, nil)
	defer C.vkFreeMemory(*app.device, *stagingBufferMem, nil)
	if err := fillIndexBuffer(app, indices, stagingBufferMem); err != nil {
		return errors.WithStack(err)
	}
	// Create index buffer in GPU memory.
	indexBufferUsage := C.VkBufferUsageFlags(C.VK_BUFFER_USAGE_TRANSFER_DST_BIT | C.VK_BUFFER_USAGE_INDEX_BUFFER_BIT)
	indexBufferProperties := C.VkMemoryPropertyFlags(C.VK_MEMORY_PROPERTY_DEVICE_LOCAL_BIT)
	indexBuffer, indexBufferMem, err := createBuffer(app, indexBufferSize, indexBufferUsage, indexBufferProperties)
	if err != nil {
		return errors.WithStack(err)
	}
	app.indexBuffer = indexBuffer
	app.indexBufferMem = indexBufferMem
	// Copy staging buffer to index buffer.
	if err := copyBuffer(app, indexBuffer, stagingBuffer, indexBufferSize); err != nil {
		return errors.WithStack(err)
	}
	return nil
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

// uniqueIndexList returns a list of indices into a list of unique vertices,
// where uniqueVertices[indices[i]] == vertices[i] for each element.
func uniqueIndexList(vertices []Vertex) ([]uint32, []Vertex) {
	indices := make([]uint32, 0, len(vertices))
	var uniqueVertices []Vertex
	// maps from vertex to uniqueVertexIndex
	uniqueVertexIndexFromVertex := make(map[Vertex]uint32)
	for _, vertex := range vertices {
		uniqueVertexIndexFromVertex[vertex] = 0 // placeholder value.
	}
	for uniqueVertex := range uniqueVertexIndexFromVertex {
		uniqueVertices = append(uniqueVertices, uniqueVertex)
	}
	sort.Slice(uniqueVertices, func(i, j int) bool {
		return uniqueVertices[i].Less(uniqueVertices[j])
	})
	for uniqueVertexIndex, uniqueVertex := range uniqueVertices {
		uniqueVertexIndexFromVertex[uniqueVertex] = uint32(uniqueVertexIndex)
	}
	// Ensure vertices[i] ==
	for _, vertex := range vertices {
		uniqueVertexIndex := uniqueVertexIndexFromVertex[vertex]
		indices = append(indices, uniqueVertexIndex)
	}
	return indices, uniqueVertices
}
