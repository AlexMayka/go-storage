# Go Storage - Docker Setup

This document describes how to run Go Storage using Docker and Docker Compose.

## Quick Start

1. **Copy environment file:**
   ```bash
   make env-copy
   # or manually: cp .env.example .env
   ```

2. **Start all services:**
   ```bash
   make quick-start
   # or manually: docker-compose up -d
   ```

3. **Access the application:**
   - **API**: http://localhost:8080
   - **Swagger Documentation**: http://localhost:8080/swagger/index.html
   - **MinIO Console**: http://localhost:9001 (admin/secret123)
   - **PostgreSQL**: localhost:5432 (admin/admin/storage)

## Services

### Application (Port 8080)
- Go-based REST API
- JWT authentication
- File upload/download with chunked support
- Swagger documentation

### PostgreSQL Database (Port 5432)
- Database for user data, companies, roles, and file metadata
- Auto-migration on startup
- Health checks enabled

### MinIO Object Storage (Ports 9000/9001)
- S3-compatible object storage for files
- Web console for management
- Automatic bucket creation

## Available Commands

```bash
# Development
make dev          # Start in development mode
make dev-down     # Stop development environment

# Production
make build        # Build all images
make up           # Start all services
make down         # Stop and remove containers
make restart      # Restart services

# Monitoring
make logs         # Show all logs
make logs-app     # Show app logs only
make logs-db      # Show database logs only
make logs-minio   # Show MinIO logs only
make status       # Show service status
make health       # Check service health

# Database
make db-shell     # Access PostgreSQL shell

# Cleanup
make clean        # Remove containers and volumes
make clean-all    # Remove everything including images
```

## Environment Configuration

Edit `.env` file to customize:

```bash
# Database
POSTGRES_USER=admin
POSTGRES_PASSWORD=admin
POSTGRES_DB=storage

# MinIO
MINIO_ROOT_USER=admin
MINIO_ROOT_PASSWORD=secret123

# Application
APP_JWT_SECRET=your-super-secret-jwt-key
APP_LOG_LEVEL=info

# File Server Limits
FILE_MAX_SIZE=5368709120           # 5GB
FILE_MAX_CONCURRENT_UPLOADS=10
FILE_CHUNK_SIZE=5242880            # 5MB
```

## Development Workflow

1. **Start dependencies only:**
   ```bash
   docker-compose up -d db minio
   ```

2. **Run app locally:**
   ```bash
   go run ./cmd/api
   ```

3. **Or rebuild and run in container:**
   ```bash
   make dev
   ```

## Production Deployment

1. **Generate SSL certificates:**
   ```bash
   ./scripts/generate-ssl.sh
   ```

2. **Use production compose file:**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

3. **Services include:**
   - Nginx reverse proxy with SSL
   - Rate limiting and security headers
   - Production-optimized settings

## File Upload Testing

Test file upload with curl:

```bash
# Small file upload
curl -X POST \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@test.txt" \
  -F "parentPath=/" \
  http://localhost:8080/api/v1/files/upload

# Check upload strategy
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  "http://localhost:8080/api/v1/files/upload-strategy?fileSize=1048576"
```

## Health Checks

All services include health checks:

```bash
# Check all services
make health

# Individual service checks
docker-compose exec app wget -q --spider http://localhost:8080/api/v1/users/register
docker-compose exec db pg_isready -U admin -d storage
docker-compose exec minio curl -f http://localhost:9000/minio/health/live
```

## Troubleshooting

### App won't start
- Check if database is ready: `make logs-db`
- Check MinIO connection: `make logs-minio`
- Verify environment variables: `docker-compose config`

### File uploads fail
- Check MinIO is running: `make logs-minio`
- Verify bucket creation in MinIO console
- Check app logs for errors: `make logs-app`

### Database connection issues
- Ensure PostgreSQL is healthy: `make health`
- Check migration logs: `make logs-db`
- Verify credentials in `.env`

### Performance issues
- Monitor resource usage: `docker stats`
- Adjust memory limits in docker-compose.yml
- Check file server settings in `.env`

## Security Considerations

### Development
- Default passwords are used
- SSL certificates are self-signed
- All ports exposed for debugging

### Production
- Change all default passwords
- Use proper SSL certificates
- Configure firewall rules
- Limit exposed ports
- Use secrets management
- Enable audit logging

## Data Persistence

Data is stored in Docker volumes:
- `pgdata`: PostgreSQL data
- `minio_data`: MinIO object storage
- `app_logs`: Application logs

To backup:
```bash
docker run --rm -v go-storage_pgdata:/data -v $(pwd):/backup alpine tar czf /backup/db-backup.tar.gz -C /data .
docker run --rm -v go-storage_minio_data:/data -v $(pwd):/backup alpine tar czf /backup/minio-backup.tar.gz -C /data .
```

## Monitoring

### Logs
```bash
# Real-time logs
docker-compose logs -f

# Specific service
docker-compose logs -f app

# Last 100 lines
docker-compose logs --tail=100 app
```

### Metrics
- Application exposes metrics at `/api/v1/files/stats`
- MinIO has built-in metrics in console
- PostgreSQL logs can be monitored

### Alerts
Set up monitoring for:
- Service health status
- Disk space (especially MinIO volumes)
- Memory and CPU usage
- Failed authentication attempts