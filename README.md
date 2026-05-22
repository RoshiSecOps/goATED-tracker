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

## Usage


