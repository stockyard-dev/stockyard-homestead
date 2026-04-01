.PHONY: build run clean
build:
	CGO_ENABLED=0 go build -o homestead ./cmd/homestead/
run: build
	./homestead
clean:
	rm -f homestead
