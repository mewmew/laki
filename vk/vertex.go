package vk

// #include <vulkan/vulkan.h>
import "C"

import "unsafe"

type Vertex struct {
	pos   Vec2
	color Vec3
}

// Less reports whether vertex a is less than vertex b, comparing each element
// of the struct depth first.
func (a Vertex) Less(b Vertex) bool {
	for i := range a.pos {
		switch {
		case a.pos[i] < b.pos[i]:
			return true
		case a.pos[i] > b.pos[i]:
			return false
		}
	}
	for i := range a.color {
		switch {
		case a.color[i] < b.color[i]:
			return true
		case a.color[i] > b.color[i]:
			return false
		}
	}
	return false
}

type Vec2 [2]float32

func vec2(x, y float32) Vec2 {
	return Vec2{x, y}
}

type Vec3 [3]float32

func vec3(x, y, z float32) Vec3 {
	return Vec3{x, y, z}
}

func getBindingDescs() ([]C.VkVertexInputBindingDescription, []C.VkVertexInputAttributeDescription) {
	dbg.Println("vk.getBindingDescs")
	const bindingNum = 0
	stride := C.uint(unsafe.Sizeof(Vertex{}))
	dbg.Println("   stride:", stride)
	bindingDescs := []C.VkVertexInputBindingDescription{
		{
			binding:   bindingNum,
			stride:    stride,
			inputRate: C.VK_VERTEX_INPUT_RATE_VERTEX,
		},
	}
	const (
		posLocationNum      = 0
		posColorLocationNum = 1
	)
	posOffset := C.uint(unsafe.Offsetof(Vertex{}.pos))
	colorOffset := C.uint(unsafe.Offsetof(Vertex{}.color))
	dbg.Println("   posOffset:", posOffset)
	dbg.Println("   colorOffset:", colorOffset)
	attrDescs := []C.VkVertexInputAttributeDescription{
		{
			location: posLocationNum,
			binding:  bindingNum,
			format:   C.VK_FORMAT_R32G32_SFLOAT,
			offset:   posOffset,
		},
		{
			location: posColorLocationNum,
			binding:  bindingNum,
			format:   C.VK_FORMAT_R32G32B32_SFLOAT,
			offset:   colorOffset,
		},
	}
	return bindingDescs, attrDescs
}
