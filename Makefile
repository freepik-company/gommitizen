BINARY_NAME=gommitizen
TAG=$(shell git describe --tags --always --dirty)

all: build install

build:
	@echo "Building..."
	@go build -o bin/$(BINARY_NAME) -v

docker:
	@echo "Building docker image..."
	@docker build -t $(BINARY_NAME):$(TAG) .
	@echo
	@echo "Done!"
	@echo "May the docker be with you..."
	@echo
	@echo "  # docker run -it $(BINARY_NAME):$(TAG) help"
	@echo
