# :shield: GoATED Pentest Tracker 

## Features
- Create and log Pentests for distinct team
- Team memberships and separation of ownership
- Findings to Pentest mapping and authorization checks
- Closure of Findings and Pentest Storage for Archiving

## Requirements
- go 1.26.1 or > 
- sqlc v1.31.1
- goose version: v3.27.0
- Postgresql 15 or higher
- github.com/alexedwards/argon2id v1.0.0
- github.com/golang-jwt/jwt/v5 v5.3.1
- github.com/google/uuid v1.6.0
- github.com/joho/godotenv v1.5.1
- github.com/lib/pq v1.12.3
- golang.org/x/crypto v0.14.0
- golang.org/x/sys v0.13.0

## Configuration / ENV
- In the .env file you need to specifiy the following variables:
    - DB_URL: Access url to your postgresql instance
    - JWT_SECRET: used for signing the tokens
    - ADMIN_SECRET: used to access the admin endpoints

## Installation and usage
```bash
git clone https://github.com/RoshiSecOps/goATED-tracker \
goose -dir sql/schema postgres "your-db-url" up # initialize the database and run migrations
go build
./goATED-tracker
```

## API Endpoints

### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/users` | Create a new user |
| DELETE | `/api/users` | Wipe all users |
| POST | `/api/login` | User login |
| GET | `/api/users/teams` | Get team membership (as a user) |

### Admin - Teams

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/admin/teams` | Create a team |
| GET | `/api/admin/teams` | Get all teams |
| DELETE | `/api/admin/teams` | Wipe all teams |
| GET | `/api/admin/teams/{TeamName}` | Get a specific team |

### Admin - Members

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/admin/teams/members` | Add a member to a team |
| GET | `/api/admin/teams/members` | Get all team members |
| DELETE | `/api/admin/teams/members` | Wipe all team members |

### Admin - Pentests

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/admin/teams/pentests` | Add a pentest |
| GET | `/api/admin/teams/pentests` | Get all pentests |
| DELETE | `/api/admin/teams/pentests` | Wipe all pentests |

### Admin - Findings

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/admin/findings` | Add a finding |
| GET | `/api/admin/findings` | Get all findings |
| DELETE | `/api/admin/findings` | Wipe all findings |

### User - Team Pentests & Findings

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/{TeamName}/pentests` | Get pentests for a team |
| POST | `/api/v1/{TeamName}/pentests` | Add a pentest to a team |
| POST | `/api/v1/{TeamName}/pentests/findings` | Add a finding to a pentest |
| GET | `/api/v1/{TeamName}/{PentestTitle}/findings` | Get findings for a pentest |
| POST | `/api/v1/{TeamName}/{PentestTitle}/findings` | Close findings for a pentest |


