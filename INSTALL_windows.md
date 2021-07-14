# Installation (Windows)

## Dependencies

* Download and install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) (e.g. `tdm64-gcc-10.3.0-2.exe`)
* Download and install the [LunarG Vulkan SDK](https://vulkan.lunarg.com/) (e.g. `VulkanSDK-1.2.182.0-Installer.exe`)
* Download and extract [GLFW](https://www.glfw.org/download) (e.g. `glfw-3.3.4.bin.WIN64.zip`)

## Configure environment

```bash
set VULKAN_DIR=C:\VulkanSDK\1.2.182.0
set GLFW_DIR=C:\libs\glfw-3.3.4.bin.WIN64
set CGO_CFLAGS=-I %VULKAN_DIR%\Include -I %GLFW_DIR%\include
set CGO_LDFLAGS=-L %VULKAN_DIR%\Lib -L %GLFW_DIR%\lib-mingw-w64
set PATH=%PATH%;%GLFW_DIR%\lib-mingw-w64
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
