// TODO: continue at https://vulkan-tutorial.com/en/Drawing_a_triangle/Setup/Physical_devices_and_queue_families

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

var REQUIRED_EXTENSIONS = []string{
	"VK_EXT_debug_utils",
}

var REQUIRED_LAYERS = []string{
	"VK_LAYER_KHRONOS_validation",
}

func Init() error {
	app := &App{}
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
	debug_messanger, err := initDebugMessanger(app.instance)
	if err != nil {
		return errors.WithStack(err)
	}
	app.debug_messanger = debug_messanger
	device, err := initDevice(app.instance)
	if err != nil {
		return errors.WithStack(err)
	}
	app.device = device
	return nil
}

func CleanupVulkan(app *App) {
	app.device = nil
	DestroyDebugUtilsMessengerEXT(*(app.instance), *(app.debug_messanger), nil)
	C.vkDestroyInstance(*(app.instance), nil)
}

func createInstance() (*C.VkInstance, error) {
	app_info := C.VkApplicationInfo{
		sType:              C.VK_STRUCTURE_TYPE_APPLICATION_INFO,
		pApplicationName:   C.CString(AppTitle),
		applicationVersion: VK_MAKE_API_VERSION(0, 1, 0, 0),
		pEngineName:        C.CString("No Engine"),
		engineVersion:      VK_MAKE_API_VERSION(0, 1, 0, 0),
		apiVersion:         C.VK_API_VERSION_1_0,
	}
	_ = app_info

	enabledExtensions := getExtensions()
	dbg.Println("nenabledExtensions:", len(enabledExtensions))
	for _, enabledExtension := range enabledExtensions {
		dbg.Println("   enabledExtension:", enabledExtension)
	}

	enabledLayers := getLayers()
	dbg.Println("nenabledLayers:", len(enabledLayers))
	for _, enabledLayer := range enabledLayers {
		dbg.Println("   enabledLayer:", enabledLayer)
	}

	create_info := C.new_VkInstanceCreateInfo()
	create_info.sType = C.VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO
	create_info.flags = 0 // reserved.
	create_info.pApplicationInfo = &app_info
	create_info.enabledExtensionCount = C.uint32_t(len(enabledExtensions))
	create_info.ppEnabledExtensionNames = getCStringSlice(enabledExtensions)
	create_info.enabledLayerCount = C.uint32_t(len(enabledLayers))
	create_info.ppEnabledLayerNames = getCStringSlice(enabledLayers)

	debug_messanger_create_info := C.new_VkDebugUtilsMessengerCreateInfoEXT()
	populateDebugMessangerCreateInfo(debug_messanger_create_info)
	create_info.pNext = unsafe.Pointer(debug_messanger_create_info)

	instance := C.new_VkInstance()
	result := C.vkCreateInstance(create_info, nil, instance)
	if result != C.VK_SUCCESS {
		return nil, errors.Errorf("unable to create Vulkan instance")
	}
	return instance, nil
}

func getExtensions() []string {
	// Get supported extensions.
	var nextensions C.uint32_t
	C.vkEnumerateInstanceExtensionProperties(nil, &nextensions, nil)
	extensions := make([]C.VkExtensionProperties, int(nextensions))
	C.vkEnumerateInstanceExtensionProperties(nil, &nextensions, &extensions[0])
	dbg.Println("nextensions:", len(extensions))
	var extensionNames []string
	for _, extension := range extensions {
		extensionName := C.GoString(&extension.extensionName[0])
		dbg.Println("   extension:", extensionName)
		extensionNames = append(extensionNames, extensionName)
	}

	// Get required extensions for GLFW.
	var nglfw_required_extensions C.uint32_t
	glfw_required_extensions := C.glfwGetRequiredInstanceExtensions(&nglfw_required_extensions)
	glfwRequiredExtensions := getStringSlice(unsafe.Pointer(glfw_required_extensions), int(nglfw_required_extensions))
	dbg.Println("nglfw_required_extensions:", len(glfwRequiredExtensions))
	for _, glfwRequiredExtension := range glfwRequiredExtensions {
		dbg.Println("   glfw_required_extension:", glfwRequiredExtension)
	}

	// Get required extensions by user.
	dbg.Println("NREQUIRED_EXTENSIONS:", len(REQUIRED_EXTENSIONS))
	for _, REQUIRED_EXTENSION := range REQUIRED_EXTENSIONS {
		dbg.Println("   REQUIRED_EXTENSION:", REQUIRED_EXTENSION)
	}

	// Check required extensions.
	var enabledExtensions []string
	for _, glfwRequiredExtension := range glfwRequiredExtensions {
		if !contains(extensionNames, glfwRequiredExtension) {
			warn.Printf("unable to locate required extension %q", glfwRequiredExtension)
			continue
		}
		enabledExtensions = append(enabledExtensions, glfwRequiredExtension)
	}
	for _, REQUIRED_EXTENSION := range REQUIRED_EXTENSIONS {
		if !contains(extensionNames, REQUIRED_EXTENSION) {
			warn.Printf("unable to locate required extension %q", REQUIRED_EXTENSION)
			continue
		}
		enabledExtensions = append(enabledExtensions, REQUIRED_EXTENSION)
	}

	return enabledExtensions
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
	dbg.Println("NREQUIRED_LAYERS:", len(REQUIRED_LAYERS))
	for _, REQUIRED_LAYER := range REQUIRED_LAYERS {
		dbg.Println("   REQUIRED_LAYER:", REQUIRED_LAYER)
	}

	// Check required layers.
	var enabledLayers []string
	for _, REQUIRED_LAYER := range REQUIRED_LAYERS {
		if !contains(layerNames, REQUIRED_LAYER) {
			warn.Printf("unable to locate required layer %q", REQUIRED_LAYER)
			continue
		}
		enabledLayers = append(enabledLayers, REQUIRED_LAYER)
	}

	return enabledLayers
}

func initDevice(instance *C.VkInstance) (*C.VkPhysicalDevice, error) {
	// TODO: rank physical devices by score if more than one is present. E.g.
	// prefer dedicated graphics card with capability for larger textures.
	//
	// ref: https://vulkan-tutorial.com/en/Drawing_a_triangle/Setup/Physical_devices_and_queue_families#page_Base-device-suitability-checks

	// Get physical devices.
	var ndevices C.uint32_t
	C.vkEnumeratePhysicalDevices(*instance, &ndevices, nil)
	if ndevices == 0 {
		return nil, errors.Errorf("unable to locate physical device (GPU)")
	}
	if ndevices > 1 {
		warn.Printf("multiple (%d) physical device (GPU) located; support for ranking physical devices not yet implemented", ndevices)
	}
	devices := make([]C.VkPhysicalDevice, int(ndevices))
	C.vkEnumeratePhysicalDevices(*instance, &ndevices, &devices[0])
	dbg.Println("ndevices:", len(devices))
	for _, device := range devices {
		if !isSuitableDevice(&device) {
			continue
		}
		_device := C.new_VkPhysicalDevice()
		*_device = device // allocate pointer on C heap.
		return _device, nil
	}
	return nil, errors.Errorf("unable to locate suitable physical device (GPU)")
}

func isSuitableDevice(device *C.VkPhysicalDevice) bool {
	// Get device properties.
	var device_properties C.VkPhysicalDeviceProperties
	C.vkGetPhysicalDeviceProperties(*device, &device_properties)
	deviceName := C.GoString(&device_properties.deviceName[0])
	dbg.Println("   deviceName:", deviceName)
	pretty.Println("   device_properties:", device_properties)

	// Get device features.
	var device_features C.VkPhysicalDeviceFeatures
	C.vkGetPhysicalDeviceFeatures(*device, &device_features)
	pretty.Println("   device_features:", device_features)

	// Find queue which supports graphics operations.
	queue_families := getQueueFamilies(device)
	if !hasQueueFlag(queue_families, C.VK_QUEUE_GRAPHICS_BIT) {
		return false
	}

	return true
}

func hasQueueFlag(queue_families []C.VkQueueFamilyProperties, queueFlags C.VkQueueFlags) bool {
	for _, queue_family := range queue_families {
		if queue_family.queueFlags&queueFlags == queueFlags {
			return true
		}
	}
	return false
}

func initDebugMessanger(instance *C.VkInstance) (*C.VkDebugUtilsMessengerEXT, error) {
	debug_messanger_create_info := C.new_VkDebugUtilsMessengerCreateInfoEXT()
	populateDebugMessangerCreateInfo(debug_messanger_create_info)
	debug_messenger := C.new_VkDebugUtilsMessengerEXT()
	result := CreateDebugUtilsMessengerEXT(*instance, debug_messanger_create_info, nil, debug_messenger)
	if result != C.VK_SUCCESS {
		return nil, errors.Errorf("unable to register debug messanger (result=%d)", result)
	}
	return debug_messenger, nil
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

func populateDebugMessangerCreateInfo(create_info *C.VkDebugUtilsMessengerCreateInfoEXT) {
	create_info.sType = C.VK_STRUCTURE_TYPE_DEBUG_UTILS_MESSENGER_CREATE_INFO_EXT
	create_info.messageSeverity = C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_VERBOSE_BIT_EXT | C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_INFO_BIT_EXT | C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_WARNING_BIT_EXT | C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_ERROR_BIT_EXT
	create_info.messageType = C.VK_DEBUG_UTILS_MESSAGE_TYPE_GENERAL_BIT_EXT | C.VK_DEBUG_UTILS_MESSAGE_TYPE_VALIDATION_BIT_EXT | C.VK_DEBUG_UTILS_MESSAGE_TYPE_PERFORMANCE_BIT_EXT
	//create_info.pfnUserCallback = (C.PFN_vkDebugUtilsMessengerCallbackEXT)(unsafe.Pointer(C.debug_callback))
	create_info.pfnUserCallback = (C.PFN_vkDebugUtilsMessengerCallbackEXT)(unsafe.Pointer(C.debugCallback))
	create_info.pUserData = nil // optional.
}

func getQueueFamilies(device *C.VkPhysicalDevice) []C.VkQueueFamilyProperties {
	var nqueue_families C.uint32_t
	C.vkGetPhysicalDeviceQueueFamilyProperties(*device, &nqueue_families, nil)
	queue_families := make([]C.VkQueueFamilyProperties, int(nqueue_families))
	C.vkGetPhysicalDeviceQueueFamilyProperties(*device, &nqueue_families, &queue_families[0])
	dbg.Println("nqueue_families:", len(queue_families))
	for _, queue_family := range queue_families {
		pretty.Println("   queue_family:", queue_family)
	}
	return queue_families
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
