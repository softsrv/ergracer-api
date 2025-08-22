# ErgRacer API

A rowing race application API built with Go and PostgreSQL that allows users to compete against each other using Concept2 rowing machines.

## Features

- **User Management**: Registration, authentication, and email verification
- **Friend System**: Invite friends after racing together
- **Race Management**: Create, join, and participate in rowing races
- **Real-time Racing**: Track progress and race completion
- **Race History**: View past races and statistics
- **JWT Authentication**: Secure API access

## Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL 17+
- Docker (optional)

### Configuration

The application uses YAML configuration. Copy `config.yaml.example` to `config.yaml` and update with your values:

```yaml
# Database Configuration
database:
  url: "postgres://username:password@localhost:5432/ergracer?sslmode=disable"

# JWT Configuration
jwt:
  secret: "your-super-secret-jwt-key-change-in-production"

# Mailgun Configuration (for email verification)
mailgun:
  domain: "yourdomain.mailgun.org"
  api_key: "your-mailgun-api-key"
  from_email: "noreply@yourdomain.com"
  from_name: "ErgRacer Team"

# Application Configuration
app:
  url: "http://localhost:8080"
  port: 8080
```

**Environment Variables:**

- `CONFIG_FILE_PATH` (optional) - Path to config file (defaults to `config.yaml`)

### Running Locally

1. Clone and setup:

```bash
git clone <repository>
cd ergracer-api
go mod download
```

2. Setup configuration:

```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your database and SMTP settings
```

3. Run the application:

**For development with hot reloading (recommended):**
note: this requires docker for the postgres container.

```bash
make fresh
```

## API Documentation

### Authentication

#### Register

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "username",
  "password": "password123"
}
```

#### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

Response:
{
  "access_token": "jwt_token_here",
  "refresh_token": "refresh_token_here",
  "user": {...}
}
```

#### Refresh Token

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "refresh_token_here"
}

Response:
{
  "access_token": "new_jwt_token",
  "refresh_token": "new_refresh_token"
}
```

#### Verify Email

```http
GET /api/v1/auth/verify-email?token=verification_token
```

### User Profile

#### Get Profile

```http
GET /api/v1/profile
Authorization: Bearer <jwt_token>
```

### Friends

#### Invite Friend

```http
POST /api/v1/friends/invite
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "friend_id": 123
}
```

#### Accept Friend Request

```http
POST /api/v1/friends/accept/123
Authorization: Bearer <jwt_token>
```

#### Get Friends List

```http
GET /api/v1/friends
Authorization: Bearer <jwt_token>
```

#### Get Pending Invitations

```http
GET /api/v1/friends/invitations
Authorization: Bearer <jwt_token>
```

### Races

#### Create Race

```http
POST /api/v1/races
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "distance": 2000
}
```

#### Join Race

```http
POST /api/v1/races/join
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "race_uuid": "race-uuid-here"
}
```

#### Get Race Details

```http
GET /api/v1/races/{uuid}
Authorization: Bearer <jwt_token>
```

#### Set Ready Status

```http
POST /api/v1/races/{raceId}/ready
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "ready": true
}
```

#### Update Race Progress

```http
POST /api/v1/races/{raceId}/progress
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "distance": 1500
}
```

#### Start Race (Admin/System)

```http
POST /api/v1/races/{raceId}/start
Authorization: Bearer <jwt_token>
```

### Race History

#### Get User Race History

```http
GET /api/v1/history
Authorization: Bearer <jwt_token>
```

## Race Flow

1. **Create Race**: User creates a race with specified distance
2. **Join Race**: Other users join using the race UUID
3. **Ready Up**: All participants mark themselves as ready
4. **Countdown**: 10-second countdown begins when all are ready
5. **Race Start**: Race becomes active, participants can submit progress
6. **Progress Updates**: Users submit their rowing distance
7. **Completion**: Users are marked finished when they reach the target distance
8. **Results**: Pace and positions calculated automatically

## Database Schema

### Users

- id, email, username, password_hash
- email_verified, email_verify_token
- created_at, updated_at

### Friendships

- user_id, friend_id, status (pending/accepted)
- created_at, accepted_at

### Races

- id, uuid, distance, status, created_by
- created_at, started_at, finished_at, countdown_at

### Race Participants

- race_id, user_id, status, current_distance
- finished_at, pace, position, joined_at

### Race Updates

- race_id, user_id, distance, timestamp

### Sessions

- user_id, refresh_token_hash, device_type
- user_agent, ip_address, expires_at
- created_at, updated_at

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o ergracer-api
```

### Linting

```bash
golangci-lint run
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details.
