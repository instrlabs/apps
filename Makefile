.PHONY: help build-all build-gateway build-auth build-image build-notification build-pdf build-product up up-gateway up-auth up-image up-notification up-pdf up-product down push-all push-gateway push-auth push-image push-notification push-pdf push-product clean dev-mode prod-mode

# Get the short commit hash
COMMIT_HASH := $(shell git rev-parse --short HEAD)

# Image names
REGISTRY = histweety
GATEWAY_IMAGE = $(REGISTRY)/instrlabs-gateway-service
AUTH_IMAGE = $(REGISTRY)/instrlabs-auth-service
IMAGE_SERVICE_IMAGE = $(REGISTRY)/instrlabs-image-service
NOTIFICATION_IMAGE = $(REGISTRY)/instrlabs-notification-service
PDF_SERVICE_IMAGE = $(REGISTRY)/instrlabs-pdf-service
PRODUCT_SERVICE_IMAGE = $(REGISTRY)/instrlabs-product-service

# Default target
help:
	@echo "Available targets:"
	@echo "  build-all           - Build all service images"
	@echo "  build-gateway       - Build gateway service image"
	@echo "  build-auth          - Build auth service image"
	@echo "  build-image         - Build image service image"
	@echo "  build-notification  - Build notification service image"
	@echo "  build-pdf           - Build PDF service image"
	@echo "  build-product       - Build product service image"
	@echo ""
	@echo "  up                  - Start all services with Docker Compose"
	@echo "  up-gateway          - Start gateway service"
	@echo "  up-auth             - Start auth service"
	@echo "  up-image            - Start image service"
	@echo "  up-notification     - Start notification service"
	@echo "  up-pdf              - Start PDF service"
	@echo "  up-product          - Start product service"
	@echo "  down                - Stop all services"
	@echo ""
	@echo "  push-all            - Push all service images"
	@echo "  push-gateway        - Push gateway service image"
	@echo "  push-auth           - Push auth service image"
	@echo "  push-image          - Push image service image"
	@echo "  push-notification   - Push notification service image"
	@echo "  push-pdf            - Push PDF service image"
	@echo "  push-product        - Push product service image"
	@echo ""
	@echo "  clean               - Remove all built images"
	@echo ""
	@echo "Mode Switching:"
	@echo "  dev-mode            Switch to development mode (use local shared library)"
	@echo "  prod-mode           Switch to production mode (use tagged shared library)"
	@echo ""
	@echo "Current commit hash: $(COMMIT_HASH)"
	@echo "Registry: $(REGISTRY)"

# Build all services
build-all: build-gateway build-auth build-image build-notification build-pdf build-product

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

build-pdf:
	@echo "Building PDF service with commit hash $(COMMIT_HASH)..."
	docker build -t $(PDF_SERVICE_IMAGE):$(COMMIT_HASH) -t $(PDF_SERVICE_IMAGE):latest ./pdf-service
	@echo "Built: $(PDF_SERVICE_IMAGE):$(COMMIT_HASH)"

build-product:
	@echo "Building product service with commit hash $(COMMIT_HASH)..."
	docker build -t $(PRODUCT_SERVICE_IMAGE):$(COMMIT_HASH) -t $(PRODUCT_SERVICE_IMAGE):latest ./product-service
	@echo "Built: $(PRODUCT_SERVICE_IMAGE):$(COMMIT_HASH)"

# Push all services
push-all: push-gateway push-auth push-image push-notification push-pdf push-product

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

push-pdf:
	@echo "Pushing PDF service..."
	docker push $(PDF_SERVICE_IMAGE):$(COMMIT_HASH)
	docker push $(PDF_SERVICE_IMAGE):latest

push-product:
	@echo "Pushing product service..."
	docker push $(PRODUCT_SERVICE_IMAGE):$(COMMIT_HASH)
	docker push $(PRODUCT_SERVICE_IMAGE):latest

# Docker Compose commands
up:
	@echo "Starting all services..."
	docker compose -p instrlabs up -d

up-gateway:
	@echo "Starting gateway service..."
	docker compose -p instrlabs up -d gateway-service

up-auth:
	@echo "Starting auth service..."
	docker compose -p instrlabs up -d --build auth-service

up-image:
	@echo "Starting image service..."
	docker compose -p instrlabs up -d image-service

up-notification:
	@echo "Starting notification service..."
	docker compose -p instrlabs up -d notification-service

up-pdf:
	@echo "Starting PDF service..."
	docker compose -p instrlabs up -d pdf-service

up-product:
	@echo "Starting product service..."
	docker compose -p instrlabs up -d product-service

down:
	@echo "Stopping all services..."
	docker compose -p instrlabs down

# Clean up images
clean:
	@echo "Removing images for commit hash $(COMMIT_HASH)..."
	-docker rmi $(GATEWAY_IMAGE):$(COMMIT_HASH) 2>/dev/null || true
	-docker rmi $(AUTH_IMAGE):$(COMMIT_HASH) 2>/dev/null || true
	-docker rmi $(IMAGE_SERVICE_IMAGE):$(COMMIT_HASH) 2>/dev/null || true
	-docker rmi $(NOTIFICATION_IMAGE):$(COMMIT_HASH) 2>/dev/null || true
	@echo "Cleanup complete"

# Mode switching
dev-mode: ## Switch to development mode (use local shared library)
	@./scripts/dev-mode.sh

prod-mode: ## Switch to production mode (use tagged shared library)
	@./scripts/prod-mode.sh
