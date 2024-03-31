# -s    disable symbol table
# -w    disable DWARF generation
BUILD_FLAGS := -s -w

# Usage: make run args="..."
.PHONY: run
run:
	eval "go run cmd/ovo.go ${args}"

.PHONY: build
build: bin
	go build -ldflags="${BUILD_FLAGS}" -o bin/ovo cmd/ovo.go

.PHONY: test
test:
	go test ./...

bin:
	mkdir -p bin
