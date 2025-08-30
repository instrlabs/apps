# Docker commands for local environment
docker-up-local:
	docker compose -f docker-compose.local.yaml up -d $(service)

docker-down-local:
	docker compose -f docker-compose.local.yaml down $(service)

docker-watch-local:
	docker compose -f docker-compose.local.yaml watch $(service)

# Docker commands for staging environment
docker-up-staging:
	docker compose -f docker-compose.staging.yaml up -d

docker-down-staging:
	docker compose -f docker-compose.staging.yaml down

clean-env-staging:
	find . -name ".env" -type f -delete
	find . -name ".env.staging" -type f -exec sh -c 'cp {} $$(dirname {})"/.env"' \;

.PHONY: docker-up-local docker-down-local docker-up-staging docker-down-staging docker-down-service-local clean-env-staging