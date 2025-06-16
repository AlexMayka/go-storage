# 🗄️ Go Storage

A high-performance, multi-tenant cloud storage system built with Go, featuring role-based access control (RBAC), JWT authentication, and seamless file management.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=for-the-badge&logo=postgresql&logoColor=white)](https://postgresql.org)
[![MinIO](https://img.shields.io/badge/MinIO-Object%20Storage-C72E29?style=for-the-badge&logo=minio&logoColor=white)](https://min.io)
[![JWT](https://img.shields.io/badge/JWT-Auth-000000?style=for-the-badge&logo=JSON%20web%20tokens)](https://jwt.io)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

## 🌟 Features

- **🔐 Multi-tenant Architecture** - Complete company isolation
- **👥 Role-Based Access Control** - Granular permissions system
- **🗂️ Unified File System** - Files and folders in single table with materialized paths
- **📡 Streaming Uploads** - Handle large files efficiently without memory issues
- **🔑 JWT Authentication** - Secure token-based authentication
- **⚡ High Performance** - Optimized queries with proper indexing
- **🧪 Comprehensive Testing** - Full test coverage with mocks
- **📊 Clean Architecture** - Well-structured codebase following best practices

## 🏗️ Architecture

```mermaid
graph TB
    subgraph "Presentation Layer"
        A[REST API Handlers]
        B[Middleware]
    end
    
    subgraph "Business Logic Layer"
        C[Use Cases]
    end
    
    subgraph "Data Layer"
        D[PostgreSQL Repository]
        E[MinIO Repository]
    end
    
    subgraph "External"
        F[(PostgreSQL)]
        G[(MinIO)]
    end
    
    A --> B
    B --> C
    C --> D
    C --> E
    D --> F
    E --> G
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+**
- **PostgreSQL 15+**
- **MinIO** (or S3-compatible storage)

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/AlexMayka/go-storage.git
cd go-storage
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. **Run database migrations**
```bash
goose -dir migrations postgres "your-postgres-connection-string" up
```

5. **Start the server**
```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

## 🎯 API Endpoints

### 🔐 Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/api/v1/auth/login` | User login | ❌ |
| `POST` | `/api/v1/auth/refresh-token` | Refresh JWT token | ❌ |
| `POST` | `/api/v1/users/register` | Register new user | ❌ |

### 🏢 Companies

| Method | Endpoint | Description | Permission Required |
|--------|----------|-------------|-------------------|
| `GET` | `/api/v1/companies/` | Get all companies | `company:read:all` |
| `POST` | `/api/v1/companies/` | Create company | `company:create` |
| `GET` | `/api/v1/companies/{id}` | Get company by ID | `company:read:all` |
| `PUT` | `/api/v1/companies/{id}` | Update company | `company:update:all` |
| `DELETE` | `/api/v1/companies/{id}` | Delete company | `company:delete` |
| `GET` | `/api/v1/companies/me` | Get my company | `company:read:own` |
| `PUT` | `/api/v1/companies/me` | Update my company | `company:update:own` |

### 👥 Users

| Method | Endpoint | Description | Permission Required |
|--------|----------|-------------|-------------------|
| `GET` | `/api/v1/users/me` | Get my profile | `user:read` |
| `PUT` | `/api/v1/users/me` | Update my profile | `user:update` |
| `PUT` | `/api/v1/users/me/password` | Change my password | `user:update` |
| `GET` | `/api/v1/users/company` | Get company users | `user:read_company` |
| `GET` | `/api/v1/users/` | Get all users | `user:read` |
| `POST` | `/api/v1/users/` | Create user | `user:create` |
| `GET` | `/api/v1/users/{id}` | Get user by ID | `user:read` |
| `PUT` | `/api/v1/users/{id}` | Update user | `user:update` |
| `DELETE` | `/api/v1/users/{id}` | Deactivate user | `user:delete` |

### 🗂️ Files (Coming Soon)

| Method | Endpoint | Description | Permission Required |
|--------|----------|-------------|-------------------|
| `GET` | `/api/v1/files/` | List files/folders | `file:read` |
| `POST` | `/api/v1/files/` | Upload file/Create folder | `file:create` |
| `GET` | `/api/v1/files/{id}` | Get file info | `file:read` |
| `GET` | `/api/v1/files/{id}/download` | Download file | `file:read` |
| `PUT` | `/api/v1/files/{id}` | Update file/folder | `file:update` |
| `DELETE` | `/api/v1/files/{id}` | Delete file/folder | `file:delete` |

## 🔑 Authentication Example

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "admin@company.com",
    "password": "SecurePass123!"
  }'
```

**Response:**
```json
{
  "status": "success",
  "answer": {
    "id": "uuid-here",
    "username": "admin",
    "email": "admin@company.com",
    "company_id": "company-uuid",
    "role_id": "role-uuid",
    "auth": {
      "token": "eyJhbGciOiJIUzI1NiIs...",
      "expire_time": "2024-01-01T12:00:00Z"
    }
  }
}
```

### Using JWT Token
```bash
curl -X GET http://localhost:8080/api/v1/companies/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

## 👥 User Roles & Permissions

### 🔱 Super Admin
- **Full system access**
- Can manage all companies, users, and files
- System-wide permissions

### 🏢 Company Admin  
- **Company-level management**
- Can manage users within their company
- Can update company information
- Full access to company files

### 👤 User
- **Personal access only**
- Can view company information (read-only)
- Can manage their own profile
- Can access assigned files/folders

## 🗄️ Database Schema

### Core Tables

| Table | Description |
|-------|-------------|
| `companies` | Company information and settings |
| `users` | User accounts with role assignments |
| `roles` | System roles (super_admin, company_admin, user) |
| `permissions` | Granular permission definitions |
| `role_permissions` | Role-permission mappings |
| `files` | Unified files and folders with materialized paths |

### Key Features

- **Materialized Paths** - Efficient folder navigation using `full_path` column
- **Company Isolation** - All data filtered by `company_id`
- **Soft Deletes** - Data safety with `is_active` flags
- **Audit Trail** - Created/updated timestamps and user tracking

## 🔧 Configuration

### Environment Variables

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=go_storage
DB_SSLMODE=disable

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
MINIO_BUCKET=files

# JWT
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRES_IN=24h

# Server
SERVER_PORT=8080
SERVER_HOST=localhost
```

## 🧪 Testing

Run the test suite:
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/delivery/http/handlers/hdAuth/...
```

## 📈 Performance Features

- **Connection Pooling** - Optimized database connections
- **Permission Caching** - 10-minute cache with automatic cleanup
- **Streaming Uploads** - Memory-efficient file handling
- **Indexed Queries** - Proper database indexing strategy
- **Materialized Paths** - Fast folder tree operations

## 🚀 Deployment

### Docker (Coming Soon)
```bash
docker-compose up -d
```

### Manual Deployment
1. Build the application: `go build -o go-storage cmd/api/main.go`
2. Run migrations: `goose -dir migrations postgres $DATABASE_URL up`
3. Start the server: `./go-storage`

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://gorm.io/) - ORM library for Golang
- [MinIO](https://min.io/) - High performance object storage
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation
- [Testify](https://github.com/stretchr/testify) - Testing toolkit

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/AlexMayka">Aleksey Mayka</a>
</p>