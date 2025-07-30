project_name = labs-auth-service
image_name = histweety/labs-auth-service:latest

run-local: ## Run the app locally
	go run main.go

requirements: ## Generate go.mod & go.sum files
	go mod tidy

clean-packages: ## Clean packages
	go clean -modcache

up:
	docker run -d -p 3000:3000 --name $(project_name) $(image_name) ./app -prod

build: ## Generate docker image
	docker build -t $(image_name) .

push:
	docker push $(image_name)

shell: ## Run interactive shell in the container
	docker exec -it $(project_name) /bin/sh
