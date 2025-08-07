# Staging Environment Setup on VPS

This document provides instructions for setting up and managing the staging environment on a VPS.

## Prerequisites

- A VPS with at least 4GB RAM and 2 CPUs
- Ubuntu 22.04 LTS or similar Linux distribution
- Docker and Docker Compose installed
- Domain name with DNS configured to point to your VPS IP
- SSH access to the VPS

## Initial Server Setup

1. Update the system:
   ```bash
   sudo apt update && sudo apt upgrade -y
   ```

2. Install Docker and Docker Compose:
   ```bash
   # Install Docker
   curl -fsSL https://get.docker.com -o get-docker.sh
   sudo sh get-docker.sh
   
   # Add your user to the docker group
   sudo usermod -aG docker $USER
   
   # Install Docker Compose
   sudo apt install docker-compose-plugin -y
   ```

3. Clone the repository:
   ```bash
   git clone <your-repository-url> /opt/labs
   cd /opt/labs
   ```

## Configuration

1. Copy the `.env.staging` template and customize it:
   ```bash
   cp .env.staging .env.staging.local
   nano .env.staging.local
   ```

2. Update the following in `.env.staging.local`:
    - Replace `yourdomain.com` with your actual domain
    - Set strong passwords for MongoDB
    - Generate hashed passwords for basic auth services using:
      ```bash
      docker run --rm httpd:2.4-alpine htpasswd -nb admin yourpassword
      ```
    - Update the `ACME_EMAIL` with your email for Let's Encrypt notifications

3. Create service-specific `.env.staging` files:
   ```bash
   # Example for auth-service
   mkdir -p auth-service
   nano auth-service/.env.staging
   
   # Repeat for other services: pdf-service, labs-service, gateway-service, web, labs-worker
   ```

## Deployment

1. Make the deployment script executable:
   ```bash
   chmod +x deploy-staging.sh
   ```

2. Run the deployment script:
   ```bash
   ./deploy-staging.sh
   ```

3. Check the status of the services:
   ```bash
   docker compose -f docker-compose.staging.yaml ps
   ```

4. View logs:
   ```bash
   # All services
   docker compose -f docker-compose.staging.yaml logs -f
   
   # Specific service
   docker compose -f docker-compose.staging.yaml logs -f service-name
   ```

## SSL Certificates

The Traefik service will automatically obtain and renew SSL certificates from Let's Encrypt. To check the status:

```bash
docker compose -f docker-compose.staging.yaml exec traefik cat /letsencrypt/acme.json
```

## Maintenance

### Updating Services

To update the services:

1. Pull the latest changes:
   ```bash
   git pull
   ```

2. Run the deployment script:
   ```bash
   ./deploy-staging.sh
   ```

### Backup

1. Backup MongoDB data:
   ```bash
   docker compose -f docker-compose.staging.yaml exec mongo mongodump --out=/data/db/backup
   ```

2. Copy the backup to a safe location:
   ```bash
   docker cp $(docker compose -f docker-compose.staging.yaml ps -q mongo):/data/db/backup ./backup
   ```

### Monitoring

Access the Prometheus dashboard at `https://prometheus.staging.yourdomain.com` using the credentials set in your `.env.staging.local` file.

## Troubleshooting

### Check Container Logs

```bash
docker compose -f docker-compose.staging.yaml logs -f service-name
```

### Restart a Service

```bash
docker compose -f docker-compose.staging.yaml restart service-name
```

### Rebuild a Service

```bash
docker compose -f docker-compose.staging.yaml build service-name
docker compose -f docker-compose.staging.yaml up -d service-name
```

### Check Traefik Logs for SSL Issues

```bash
docker compose -f docker-compose.staging.yaml logs -f traefik
```

## Security Considerations

1. Ensure your VPS firewall allows only ports 80 and 443
2. Regularly update your system and Docker images
3. Use strong, unique passwords for all services
4. Consider setting up automated backups
5. Monitor system resources and logs regularly