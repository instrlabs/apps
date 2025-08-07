# Docker commands for local environment
docker-up-local:
	docker-compose -f docker-compose.local.yaml up -d

docker-down-local:
	docker-compose -f docker-compose.local.yaml down

# Docker commands for staging environment
docker-up-staging:
	docker-compose -f docker-compose.staging.yaml up -d

docker-down-staging:
	docker-compose -f docker-compose.staging.yaml down

.PHONY: docker-up-local docker-down-local docker-up-staging docker-down-staging