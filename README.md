# FinTech Solutions Backend

This project implements a backend service for FinTech Solutions Inc., providing APIs for transaction migration and balance information.

## Features

- CSV transaction data migration
- User balance queries with optional date range filtering
- Docker containerization
- PostgreSQL database integration
- RESTful API endpoints
- Error handling and validation

## Prerequisites

- Docker
- Docker Compose

## Setup and Running

1. Clone the repository
2. Navigate to the project directory
3. Run:
```bash
docker-compose up --build
```

The service will be available at `http://localhost:8080`

## API Documentation

### Migration Service

#### POST /migrate
Uploads and processes a CSV file containing transaction records.

Request:
- Method: POST
- Content-Type: multipart/form-data
- Body: CSV file with columns: id, user_id, amount, datetime

Response:
- 200: Migration successful
- 400: Invalid file format or data
- 500: Server error

### Balance Service

#### GET /users/{user_id}/balance
Returns balance information for a specific user.

Parameters:
- user_id: User ID (path parameter)
- from: Start datetime (optional query parameter, RFC3339 format)
- to: End datetime (optional query parameter, RFC3339 format)

Response:
```json
{
    "balance": 25.21,
    "total_debits": 10,
    "total_credits": 15
}
```

Status Codes:
- 200: Success
- 400: Invalid parameters
- 500: Server error
