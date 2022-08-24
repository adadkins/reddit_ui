.PHONY: default build image check publish-images

TAG_NAME := $(shell git tag -l --contains HEAD)

IMAGE_NAME := reddit_comment_api

default: check test build

build:
	CGO_ENABLED=0 go build -a --trimpath --installsuffix cgo --ldflags="-s" -o reddit_comment_api

build_pi:
	go build -o reddit_comment_api ./cmd/main.go
test:
	go test -v -cover ./...

check:
	golangci-lint run

image:
	docker buildx build --push --platform linux/arm/v7,linux/arm64/v8,linux/amd64,linux/arm/v6 --tag ghcr.io/adadkins/reddit_comment_api:latest .

publish-images:
	seihon publish -v "$(TAG_NAME)" -v "latest" --image-name $(IMAGE_NAME) --dry-run=false
