all: laki liblaki.a

SRC=$(wildcard src/*.c)

OBJ=$(SRC:.c=.o)
A_OBJ=$(filter-out src/main.o,${OBJ})

LIBS=`pkg-config --libs glfw3`
LIBS+=`pkg-config --libs vulkan`

run: laki
	./laki

laki: liblaki.a
	go build -o $@ ./cmd/laki

liblaki.a: ${A_OBJ}
	ar rc $@ $^
	ranlib $@

%.o: %.c
	clang -c -o $@ $<

clean:
	$(RM) laki
	$(RM) liblaki.a
	$(RM) ${OBJ}
	$(RM) ${A_OBJ}

.PHONY: all clean
