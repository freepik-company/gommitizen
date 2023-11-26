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

install: build
	@echo "Installing..."
	@go install
	@echo
	@echo "Done!"

scan: start-sonar
	@echo "Scanning..."
	@curl -X POST -u admin:admin 'http://localhost:9000/api/users/create?login=user&password=password&name=user'
	@sonar-scanner \
		-Dsonar.projectKey=${BINARY_NAME} \
		-Dsonar.sources=. \
		-Dsonar.host.url=http://localhost:9000 \
		-Dsonar.login=user \
		-Dsonar.password=password
	@golangci-lint run
	@echo
	@echo "Done!"

start-sonar:
	@docker-compose up -d --wait sonarqube

stop-sonar:
	@docker-compose down

clean: stop-sonar
	@echo "Cleaning..."
	@rm -rf bin/*
	@echo
	@echo "Done!"
