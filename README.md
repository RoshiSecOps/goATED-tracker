# :shield: GoATED Pentest Tracker 

## Project Idea
A robust API with multiple endpoints that allows tracking of Penetration tests across multiple SWE Teams, along with all vulnerability findings associated with the conducted test. For authentication we use JWTs and there is authorization in place across team <-> member , pentest <-> team and finding <-> pentest. Users are isolated to create records for Penetration tests and findings only to those that are assigned to their respective team.

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
| POST | `/api/login` | User login |
| GET | `/api/users/teams` | Get team membership (as a user) |

### Admin - Teams

| Method | Endpoint | Description |
|--------|----------|-------------|
| DELETE | `/api/admin/users` | Wipe all users |
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

### Example Flow
```bash
#1. Use /api/login to get the JWT
curl http://localhost:8080/api/login -d '{"username":"x","password":"y"}' # save the JWT

#2. You need to be a member of a team, this is done via the admin endpoints using the ADMIN_SECRET as a token bearer under Authorization Header.
curl http://localhost:8080/api/users/teams -H "Authorization: Bearer <yourJWT>" # gives you the teams you are a member of.

#3. Create a Pentest Record for a team.
curl -X POST http://localhost:8080/api/v1/{TeamName}/pentests -H "Authorization: Bearer <yourJWT>" -d '{"title":"Pentest-Title","team_id":"ID-Of-Team-To-Map-Pentest"}'

#4. Add a Finding to the Pentest Record.
curl -X POST http://localhost:8080/api/v1/{TeamName}/pentests/findings -H "Authorization: Bearer <yourJWT>" -d '{"title": "SQL Injection in Search Bar", "status": "open", "severity": "critical", "severity_score": 10, "file": "search.go", "at_line": 124, "description": "User input is being concatenated directly into the query string.", "pentest_id": "your-pentest-id"}'

#5. Get All findings for a given pentest.
curl http://localhost:8080/api/v1/{TeamName}/{PentestTitle}/findings -H "Authorization: Bearer <yourJWT>"

#6. Close a finding.
curl -X POST http://localhost:8080/api/v1/{TeamName}/{PentestTitle}/findings -H "Authorization: Bearer <yourJWT>" -d '{"finding_id":"Your-Finding-ID"}'
```

### Add Members to teams via Admin endpoint
```bash
#1. Add User to a team:
curl -X POST http://localhost:8080/api/admin/teams/members -H "Authorization: Bearer <adminToken>" -d '{"teamname":"name-of-team","username":"user-to-add"}'

#2. Create a new team:
curl -X POST http://localhost:8080/api/admin/teams -H "Authorization: Bearer <adminToken>" -d '{"teamname":"name-of-team"}'
```

### Running Locally via Docker Compose

```bash
git clone https://github.com/RoshiSecOps/goATED-tracker.git
cd goATED-tracker.git
# rename the .docker.env-example file to .docker.env and generate your own JWT_SECRET and ADMIN_KEY, for example via 'openssl rand -base64 64'
docker compose up --build #Starts the API and databes
```
Once done, you should be able to use the provided commands with your ADMIN key or via JWT once you login.