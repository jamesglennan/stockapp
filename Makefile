# Variables
APP_NAME := "ghcr.io/jamesglennan/stockapp"
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

deploy:
		kubectl apply -Rf deploy

.PHONY: build image push run clean deploy