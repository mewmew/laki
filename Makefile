all: laki

SRC=$(wildcard src/*.c)

laki: ${SRC}
	clang -o $@ $^

clean:
	$(RM) laki

.PHONY: all clean
