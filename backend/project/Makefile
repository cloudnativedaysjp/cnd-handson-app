all: bin/project-service
test: lint unit-test
IMAGE_NAME=project-service
TAG=latest
PLATFORM?=linux/amd64

.PHONY: bin/project-service
bin/project-service:
	@docker build -f Dockerfile ../.. --target bin -t $(IMAGE_NAME):$(TAG) \
	--output bin/ \
	--platform ${PLATFORM}

.PHONY: bin/project-client
bin/project-client:
	@docker build -f Dockerfile ../.. --target client -t $(IMAGE_NAME):$(TAG) \
	--output bin/ \
	--platform ${PLATFORM} --no-cache

.PHONY: unit-test
unit-test:
	@docker build -f Dockerfile ../.. --target unit-test \
	$(BUILD_DIR)

.PHONY: unit-test-coverage
unit-test-coverage:
	@docker build -f Dockerfile ../.. --target unit-test-coverage \
	$(BUILD_DIR) \
	--output coverage/ 
	cat coverage/cover.out

.PHONY: lint
lint:
	@docker build -f Dockerfile ../.. --target lint

.PHONY: proto
proto:
	@$(MAKE) -C ../.. gen