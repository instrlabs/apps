#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting deployment of staging environment...${NC}"

# Check if .env.staging exists
if [ ! -f .env.staging ]; then
    echo -e "${RED}Error: .env.staging file not found!${NC}"
    echo -e "Please create .env.staging file with appropriate values."
    exit 1
fi

# Create necessary directories
echo -e "${YELLOW}Creating necessary directories...${NC}"
mkdir -p traefik/letsencrypt
mkdir -p prometheus

# Check if prometheus.yml exists, if not create a basic one
if [ ! -f prometheus/prometheus.yml ]; then
    echo -e "${YELLOW}Creating basic prometheus.yml...${NC}"
    cat > prometheus/prometheus.yml << EOF
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'services'
    dns_sd_configs:
      - names:
        - 'auth-service'
        - 'pdf-service'
        - 'labs-service'
        - 'gateway-service'
        type: 'A'
        port: 3000
EOF
fi

# Pull latest changes from git if it's a git repository
if [ -d .git ]; then
    echo -e "${YELLOW}Pulling latest changes from git...${NC}"
    git pull
fi

# Build and start the services
echo -e "${YELLOW}Building and starting services...${NC}"
docker compose -f docker-compose.staging.yaml build
docker compose -f docker-compose.staging.yaml up -d

echo -e "${GREEN}Deployment completed successfully!${NC}"
echo -e "${YELLOW}Please check the logs with: docker compose -f docker-compose.staging.yaml logs -f${NC}"