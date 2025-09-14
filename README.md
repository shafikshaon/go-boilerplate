# go-boilerplate

A production-ready RESTful API built with Go, Gin framework, PostgreSQL, Redis, and clean architecture principles.

## 🚀 Features

- **🏗️ Clean Architecture** - Layered architecture with dependency injection
- **🔥 Gin Framework** - High-performance HTTP web framework
- **🐘 PostgreSQL** - Robust relational database
- **🔴 Redis** - In-memory caching and session management
- **🔐 JWT Authentication** - Secure token-based authentication with Redis sessions
- **📊 GORM** - Feature-rich ORM for database operations
- **🐳 Docker Support** - Full containerization with Docker Compose
- **✅ Input Validation** - Comprehensive request validation
- **📝 Structured Logging** - Production-ready logging
- **⚡ Caching** - Redis-based user and session caching

## 📋 Prerequisites

- **Go 1.21+**
- **Docker & Docker Compose**
- **Make** (optional but recommended)

## 🛠️ Quick Start

### 🎯 **One-Command Setup**
```bash
./run.sh
```

### **Step-by-Step Setup**
```bash
# 1. Install dependencies
make setup

# 2. Start databases
make db-up

# 3. Run the application
make run
```

### **Development Workflow**
```bash
make dev  # Starts databases and runs app
```

## 🧪 Test the API

```bash
./test_api.sh
```

## 📚 API Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Service health status | No |
| POST | `/api/v1/auth/register` | Register new user | No |
| POST | `/api/v1/auth/login` | Login user | No |
| GET | `/api/v1/auth/me` | Get current user profile | Yes |
| POST | `/api/v1/auth/logout` | Logout user | Yes |
| POST | `/api/v1/users` | Create user | No |
| GET | `/api/v1/users` | Get all users (paginated) | No |
| GET | `/api/v1/users/:id` | Get user by ID (cached) | No |
| PUT | `/api/v1/users/:id` | Update user | Yes |
| DELETE | `/api/v1/users/:id` | Delete user | Yes |

## 💻 Example Requests

### Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get User (Cached)
```bash
curl -X GET http://localhost:8080/api/v1/users/1
```

## 🏗️ Project Structure

```
go-boilerplate/
├── cmd/                    # Application initialization
├── handlers/               # HTTP request handlers
├── routes/                 # Route definitions
├── models/                 # Data models & DTOs
│   ├── request/           # Request models
│   └── response/          # Response models
├── services/               # Business logic layer
│   └── interfaces/        # Service interfaces
├── repository/             # Data access layer
│   └── interfaces/        # Repository interfaces
├── database/               # Database connections
├── middleware/             # Custom middleware
├── utilities/              # Helper functions & Redis utils
├── config/                 # Configuration management
└── scripts/               # Utility scripts
```

## ⚙️ Configuration

### Environment Variables (.env)
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=app_db

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0

# Application
PORT=8080
JWT_SECRET=your-super-secret-jwt-key
```

## 🛠️ Development Commands

### **Database Management**
```bash
make db-up         # Start PostgreSQL & Redis
make db-down       # Stop databases
make psql          # Connect to PostgreSQL
make redis-cli     # Connect to Redis CLI
```

### **Development**
```bash
make setup         # Install dependencies
make run           # Run application
make dev           # Start databases and run app
make test          # Test API endpoints
make build         # Build binary
```

### **Docker Operations**
```bash
make docker-up     # Start all services
make docker-down   # Stop all services
```

## 🗄️ Database Features

### **PostgreSQL**
- Auto-migrations with GORM
- Connection pooling
- Transaction support

### **Redis Caching**
- User profile caching (30 min TTL)
- JWT session management (24 hour TTL)
- Cache invalidation on updates

## 🔒 Security Features

- **JWT Authentication** with Redis session storage
- **Password Hashing** using bcrypt
- **Input Validation** with comprehensive error handling
- **CORS** middleware configuration
- **SQL Injection** protection via GORM

## 📦 Deployment

### **Production with Docker**
```bash
docker-compose up -d --build
```

### **Manual Deployment**
```bash
make build
./bin/app
```

## 📋 Available Make Commands

| Command | Description |
|---------|-------------|
| `make setup` | Install dependencies |
| `make db-up` | Start databases |
| `make run` | Run the application |
| `make dev` | Start databases and run app |
| `make test` | Test API endpoints |
| `make clean` | Clean up containers |

## 🎉 You're All Set!

Your Go API with PostgreSQL and Redis is ready for production!

**Happy coding!** 🚀
