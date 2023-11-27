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

docker-build:
	@docker build -t $(IMAGE_NAME):$(TAG) -f $(DOCKERFILE_PATH) $(BUILD_CONTEXT)

docker-run:
	@docker run -it --rm $(IMAGE_NAME):$(TAG)

api:
	go run cmd/$(BINARY_NAME)/main.go

tests:
	go test -v ./...