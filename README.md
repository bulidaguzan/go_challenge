# go_challenge

## Project Structure
.
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── cmd
│   └── main.go
├── internal
│   ├── api
│   │   ├── handlers
│   │   │   ├── balance.go
│   │   │   └── migration.go
│   │   ├── middleware
│   │   │   └── middleware.go
│   │   └── router.go
│   ├── config
│   │   └── config.go
│   ├── models
│   │   └── transaction.go
│   ├── repository
│   │   └── postgres
│   │       └── transaction.go
│   └── service
│       ├── balance.go
│       └── migration.go
└── README.md

## API Endpoints:

- POST /migrate for CSV file processing
- GET /users/{user_id}/balance for balance queries
Support for date range filtering

## To run the project:

Just use docker-compose up --build

The API will be available at http://localhost:8080