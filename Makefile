all: laki

SRC=$(wildcard src/*.c)

LIBS=`pkg-config --libs glfw3`
LIBS+=`pkg-config --libs vulkan`

laki: ${SRC}
	clang -o $@ $^ ${LIBS}

clean:
	$(RM) laki

.PHONY: all clean
