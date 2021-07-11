#version 450

// output to framebuffer index 0.
layout(location = 0) out vec3 fragColor;

vec2 positions[3] = vec2[](
	vec2(0.0, -0.5), // x, y
	vec2(0.5, 0.5),  // x, y
	vec2(-0.5, 0.5)  // x, y
);

vec3 colors[3] = vec3[](
	vec3(1.0, 0.0, 0.0), // red
	vec3(0.0, 1.0, 0.0), // green
	vec3(0.0, 0.0, 1.0)  // blue
);

// main called for every vertex.
void main() {
	gl_Position = vec4(positions[gl_VertexIndex], 0.0, 1.0); // xy, z, w
	fragColor = colors[gl_VertexIndex];
}
