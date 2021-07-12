# Installation (Linux)

## Dependencies

```bash
pacman -S vulkan-headers vulkan-tools vulkan-validation-layers glslang
pacman -S glfw-wayland
```

## Installation

Clone repository.
```bash
git clone https://github.com/mewmew/laki
cd laki
```

Compile using `make`
```bash
make
```

or use corresponding commands.
```bash
glslangValidator -V shaders/shader.vert -o shaders/shader_vert.spv
glslangValidator -V shaders/shader.frag -o shaders/shader_frag.spv
go build -v ./cmd/laki
```
