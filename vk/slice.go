// Go slices backed by C memory.

package vk

// #include <vulkan/vulkan.h>
//
// #include <stdlib.h>
//
// #include "malloc.h"
import "C"

import (
	"reflect"
	"unsafe"
)

func newVkPipelineSlice(elems ...C.VkPipeline) []C.VkPipeline {
	n := len(elems)
	data := C.new_VkPipelines(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkPipeline)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newVkAttachmentDescriptionSlice(elems ...C.VkAttachmentDescription) []C.VkAttachmentDescription {
	n := len(elems)
	data := C.new_VkAttachmentDescriptions(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkAttachmentDescription)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newVkAttachmentReferenceSlice(elems ...C.VkAttachmentReference) []C.VkAttachmentReference {
	n := len(elems)
	data := C.new_VkAttachmentReferences(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkAttachmentReference)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newVkSubpassDescriptionSlice(elems ...C.VkSubpassDescription) []C.VkSubpassDescription {
	n := len(elems)
	data := C.new_VkSubpassDescriptions(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkSubpassDescription)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newVkViewportSlice(elems ...C.VkViewport) []C.VkViewport {
	n := len(elems)
	data := C.new_VkViewports(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkViewport)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newVkRect2DSlice(elems ...C.VkRect2D) []C.VkRect2D {
	n := len(elems)
	data := C.new_VkRect2Ds(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkRect2D)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newVkPipelineColorBlendAttachmentStateSlice(elems ...C.VkPipelineColorBlendAttachmentState) []C.VkPipelineColorBlendAttachmentState {
	n := len(elems)
	data := C.new_VkPipelineColorBlendAttachmentStates(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkPipelineColorBlendAttachmentState)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newVkGraphicsPipelineCreateInfoSlice(elems ...C.VkGraphicsPipelineCreateInfo) []C.VkGraphicsPipelineCreateInfo {
	n := len(elems)
	data := C.new_VkGraphicsPipelineCreateInfos(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkGraphicsPipelineCreateInfo)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newVkFramebufferSlice(elems ...C.VkFramebuffer) []C.VkFramebuffer {
	n := len(elems)
	data := C.new_VkFramebuffers(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkFramebuffer)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newVkImageViewSlice(elems ...C.VkImageView) []C.VkImageView {
	n := len(elems)
	data := C.new_VkImageViews(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkImageView)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newVkCommandBufferSlice(elems ...C.VkCommandBuffer) []C.VkCommandBuffer {
	n := len(elems)
	data := C.new_VkCommandBuffers(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkCommandBuffer)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newVkClearValueSlice(elems ...C.VkClearValue) []C.VkClearValue {
	n := len(elems)
	data := C.new_VkClearValues(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkClearValue)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newVkSemaphoreSlice(elems ...C.VkSemaphore) []C.VkSemaphore {
	n := len(elems)
	data := C.new_VkSemaphores(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkSemaphore)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newVkSubmitInfoSlice(elems ...C.VkSubmitInfo) []C.VkSubmitInfo {
	n := len(elems)
	data := C.new_VkSubmitInfos(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkSubmitInfo)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newVkSubpassDependencySlice(elems ...C.VkSubpassDependency) []C.VkSubpassDependency {
	n := len(elems)
	data := C.new_VkSubpassDependencys(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkSubpassDependency)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newVkSwapchainKHRSlice(elems ...C.VkSwapchainKHR) []C.VkSwapchainKHR {
	n := len(elems)
	data := C.new_VkSwapchainKHRs(C.size_t(n))
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.VkSwapchainKHR)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}

func newCUint32Slice(elems ...C.uint32_t) []C.uint32_t {
	n := len(elems)
	const sizeof_uint32_t = 4
	data := C.calloc(C.size_t(n), sizeof_uint32_t)
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(data)),
		Len:  n,
		Cap:  n,
	}
	dst := *(*[]C.uint32_t)(unsafe.Pointer(&sh))
	for i := range elems {
		dst[i] = elems[i]
	}
	return dst
}
