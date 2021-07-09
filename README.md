# Laki

Sprickvulkaner, lek och grafikstacken.

## Dependencies

```bash
pacman -S vulkan-headers vulkan-tools vulkan-validation-layers glslang
pacman -S glfw-wayland
```

## Installation

```bash
git clone https://github.com/mewmew/laki
cd laki
go install -v ./...
```

## Usage

```bash
go run ./cmd/laki
```
