# -s    disable symbol table
# -w    disable DWARF generation
BUILD_FLAGS := -s -w

build: phony bin
	go build -ldflags="${BUILD_FLAGS}" -o bin/ovo cmd/ovo.go

test: phony
	go test ./...

bin:
	mkdir -p bin

phony:
