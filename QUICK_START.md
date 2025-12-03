# Event-Coming - Quick Start Guide

## Prerequisites
- Docker & Docker Compose
- Go 1.23+ (for local development)
- Make (optional, for convenience)

## 5-Minute Setup

### 1. Clone & Setup
```bash
git clone https://github.com/LeonardoGrigolettoDev/event-coming.git
cd event-coming
cp .env.example .env
```

### 2. Start Infrastructure
```bash
make docker-up
# or
docker-compose up -d
```

This starts:
- PostgreSQL with TimescaleDB & PostGIS (port 5432)
- Redis (port 6379)
- API Server (port 8080)
- Worker Processes

### 3. Check Health
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "event-coming"
}
```

## Development

### Build Locally
```bash
make build
# Builds to: bin/api and bin/worker
```

### Run API Server
```bash
make run
# or
go run cmd/api/main.go
```

### Run Workers
```bash
make run-worker
# or
go run cmd/worker/main.go
```

### Run Migrations
```bash
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/event_coming?sslmode=disable"
make migrate-up
```

## Project Structure at a Glance

```
event-coming/
├── cmd/                    # Application entry points
│   ├── api/               # API server
│   └── worker/            # Background workers
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── domain/           # Business domain models
│   ├── repository/       # Data access layer
│   ├── cache/            # Redis caching
│   ├── service/          # Business logic
│   ├── handler/          # HTTP handlers & middleware
│   ├── whatsapp/         # WhatsApp integration
│   └── router/           # Route definitions
├── pkg/                   # Public utilities
├── migrations/            # Database migrations
├── scripts/              # Helper scripts
└── docs/                 # Documentation
```

## Key Commands

```bash
# Build
make build              # Build both binaries
make run               # Run API server
make run-worker        # Run workers

# Test
make test              # Run tests
make test-coverage     # Generate coverage

# Database
make migrate-up        # Run migrations
make migrate-down      # Rollback migrations

# Docker
make docker-up         # Start all services
make docker-down       # Stop all services
make docker-logs       # View logs

# Quality
make lint              # Run linter
make tidy              # Tidy go modules

# Clean
make clean             # Remove build artifacts
```

## Configuration

Key environment variables (see `.env.example` for complete list):

```bash
# Database
EVENT_COMING_DATABASE_HOST=localhost
EVENT_COMING_DATABASE_PORT=5432
EVENT_COMING_DATABASE_USER=postgres
EVENT_COMING_DATABASE_PASSWORD=postgres

# Redis
EVENT_COMING_REDIS_HOST=localhost
EVENT_COMING_REDIS_PORT=6379

# JWT
EVENT_COMING_JWT_ACCESS_SECRET=your-secret-here
EVENT_COMING_JWT_REFRESH_SECRET=your-secret-here

# WhatsApp
EVENT_COMING_WHATSAPP_ACCESS_TOKEN=your-token
EVENT_COMING_WHATSAPP_PHONE_NUMBER_ID=your-id
```

## API Endpoints

### Public
- `POST /api/v1/auth/register` - Register user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh token

### Protected (requires JWT)
- `GET /api/v1/organizations` - List organizations
- `POST /api/v1/events` - Create event
- `GET /api/v1/events/:id` - Get event
- `POST /api/v1/events/:id/participants` - Add participant
- `POST /api/v1/participants/:id/locations` - Submit location
- `GET /api/v1/eta/events/:event_id` - Get ETAs

### Health
- `GET /health` - Health check

## Quick Examples

### Create Organization (Coming Soon)
```bash
curl -X POST http://localhost:8080/api/v1/organizations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "My School",
    "type": "school",
    "subscription_plan": "basic"
  }'
```

### Create Event (Coming Soon)
```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "School Bus Route",
    "type": "periodic",
    "location_lat": -23.5505,
    "location_lng": -46.6333,
    "start_time": "2024-12-04T07:00:00Z",
    "rrule_string": "RRULE:FREQ=DAILY;BYDAY=MO,TU,WE,TH,FR"
  }'
```

### Submit Location (Coming Soon)
```bash
curl -X POST http://localhost:8080/api/v1/participants/:id/locations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "latitude": -23.5505,
    "longitude": -46.6333,
    "accuracy": 10.5
  }'
```

## Database Access

```bash
# Connect to PostgreSQL
docker exec -it event_coming_postgres psql -U postgres event_coming

# Useful queries
SELECT * FROM organizations;
SELECT * FROM events;
SELECT * FROM participants;
SELECT * FROM locations ORDER BY time DESC LIMIT 10;
```

## Redis Access

```bash
# Connect to Redis
docker exec -it event_coming_redis redis-cli

# Check location buffer
LLEN location:buffer:<org_id>
LRANGE location:buffer:<org_id> 0 -1

# Check cached locations
GET location:latest:<event_id>:<participant_id>
```

## Troubleshooting

### API won't start
```bash
# Check logs
docker-compose logs api

# Check database connection
docker exec event_coming_postgres pg_isready

# Check Redis connection
docker exec event_coming_redis redis-cli ping
```

### Migrations fail
```bash
# Check current version
migrate -path migrations -database "$DATABASE_URL" version

# Force version (use carefully)
migrate -path migrations -database "$DATABASE_URL" force 2
```

### Build fails
```bash
# Clean and rebuild
make clean
go mod tidy
make build
```

## Project Status

### ✅ Complete
- Project structure
- Domain models
- Database schema
- Configuration system
- Middleware
- WhatsApp client
- ETA calculations
- Docker setup
- Documentation

### 🚧 In Progress
- Repository implementations
- Service layer
- HTTP handlers
- Background workers
- Test suite

### 📋 Planned
- Swagger documentation
- Admin dashboard
- Monitoring
- Analytics
- Mobile integration

## Learn More

- [README.md](README.md) - Comprehensive documentation
- [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) - Architecture details
- [IMPLEMENTATION_REPORT.md](IMPLEMENTATION_REPORT.md) - Implementation status

## Support

- Issues: https://github.com/LeonardoGrigolettoDev/event-coming/issues
- Docs: See documentation files in repo

## License

MIT License - See LICENSE file for details

---

**Ready to build something amazing!** 🚀
