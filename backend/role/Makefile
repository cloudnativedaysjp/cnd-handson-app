all: bin/role-service
test: lint unit-test
IMAGE_NAME=role-service
TAG=latest
PLATFORM?=linux/amd64

.PHONY: bin/role-service
runtime:
	@docker build . --target runtime -t $(IMAGE_NAME):$(TAG)
.PHONY: lint
lint:
	@docker build . --target lint

.PHONY: proto
proto:
	@$(MAKE) -C ../.. gen