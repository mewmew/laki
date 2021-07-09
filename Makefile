all: laki

laki:
	go build -o $@ ./cmd/laki

run: laki
	./laki

clean:
	$(RM) laki

.PHONY: all clean
