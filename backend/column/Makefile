all: bin/column-service
test: lint unit-test
IMAGE_NAME=column-service
TAG=latest
PLATFORM?=linux/amd64

.PHONY: bin/column-service
runtime:
	@docker build . --target runtime -t $(IMAGE_NAME):$(TAG)
.PHONY: lint
lint:
	@docker build . --target lint

.PHONY: proto
proto:
	@$(MAKE) -C ../.. gen

