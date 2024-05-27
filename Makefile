# Variables
DOCKER_IMAGE_NAME = gc-aws-agent
DOCKERFILE_PATH = Docker/Dockerfile
BROKER_URL = ec2-44-204-116-223.compute-1.amazonaws.com

# Targets
.PHONY: all build run

all: build run

build:
	docker buildx build -f $(DOCKERFILE_PATH) -t $(DOCKER_IMAGE_NAME) . --load

run:
	docker run -ti --network=host -e GCAPP_BROKER=$(BROKER_URL) --rm $(DOCKER_IMAGE_NAME)
