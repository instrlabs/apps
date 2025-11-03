.PHONY: help build-all build-gateway build-auth build-image build-notification push-all push-gateway push-auth push-image push-notification clean

# Get the short commit hash
COMMIT_HASH := $(shell git rev-parse --short HEAD)

# Image names
REGISTRY = histweety
GATEWAY_IMAGE = $(REGISTRY)/instrlabs-gateway-service
AUTH_IMAGE = $(REGISTRY)/instrlabs-auth-service
IMAGE_SERVICE_IMAGE = $(REGISTRY)/instrlabs-image-service
NOTIFICATION_IMAGE = $(REGISTRY)/instrlabs-notification-service

# Default target
help:
	@echo "Available targets:"
	@echo "  build-all           - Build all service images"
	@echo "  build-gateway       - Build gateway service image"
	@echo "  build-auth          - Build auth service image"
	@echo "  build-image         - Build image service image"
	@echo "  build-notification  - Build notification service image"
	@echo ""
	@echo "  push-all            - Push all service images"
	@echo "  push-gateway        - Push gateway service image"
	@echo "  push-auth           - Push auth service image"
	@echo "  push-image          - Push image service image"
	@echo "  push-notification   - Push notification service image"
	@echo ""
	@echo "  clean               - Remove all built images"
	@echo ""
	@echo "Current commit hash: $(COMMIT_HASH)"
	@echo "Registry: $(REGISTRY)"

# Build all services
build-all: build-gateway build-auth build-image build-notification

# Build individual services
build-gateway:
	@echo "Building gateway service with commit hash $(COMMIT_HASH)..."
	docker build -t $(GATEWAY_IMAGE):$(COMMIT_HASH) -t $(GATEWAY_IMAGE):latest ./gateway-service
	@echo "Built: $(GATEWAY_IMAGE):$(COMMIT_HASH)"

build-auth:
	@echo "Building auth service with commit hash $(COMMIT_HASH)..."
	docker build -t $(AUTH_IMAGE):$(COMMIT_HASH) -t $(AUTH_IMAGE):latest ./auth-service
	@echo "Built: $(AUTH_IMAGE):$(COMMIT_HASH)"

build-image:
	@echo "Building image service with commit hash $(COMMIT_HASH)..."
	docker build -t $(IMAGE_SERVICE_IMAGE):$(COMMIT_HASH) -t $(IMAGE_SERVICE_IMAGE):latest ./image-service
	@echo "Built: $(IMAGE_SERVICE_IMAGE):$(COMMIT_HASH)"

build-notification:
	@echo "Building notification service with commit hash $(COMMIT_HASH)..."
	docker build -t $(NOTIFICATION_IMAGE):$(COMMIT_HASH) -t $(NOTIFICATION_IMAGE):latest ./notification-service
	@echo "Built: $(NOTIFICATION_IMAGE):$(COMMIT_HASH)"

# Push all services
push-all: push-gateway push-auth push-image push-notification

# Push individual services
push-gateway:
	@echo "Pushing gateway service..."
	docker push $(GATEWAY_IMAGE):$(COMMIT_HASH)
	docker push $(GATEWAY_IMAGE):latest

push-auth:
	@echo "Pushing auth service..."
	docker push $(AUTH_IMAGE):$(COMMIT_HASH)
	docker push $(AUTH_IMAGE):latest

push-image:
	@echo "Pushing image service..."
	docker push $(IMAGE_SERVICE_IMAGE):$(COMMIT_HASH)
	docker push $(IMAGE_SERVICE_IMAGE):latest

push-notification:
	@echo "Pushing notification service..."
	docker push $(NOTIFICATION_IMAGE):$(COMMIT_HASH)
	docker push $(NOTIFICATION_IMAGE):latest

# Clean up images
clean:
	@echo "Removing images for commit hash $(COMMIT_HASH)..."
	-docker rmi $(GATEWAY_IMAGE):$(COMMIT_HASH) 2>/dev/null || true
	-docker rmi $(AUTH_IMAGE):$(COMMIT_HASH) 2>/dev/null || true
	-docker rmi $(IMAGE_SERVICE_IMAGE):$(COMMIT_HASH) 2>/dev/null || true
	-docker rmi $(NOTIFICATION_IMAGE):$(COMMIT_HASH) 2>/dev/null || true
	@echo "Cleanup complete"
