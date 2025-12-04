# Event-Coming

Event-Coming is a comprehensive event management and geolocation tracking system integrated with WhatsApp Cloud API. It supports both on-demand and recurring events with real-time participant tracking and ETA calculations.

## Features

- 🎯 **Multi-tenant Architecture**: Complete organization isolation with RBAC
- 🔐 **JWT Authentication**: Secure access with refresh tokens
- 📍 **Geolocation Tracking**: Real-time participant location tracking with TimescaleDB
- 💬 **WhatsApp Integration**: Cloud API for confirmations and notifications
- 📊 **ETA Calculations**: Multiple strategies (OSRM, velocity-based, simple)
- 🔄 **Recurring Events**: RRULE support for periodic events
- 🚀 **High Performance**: Redis caching and write-behind buffer pattern
- 🌐 **Real-time Updates**: WebSocket support for live location updates
- 📦 **Docker Ready**: Complete containerized deployment

## Tech Stack

- **Language**: Go 1.23+
- **Web Framework**: Gin
- **Database**: PostgreSQL with TimescaleDB (time-series) and PostGIS (geospatial)
- **Cache/Pub-Sub**: Redis
- **Messaging**: WhatsApp Cloud API
- **Logging**: Zap (structured logging)

## Architecture

```
event-coming/
├── cmd/
│   ├── api/          # API server entry point
│   └── worker/       # Background workers entry point
├── internal/
│   ├── config/       # Configuration management
│   ├── domain/       # Domain models and business logic
│   ├── repository/   # Data access layer
│   ├── cache/        # Redis cache layer
│   ├── service/      # Business services
│   ├── handler/      # HTTP handlers
│   ├── middleware/   # HTTP middleware
│   ├── whatsapp/     # WhatsApp client
│   ├── worker/       # Background workers
│   └── router/       # Route configuration
├── pkg/              # Shared utilities
├── migrations/       # Database migrations
└── docs/             # Documentation
```

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- PostgreSQL 15+ with TimescaleDB and PostGIS extensions
- Redis 7+
- Make (optional, for convenience commands)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/LeonardoGrigolettoDev/event-coming.git
cd event-coming
```

### 2. Setup Environment

```bash
cp .env.example .env
# Edit .env with your configuration
```

### 3. Start with Docker Compose

```bash
make docker-up
```

This will start:
- PostgreSQL with TimescaleDB (port 5432)
- Redis (port 6379)
- API server (port 8080)
- Worker processes

### 4. Run Migrations

```bash
make migrate-up
```

### 5. Test the API

```bash
curl http://localhost:8080/health
```

## Development

### Local Development Setup

#### Install Dependencies

```bash
go mod download
```

#### Install Development Tools

```bash
make install-tools
```

This installs:
- `swag` - Swagger documentation generator
- `golangci-lint` - Go linter
- `migrate` - Database migration tool

#### Run API Server

```bash
make run
```

#### Run Workers

```bash
make run-worker
```

### Database Migrations

#### Create New Migration

```bash
make migrate-create name=add_new_table
```

#### Run Migrations

```bash
make migrate-up
```

#### Rollback Migrations

```bash
make migrate-down
```

### Testing

#### Run Tests

```bash
make test
```

#### Run Tests with Coverage

```bash
make test-coverage
```

### Code Quality

#### Run Linter

```bash
make lint
```

#### Format Code

```bash
go fmt ./...
```

## Configuration

Configuration is managed through environment variables with the prefix `EVENT_COMING_`. See `.env.example` for all available options.

### Key Configuration Sections

#### Application
- `EVENT_COMING_APP_NAME`: Application name
- `EVENT_COMING_APP_ENVIRONMENT`: Environment (development/production)
- `EVENT_COMING_APP_DEBUG`: Debug mode

#### Database
- `EVENT_COMING_DATABASE_HOST`: PostgreSQL host
- `EVENT_COMING_DATABASE_PORT`: PostgreSQL port
- `EVENT_COMING_DATABASE_USER`: Database user
- `EVENT_COMING_DATABASE_PASSWORD`: Database password
- `EVENT_COMING_DATABASE_DATABASE`: Database name

#### Redis
- `EVENT_COMING_REDIS_HOST`: Redis host
- `EVENT_COMING_REDIS_PORT`: Redis port
- `EVENT_COMING_REDIS_PASSWORD`: Redis password (optional)

#### JWT
- `EVENT_COMING_JWT_ACCESS_SECRET`: Secret for access tokens
- `EVENT_COMING_JWT_REFRESH_SECRET`: Secret for refresh tokens
- `EVENT_COMING_JWT_ACCESS_TOKEN_TTL`: Access token TTL (default: 15m)
- `EVENT_COMING_JWT_REFRESH_TOKEN_TTL`: Refresh token TTL (default: 168h)

#### WhatsApp Cloud API
- `EVENT_COMING_WHATSAPP_ACCESS_TOKEN`: WhatsApp API access token
- `EVENT_COMING_WHATSAPP_PHONE_NUMBER_ID`: Phone number ID
- `EVENT_COMING_WHATSAPP_VERIFY_TOKEN`: Webhook verification token
- `EVENT_COMING_WHATSAPP_API_VERSION`: API version (default: v18.0)

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/forgot-password` - Request password reset
- `POST /api/v1/auth/reset-password` - Reset password

### Organizations
- `POST /api/v1/organizations` - Create organization
- `GET /api/v1/organizations/:id` - Get organization
- `PUT /api/v1/organizations/:id` - Update organization
- `GET /api/v1/organizations` - List organizations

### Events
- `POST /api/v1/events` - Create event
- `GET /api/v1/events/:id` - Get event
- `PUT /api/v1/events/:id` - Update event
- `DELETE /api/v1/events/:id` - Delete event
- `GET /api/v1/events` - List events

### Participants
- `POST /api/v1/events/:id/participants` - Add participant
- `GET /api/v1/events/:id/participants` - List participants
- `GET /api/v1/participants/:id` - Get participant
- `PUT /api/v1/participants/:id` - Update participant
- `DELETE /api/v1/participants/:id` - Remove participant

### Locations
- `POST /api/v1/participants/:id/locations` - Submit location
- `GET /api/v1/participants/:id/locations` - Get location history

### ETA
- `GET /api/v1/eta/events/:event_id` - Get ETAs for all participants
- `GET /api/v1/eta/participants/:participant_id` - Get ETA for participant

### WebSocket
- `GET /api/v1/ws/events/:event_id` - Real-time event updates

### Webhooks
- `GET /api/v1/webhook/whatsapp` - WhatsApp webhook verification
- `POST /api/v1/webhook/whatsapp` - WhatsApp webhook handler

## User Roles

- **super_admin**: Full system access
- **org_owner**: Organization owner with full access
- **org_admin**: Can manage events and users
- **org_manager**: Can manage events and participants
- **org_operator**: Can view and update events
- **org_viewer**: Read-only access

## WhatsApp Configuration

### 1. Create WhatsApp Business Account
1. Go to [Facebook for Developers](https://developers.facebook.com/)
2. Create a new app with WhatsApp product
3. Get your access token and phone number ID

### 2. Configure Webhook
1. Set webhook URL: `https://your-domain.com/api/v1/webhook/whatsapp`
2. Set verify token (same as `EVENT_COMING_WHATSAPP_VERIFY_TOKEN`)
3. Subscribe to `messages` webhook field

### 3. Create Message Templates
Create templates in WhatsApp Manager for:
- Event confirmations
- Location requests
- Reminders
- Event closure notifications

## Performance Considerations

### Location Buffer Pattern
Locations are first stored in Redis buffer, then batch-inserted into PostgreSQL:
1. Client submits location → Redis buffer
2. Worker periodically flushes buffer → PostgreSQL
3. Benefits: High write throughput, reduced DB load

### TimescaleDB Optimizations
- Hypertable partitioning by time
- Compression policy (data older than 7 days)
- Retention policy (90 days)
- Continuous aggregates for hourly summaries

### Redis Caching
- Latest participant locations cached
- Pub/sub for real-time updates
- Session management

## Monitoring

### Health Check
```bash
curl http://localhost:8080/health
```

### Metrics
Prometheus metrics are available at `/metrics` endpoint (when enabled).

## Production Deployment

### Build Production Image

```bash
docker build -t event-coming:latest .
```

### Run with Docker Compose

```bash
docker-compose up -d
```

### Database Backups

```bash
# Backup
docker exec event_coming_postgres pg_dump -U postgres event_coming > backup.sql

# Restore
docker exec -i event_coming_postgres psql -U postgres event_coming < backup.sql
```

## Troubleshooting

### Connection Issues
- Verify PostgreSQL is running and accessible
- Check Redis connection
- Verify environment variables are set correctly

### Migration Issues
```bash
# Check migration version
migrate -path migrations -database "$(DATABASE_URL)" version

# Force version (use with caution)
migrate -path migrations -database "$(DATABASE_URL)" force <version>
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support, email support@event-coming.com or open an issue on GitHub.

## Roadmap

- [ ] Complete handler implementations
- [ ] Add comprehensive test suite
- [ ] Implement OSRM routing integration
- [ ] Add Swagger documentation
- [ ] Implement audit logging
- [ ] Add metrics and monitoring
- [ ] Create admin dashboard
- [ ] Mobile app integration
- [ ] Multi-language support
- [ ] Advanced analytics

## References

- [WhatsApp Cloud API Documentation](https://developers.facebook.com/docs/whatsapp/cloud-api)
- [TimescaleDB Documentation](https://docs.timescale.com/)
- [PostGIS Documentation](https://postgis.net/documentation/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [OSRM Documentation](http://project-osrm.org/docs/v5.24.0/api/)
