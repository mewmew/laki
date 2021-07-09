// TODO: continue at https://vulkan-tutorial.com/en/Drawing_a_triangle/Setup/Physical_devices_and_queue_families

package vk

// #include "app.h"
// #include "malloc.h"
// #include "vk.h"
//
//#cgo CFLAGS: -I../src
//#cgo LDFLAGS: -llaki -L../
//#cgo pkg-config: glfw3
//#cgo pkg-config: vulkan
import "C"

import (
	"fmt"
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
	app := C.new_app()
	app.win = InitWindow()
	defer CleanupWindow(app.win)
	if err := InitVulkan(app); err != nil {
		return errors.WithStack(err)
	}
	defer CleanupVulkan(app)

	EventLoop(app.win)
	return nil
}

func InitVulkan(app *C.App) error {
	instance, err := createInstance()
	if err != nil {
		return errors.WithStack(err)
	}
	app.instance = instance
	app.debug_messanger = C.init_debug_messanger(app.instance)
	device, err := initDevice(app.instance)
	if err != nil {
		return errors.WithStack(err)
	}
	app.device = device
	return nil
}

func CleanupVulkan(app *C.App) {
	app.device = nil
	C.DestroyDebugUtilsMessengerEXT(*(app.instance), *(app.debug_messanger), nil)
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
	fmt.Println("nenabledExtensions:", len(enabledExtensions))
	for _, enabledExtension := range enabledExtensions {
		fmt.Println("   enabledExtension:", enabledExtension)
	}

	enabledLayers := getLayers()
	fmt.Println("nenabledLayers:", len(enabledLayers))
	for _, enabledLayer := range enabledLayers {
		fmt.Println("   enabledLayer:", enabledLayer)
	}

	create_info := C.new_VkInstanceCreateInfo()
	create_info.sType = C.VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO
	create_info.flags = 0 // reserved.
	create_info.pApplicationInfo = &app_info
	create_info.enabledExtensionCount = C.uint32_t(len(enabledExtensions))
	create_info.ppEnabledExtensionNames = getCStringSlice(enabledExtensions)
	create_info.enabledLayerCount = C.uint32_t(len(enabledLayers))
	create_info.ppEnabledLayerNames = getCStringSlice(enabledLayers)

	//VkDebugUtilsMessengerCreateInfoEXT *debug_messanger_create_info = calloc(1, sizeof(VkDebugUtilsMessengerCreateInfoEXT));
	//populate_debug_messanger_create_info(debug_messanger_create_info);
	//create_info.pNext = (VkDebugUtilsMessengerCreateInfoEXT*)debug_messanger_create_info;

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
	fmt.Println("nextensions:", len(extensions))
	var extensionNames []string
	for _, extension := range extensions {
		extensionName := C.GoString(&extension.extensionName[0])
		fmt.Println("   extension:", extensionName)
		extensionNames = append(extensionNames, extensionName)
	}

	// Get required extensions for GLFW.
	var nglfw_required_extensions C.uint32_t
	glfw_required_extensions := C.glfwGetRequiredInstanceExtensions(&nglfw_required_extensions)
	glfwRequiredExtensions := getStringSlice(unsafe.Pointer(glfw_required_extensions), int(nglfw_required_extensions))
	fmt.Println("nglfw_required_extensions:", len(glfwRequiredExtensions))
	for _, glfwRequiredExtension := range glfwRequiredExtensions {
		fmt.Println("   glfw_required_extension:", glfwRequiredExtension)
	}

	// Get required extensions by user.
	fmt.Println("NREQUIRED_EXTENSIONS:", len(REQUIRED_EXTENSIONS))
	for _, REQUIRED_EXTENSION := range REQUIRED_EXTENSIONS {
		fmt.Println("   REQUIRED_EXTENSION:", REQUIRED_EXTENSION)
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
	fmt.Println("nlayers:", len(layers))
	var layerNames []string
	for _, layer := range layers {
		layerName := C.GoString(&layer.layerName[0])
		layerDesc := C.GoString(&layer.description[0])
		fmt.Println("   layer:", layerName)
		fmt.Println("      desc:", layerDesc)
		layerNames = append(layerNames, layerName)
	}

	// Get required layers by user.
	fmt.Println("NREQUIRED_LAYERS:", len(REQUIRED_LAYERS))
	for _, REQUIRED_LAYER := range REQUIRED_LAYERS {
		fmt.Println("   REQUIRED_LAYER:", REQUIRED_LAYER)
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
	devices := make([]C.VkPhysicalDevice, int(ndevices))
	C.vkEnumeratePhysicalDevices(*instance, &ndevices, &devices[0])
	fmt.Println("ndevices:", len(devices))
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
	fmt.Println("   deviceName:", deviceName)
	pretty.Println("   device_properties:", device_properties)

	// Get device features.
	var device_features C.VkPhysicalDeviceFeatures
	C.vkGetPhysicalDeviceFeatures(*device, &device_features)
	pretty.Println("   device_features:", device_features)

	return true
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
