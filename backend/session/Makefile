all: bin/session-service
test: lint unit-test
IMAGE_NAME=session-service
TAG=latest
PLATFORM?=linux/amd64

.PHONY: bin/session-service
bin/session-service:
	@docker build -f Dockerfile ../.. --target bin -t $(IMAGE_NAME):$(TAG) \
	--output bin/ \
	--platform ${PLATFORM}

.PHONY: bin/session-client
bin/session-client:
	@docker build -f Dockerfile ../.. --target client -t $(IMAGE_NAME):$(TAG) \
	--output bin/ \
	--platform ${PLATFORM}

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