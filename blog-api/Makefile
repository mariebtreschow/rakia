#!make
include .env
# Folder name
BINARY_NAME=blog-api

# Name of the Docker image
IMAGE_NAME=blog-api

# Docker tag for the image
TAG=1.0.0

# Docker build context
BUILD_CONTEXT=.

# Dockerfile location
DOCKERFILE_PATH=./Dockerfile

# Build the application
docker-build:
	@docker build -t $(IMAGE_NAME):$(TAG) -f $(DOCKERFILE_PATH) $(BUILD_CONTEXT)

# Run the docker image
docker-run:
	@docker run -p 8080:8080 -d -it --rm $(IMAGE_NAME):$(TAG) 

# Run the application locally
api:
	go run cmd/$(BINARY_NAME)/main.go

# Run all tests
tests:
	go test -v ./...