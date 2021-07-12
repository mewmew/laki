all: laki

#VERT_SHADERS=$(wildcard shaders/*.vert)
#FRAG_SHADERS=$(wildcard shaders/*.frag)

shaders/%_vert.spv: shaders/%.vert
	glslangValidator -V $< -o $@

shaders/%_frag.spv: shaders/%.frag
	glslangValidator -V $< -o $@

laki: shaders/shader_vert.spv shaders/shader_frag.spv
	go build -v ./cmd/laki

run: laki
	./laki

clean:
	$(RM) laki

.PHONY: all clean
