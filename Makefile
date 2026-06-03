.PHONY: clean build

clean:
	@echo "Cleaning up..."
	@rm vkc

build: clean
	@echo "Building vkc..."
	@go build -o vkc ./cmd/vkc/main.go
