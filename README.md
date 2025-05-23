# WMD.Consistency.Service

A habit tracking backend service built with Go and Gin framework. This service helps users track their habits and maintain consistency streaks.

## Features (Planned)
- Set habit tracking goals with customizable durations (up to 7 days)
- Daily check-ins to mark habit completion
- Streak reset on missed check-ins
- RESTful API endpoints

## Requirements
- Go 1.22.2 or higher
- Docker and Docker Compose (for containerized deployment)
- PostgreSQL (for database, included in Docker Compose)

## Setup

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/WomenMobileDev/WMD.Consistency.Service.git
cd WMD.Consistency.Service
```

2. Copy the example environment file:
```bash
cp .env.example .env
```

3. Install dependencies:
```bash
go mod tidy
```

4. Run the server:
```bash
go run cmd/main.go
```

The server will start on port 8080 (or the port specified in your .env file). You can test the health endpoint at:
```
GET http://localhost:8080/health
```

### Development with Live Reloading (Air)

This project supports live reloading using [Air](https://github.com/cosmtrek/air), which automatically rebuilds and restarts the application when code changes are detected.

1. Install Air:
```bash
./scripts/install_air.sh
# Or use the Makefile
make install-air
```

2. Run the application with Air:
```bash
air -c .air.toml
# Or use the Makefile
make dev
```

3. Test live reloading by making changes to any Go file. The application will automatically rebuild and restart.

4. You can test the live reloading endpoint at:
```
GET http://localhost:8080/test-reload
```

### Docker Development with Live Reloading

You can also use Docker for development with live reloading:

1. Start the development environment with Docker:
```bash
docker-compose -f docker-compose.dev.yml up --build
# Or use the Makefile
make dev-docker
```

2. The application code is mounted as a volume, so any changes you make to the code will trigger a rebuild and restart of the application.

3. Test the live reloading by making changes to the code and accessing the test endpoint:
```
GET http://localhost:8080/test-reload
```

### Production Docker Deployment

1. Build and start the containers for production:
```bash
docker-compose up -d
# Or use the Makefile
make docker-prod
```

2. Check container status:
```bash
docker-compose ps
```

3. View logs:
```bash
docker-compose logs -f app
```

4. Stop the containers:
```bash
docker-compose down
# Or use the Makefile
make docker-stop
```

## Environment Variables

The application uses the following environment variables (defined in .env file):

- `PORT`: Server port (default: 8080)
- `ENV`: Environment (development, test, production)
- `LOG_LEVEL`: Logging level (debug, info, warn, error)
- `LOG_PRETTY`: Enable pretty logging (true/false)
- `DB_*`: Database connection parameters

See `.env.example` for all available configuration options.
