#version 450

// input variables.
layout(location = 0) in vec2 inPosition;
layout(location = 1) in vec3 inColor;

// output to framebuffer index 0.
layout(location = 0) out vec3 fragColor;

// main called for every vertex.
void main() {
	gl_Position = vec4(inPosition, 0.0, 1.0); // xy, z, w
	fragColor = inColor;
}
