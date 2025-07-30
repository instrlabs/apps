project_name = labs-web
image_name = histweety/labs-web:latest

.PHONY: build
build: ## Generate docker image
	docker build -t $(image_name) .

.PHONY: push
push: ## Push docker image
	docker push $(image_name)
