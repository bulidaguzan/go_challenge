# FinTech Solutions Backend

## 📋 Description
A high-performance Go backend service that provides REST APIs for financial transaction management and balance reporting. Built with Gin framework and PostgreSQL, it offers robust transaction data migration from CSV files and detailed balance querying capabilities.

## ✨ Key Features
- CSV transaction data migration with upsert support
- User balance calculation with date range filtering
- Automatic database schema management
- Index optimization for query performance
- Environment-based configuration
- Docker-ready architecture with Docker Compose
- Decimal precision handling for financial data
- Integrated pgAdmin for database management

## 🛠️ Technical Stack
- **Language:** Go
- **Web Framework:** Gin
- **Database:** PostgreSQL 14 (Alpine)
- **Administration:** pgAdmin 4
- **Dependencies:**
  - `github.com/gin-gonic/gin`
  - `github.com/lib/pq`
  - Standard Go libraries

## 🚀 Installation and Setup

### Prerequisites
- Docker 20.10 or higher
- Docker Compose 2.0 or higher (supporting version 3.8 compose files)
- 4GB RAM recommended

### Docker Compose Setup
```bash
# Clone the repository
git clone https://github.com/bulidaguzan/go_challenge.git

# Navigate to project directory
cd go_challenge

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```
# Access Swagger UI:
- Navigate to `http://localhost:8080/swagger/index.html`
- Test endpoints directly from the interface


### Service Architecture
```yml
Services:
  - app: Financial transaction API (Port 8080)
  - db: PostgreSQL database (Port 5432)
  - pgadmin: Database administration (Port 5050)
```

### Environment Configuration
```yaml
# Application Service
DB_HOST: db
DB_PORT: 5432
DB_USER: postgres
DB_PASSWORD: postgres
DB_NAME: fintech

# PostgreSQL Service
POSTGRES_USER: postgres
POSTGRES_PASSWORD: postgres
POSTGRES_DB: fintech

# pgAdmin Service
PGADMIN_DEFAULT_EMAIL: admin@admin.com
PGADMIN_DEFAULT_PASSWORD: admin
```

## 📚 API Documentation

### Transaction Migration API

#### Upload Transactions
```http
POST /migrate
```

**Request Details:**
- Method: `POST`
- Content-Type: `multipart/form-data`
- Form Field: `file` (CSV file)

**CSV Format Example:**
```csv
id,user_id,amount,datetime
1,2,10.00,2024-07-01T15:00:11Z
2,1,-5.00,2024-07-01T15:00:12Z
3,4,-5.52,2024-06-01T14:59:59Z
```

**Important CSV Notes:**
- Supports transaction updates (upsert) based on transaction ID
- Datetime must be in RFC3339 format (UTC)
- Amount supports 2 decimal places
- Negative amounts represent debits
- Positive amounts represent credits

### Balance API

#### Get User Balance
```http
GET /users/{user_id}/balance
```

**Parameters:**
| Name | Location | Type | Required | Description |
|------|----------|------|----------|-------------|
| user_id | path | integer | Yes | User identifier |
| from | query | string | No | Start date (RFC3339) |
| to | query | string | No | End date (RFC3339) |

**Response Example:**
```json
{
    "balance": 100.00,
    "total_debits": 3,
    "total_credits": 5
}
```

## 🗄️ Database Management

### PostgreSQL Configuration
- Version: 14 (Alpine)
- Port: 5432
- Persistent storage: Docker volume `postgres_data`
- Health check: Every 5 seconds
- Network: Internal Docker network (fintech-network)

### pgAdmin Access
1. Access URL: `http://localhost:5050`
2. Login credentials:
   - Email: `admin@admin.com`
   - Password: `admin`

### Database Connection (via pgAdmin)
```yaml
Connection Details:
  - Host: db
  - Port: 5432
  - Database: fintech
  - Username: postgres
  - Password: postgres
```

## 🔍 Monitoring and Logs

### View Service Logs
```bash
# All services
docker-compose logs

# Specific services
docker-compose logs app
docker-compose logs db
docker-compose logs pgadmin

# Follow logs
docker-compose logs -f [service]
```



## 🔧 Troubleshooting

### Common Issues
1. **Service Startup Order**
   - The application waits for database health check before starting
   - Check logs if services aren't starting properly:
     ```bash
     docker-compose logs app
     ```

2. **Database Connection Issues**
   ```bash
   # Check if database is healthy
   docker-compose ps db
   
   # Check database logs
   docker-compose logs db
   ```

3. **Data Migration Issues**
   - Verify CSV format matches the example provided
   - Check file encoding (UTF-8 recommended)
   - Ensure datetime values are in UTC

### Container Management
```bash
# remove volume
docker volume rm go_challenge_postgres_data

# Restart specific service
docker-compose restart [service]

# View service status
docker-compose ps

# Remove containers and volumes
docker-compose down -v
```

## 🤝 Contributing
1. Fork the repository
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Implement tests for new features
4. Submit a Pull Request



