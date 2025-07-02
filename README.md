# Consistency API

A RESTful API service for tracking habits, streaks, and achievements built with Go, Gin, and PostgreSQL.

## 🚀 Features

- **User Management**: Registration, authentication, and profile management
- **Habit Tracking**: Create, update, and delete habits
- **Streak Management**: Track consecutive habit completions
- **Check-ins**: Log habit completions with timestamps
- **Achievements**: Earn badges and rewards for milestones
- **RESTful API**: Clean, well-documented API endpoints
- **Health Checks**: Built-in health monitoring endpoints

## 🏗️ Architecture

- **Backend**: Go 1.23 with Gin web framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens
- **Logging**: Structured logging with zerolog
- **Deployment**: Docker containers on AWS EC2
- **CI/CD**: GitHub Actions for automated deployment

## 💰 Cost-Effective Deployment

**Monthly Costs (~$8-12/month):**
- EC2 t3.micro (free tier): $0-8/month
- RDS db.t4g.micro: $8/month
- Data transfer: ~$1/month
- **Total: $8-12/month** (92% cost reduction from previous ECS setup)

## 🛠️ Requirements

- Go 1.23+
- Docker & Docker Compose
- PostgreSQL (or use Docker)

## 🚀 Deployment

### AWS EC2 Deployment (Recommended)

See [EC2_DEPLOYMENT.md](EC2_DEPLOYMENT.md) for detailed deployment instructions.

**Automatic Deployment:**
- Push to `main` branch triggers GitHub Actions
- Builds Docker image and pushes to ECR
- Deploys to EC2 instance automatically

**Required GitHub Secrets:**
- `EC2_SSH_PRIVATE_KEY`: EC2 private key
- `EC2_HOST`: EC2 public IP or domain
- `EC2_USER`: EC2 username (usually `ec2-user`)
- `AWS_ACCESS_KEY_ID`: AWS access key
- `AWS_SECRET_ACCESS_KEY`: AWS secret key

## 🛠️ Local Development

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/WomenMobileDev/WMD.Consistency.Service.git
   cd WMD.Consistency.Service
   ```

2. **Run with Docker Compose**
   ```bash
   docker-compose up -d
   ```

3. **Access the API**
   ```bash
   curl http://localhost:8080/health
   ```

4. **View API Documentation**
   Open http://localhost:8080/swagger in your browser

### Manual Setup

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

## API Documentation

### Authentication Endpoints

- **Register User**
  - `POST /api/v1/auth/register`
  - Creates a new user account
  - Request body: `{"name": "string", "email": "string", "password": "string"}`

- **Login**
  - `POST /api/v1/auth/login`
  - Authenticates a user and returns a JWT token
  - Request body: `{"email": "string", "password": "string"}`

- **Forgot Password**
  - `POST /api/v1/auth/forgot-password`
  - Initiates password reset process
  - Request body: `{"email": "string"}`

- **Reset Password**
  - `POST /api/v1/auth/reset-password`
  - Resets user password using token
  - Request body: `{"token": "string", "password": "string"}`

### User Endpoints

- **Get User Profile**
  - `GET /api/v1/users/me`
  - Returns the current user's profile
  - Requires authentication

- **Update User Profile**
  - `PUT /api/v1/users/me`
  - Updates the current user's profile
  - Requires authentication
  - Request body: `{"name": "string"}`

### Habit Endpoints

- **List Habits**
  - `GET /api/v1/habits`
  - Returns all habits for the current user
  - Requires authentication

- **Create Habit**
  - `POST /api/v1/habits`
  - Creates a new habit
  - Requires authentication
  - Request body: `{"name": "string", "description": "string", "color": "#RRGGBB", "icon": "string"}`

- **Get Habit**
  - `GET /api/v1/habits/:id`
  - Returns a specific habit
  - Requires authentication

- **Update Habit**
  - `PUT /api/v1/habits/:id`
  - Updates a habit
  - Requires authentication
  - Request body: `{"name": "string", "description": "string", "color": "#RRGGBB", "icon": "string", "is_active": boolean}`

- **Delete Habit**
  - `DELETE /api/v1/habits/:id`
  - Deletes a habit
  - Requires authentication

### Streak Endpoints

- **Create Streak**
  - `POST /api/v1/habits/:id/streaks`
  - Creates a new streak for a habit
  - Requires authentication
  - Request body: `{"target_days": number}`

- **List Streaks**
  - `GET /api/v1/habits/:id/streaks`
  - Returns all streaks for a habit
  - Requires authentication

- **Get Current Streak**
  - `GET /api/v1/habits/:id/streaks/current`
  - Returns the current active streak for a habit
  - Requires authentication

### Check-in Endpoints

- **Check In**
  - `POST /api/v1/habits/:id/checkin`
  - Checks in for a habit for the current day
  - Requires authentication
  - Request body: `{"notes": "string"}`

- **List Check-ins**
  - `GET /api/v1/habits/:id/checkins`
  - Returns all check-ins for a habit
  - Requires authentication

### Achievement Endpoints

- **List Achievements**
  - `GET /api/v1/achievements`
  - Returns all achievements for the current user
  - Requires authentication

## Environment Variables

The application uses the following environment variables (defined in .env file):

- `PORT`: Server port (default: 8080)
- `ENV`: Environment (development, test, production)
- `LOG_LEVEL`: Logging level (debug, info, warn, error)
- `LOG_PRETTY`: Enable pretty logging (true/false)
- `DB_*`: Database connection parameters
- `AUTH_JWT_SECRET`: Secret key for JWT token generation
- `AUTH_JWT_EXPIRY_HOURS`: JWT token expiry in hours (default: 72)
- `AUTH_PASSWORD_RESET_EXPIRY`: Password reset token expiry (default: 24h)
- `AUTH_TOKEN_ISSUER`: JWT token issuer name

See `.env.example` for all available configuration options.

## Deployment (Staging & Production)

This project uses AWS ECS Fargate, RDS, and GitHub Actions for CI/CD. Two environments are supported: `staging` (branch: staging) and `production` (branch: main).

- See [docs/aws-setup.md](docs/aws-setup.md) for AWS resource setup instructions.
- On push to `main` or `staging`, GitHub Actions builds and deploys the app to the corresponding ECS service.
- Environment variables and secrets are managed via ECS and AWS Secrets Manager.
