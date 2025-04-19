all: bin/user-service
test: lint unit-test
IMAGE_NAME=user-service
TAG=latest
PLATFORM?=linux/amd64

.PHONY: bin/user-service
bin/user-service:
	@docker build . --target bin -t $(IMAGE_NAME):$(TAG) \
	--output bin/ \
	--platform ${PLATFORM}

.PHONY: bin/user-client
bin/user-client:
	@docker build . --target client -t $(IMAGE_NAME):$(TAG) \
	--output bin/ \
	--platform ${PLATFORM} --no-cache

.PHONY: unit-test
unit-test:
	@docker build . --target unit-test \
	$(BUILD_DIR)

.PHONY: unit-test-coverage
unit-test-coverage:
	@docker build . --target unit-test-coverage \
	$(BUILD_DIR) \
	--output coverage/ 
	cat coverage/cover.out

.PHONY: lint
lint:
	@docker build . --target lint
