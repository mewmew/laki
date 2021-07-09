package vk

// #include <stdint.h>
import "C"

// ref: VK_MAKE_API_VERSION
func VK_MAKE_API_VERSION(variant, major, minor, patch C.uint32_t) C.uint32_t {
	return (variant << 29) | (major << 22) | (minor << 12) | patch
}
