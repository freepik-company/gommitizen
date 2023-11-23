BINARY_NAME=fc-version

all: build install

build:
	@echo "Building..."
	@go build -o bin/$(BINARY_NAME) -v

docker:
	@echo "Building docker image..."
	@docker build -t $(BINARY_NAME) .
