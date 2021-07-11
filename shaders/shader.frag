#version 450

// input from framebuffer index 0.
layout(location = 0) in vec3 fragColor;

// output to framebuffer index 0.
layout(location = 0) out vec4 outColor;

// main called for every fragment.
void main() {
	outColor = vec4(fragColor, 1.0); // rgb, a
}
