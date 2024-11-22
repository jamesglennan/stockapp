# Variables
APP_NAME := "stockapp"
DOCKER_IMAGE := $(APP_NAME)
TAG := latest

build:
		go build -o main .

image:
		docker build -t $(DOCKER_IMAGE):$(TAG) .

push: image
		docker push $(DOCKER_IMAGE):$(TAG)

run:
		go run main.go

clean:
		rm -f main

.PHONY: build image push run clean