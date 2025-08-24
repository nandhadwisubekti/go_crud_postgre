# Employee Management API

Sebuah REST API untuk manajemen data pegawai yang dibangun dengan Go, Gin framework, PostgreSQL, dan JWT authentication.

## 🚀 Fitur

- **Authentication & Authorization**: JWT-based authentication dengan secure token
- **CRUD Operations**: Create, Read, Update, Delete data pegawai
- **Advanced Filtering**: Filter berdasarkan department, position, status aktif
- **Search Functionality**: Pencarian berdasarkan nama, email, atau NIP
- **Pagination**: Dukungan pagination untuk performa optimal
- **Data Validation**: Validasi input yang komprehensif
- **Security**: Password hashing dengan bcrypt, CORS middleware
- **Soft Delete**: Penghapusan data dengan soft delete (is_active flag)

## 🛠️ Tech Stack

- **Backend**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Environment Management**: godotenv

## 📋 Prerequisites

Sebelum menjalankan aplikasi, pastikan Anda telah menginstall:

- Go 1.21 atau lebih baru
- PostgreSQL 12 atau lebih baru
- Git

## 🔧 Installation

### 1. Clone Repository

```bash
git clone <repository-url>
cd go_crud_postgre
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Setup Database

Buat database PostgreSQL:

```sql
CREATE DATABASE employee_db;
```

### 4. Environment Configuration

Salin file `.env` dan sesuaikan konfigurasi:

```bash
cp .env.example .env
```

Edit file `.env`:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=employee_db
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY=24h

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Environment
ENV=development
```

### 5. Run Application

```bash
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:8080`

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication

Semua endpoint employee memerlukan JWT token di header:
```
Authorization: Bearer <your-jwt-token>
```

### Default Admin User

Aplikasi akan membuat user admin default:
- **Username**: `admin`
- **Password**: `admin123`
- **Email**: `admin@company.com`

### Endpoints

#### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register user baru |
| POST | `/auth/login` | Login user |
| GET | `/auth/profile` | Get profile user (protected) |

#### Employees

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/employees/` | Buat pegawai baru |
| GET | `/employees/` | Get semua pegawai dengan filter |
| GET | `/employees/:id` | Get pegawai berdasarkan ID |
| PUT | `/employees/:id` | Update data pegawai |
| DELETE | `/employees/:id` | Hapus pegawai (soft delete) |

#### Query Parameters untuk GET /employees/

- `limit`: Jumlah data per halaman (default: 10, max: 100)
- `offset`: Offset untuk pagination (default: 0)
- `department`: Filter berdasarkan department
- `position`: Filter berdasarkan position
- `is_active`: Filter berdasarkan status aktif (true/false)
- `search`: Pencarian berdasarkan nama, email, atau NIP

### Request/Response Examples

#### Register User

**Request:**
```json
POST /api/v1/auth/register
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 2,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

#### Login

**Request:**
```json
POST /api/v1/auth/login
{
  "username": "admin",
  "password": "admin123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@company.com"
    },
    "expires_at": "2024-08-25T08:22:00Z"
  }
}
```

#### Create Employee

**Request:**
```json
POST /api/v1/employees/
Authorization: Bearer <token>
{
  "nip": "EMP001",
  "name": "John Doe",
  "email": "john.doe@company.com",
  "phone": "+62812345678",
  "position": "Software Engineer",
  "department": "IT",
  "salary": 15000000,
  "hire_date": "2024-01-15"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Employee created successfully",
  "data": {
    "id": 1,
    "nip": "EMP001",
    "name": "John Doe",
    "email": "john.doe@company.com",
    "phone": "+62812345678",
    "position": "Software Engineer",
    "department": "IT",
    "salary": 15000000,
    "hire_date": "2024-01-15",
    "is_active": true,
    "created_at": "2024-08-24 08:22:00",
    "updated_at": "2024-08-24 08:22:00"
  }
}
```

#### Get Employees with Filters

**Request:**
```
GET /api/v1/employees/?department=IT&limit=5&offset=0
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "message": "Employees retrieved successfully",
  "data": {
    "employees": [
      {
        "id": 1,
        "nip": "EMP001",
        "name": "John Doe",
        "email": "john.doe@company.com",
        "phone": "+62812345678",
        "position": "Software Engineer",
        "department": "IT",
        "salary": 15000000,
        "hire_date": "2024-01-15",
        "is_active": true,
        "created_at": "2024-08-24 08:22:00",
        "updated_at": "2024-08-24 08:22:00"
      }
    ],
    "total": 1,
    "limit": 5,
    "offset": 0
  }
}
```

## 🧪 Testing

### Menggunakan HTTP Files

Gunakan file `api_test.http` dengan VS Code REST Client extension:

1. Install extension "REST Client" di VS Code
2. Buka file `api_test.http`
3. Klik "Send Request" pada request yang ingin ditest
4. Update token setelah login

### Menggunakan Postman

1. Import file `Employee_API.postman_collection.json` ke Postman
2. Set environment variable `baseUrl` ke `http://localhost:8080`
3. Login terlebih dahulu untuk mendapatkan token
4. Token akan otomatis tersimpan di environment variable

### Menggunakan cURL

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# Create Employee (ganti <token> dengan token dari login)
curl -X POST http://localhost:8080/api/v1/employees/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "nip": "EMP001",
    "name": "John Doe",
    "email": "john.doe@company.com",
    "position": "Software Engineer",
    "department": "IT",
    "hire_date": "2024-01-15"
  }'

# Get All Employees
curl -X GET "http://localhost:8080/api/v1/employees/?limit=10" \
  -H "Authorization: Bearer <token>"
```

## 📁 Project Structure

```
go_crud_postgre/
├── cmd/
│   └── api/
│       └── main.go              # Entry point aplikasi
├── config/
│   └── config.go               # Konfigurasi aplikasi
├── database/
│   └── database.go             # Database connection dan setup
├── handlers/
│   ├── auth.go                 # Authentication handlers
│   └── employee.go             # Employee CRUD handlers
├── middleware/
│   ├── auth.go                 # JWT authentication middleware
│   └── cors.go                 # CORS middleware
├── models/
│   ├── employee.go             # Employee models
│   ├── user.go                 # User models
│   └── response.go             # API response models
├── utils/
│   ├── jwt.go                  # JWT utilities
│   └── password.go             # Password hashing utilities
├── .env                        # Environment variables
├── api_test.http              # HTTP test requests
├── Employee_API.postman_collection.json  # Postman collection
├── go.mod                      # Go module dependencies
├── go.sum                      # Go module checksums
└── README.md                   # Dokumentasi project
```

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: bcrypt untuk hashing password
- **CORS Protection**: Cross-Origin Resource Sharing middleware
- **Input Validation**: Validasi input menggunakan Gin binding
- **SQL Injection Protection**: Menggunakan prepared statements
- **Environment Variables**: Sensitive data disimpan di environment variables

## 🚀 Deployment

### Production Checklist

1. **Environment Variables**:
   - Ganti `JWT_SECRET` dengan secret key yang kuat
   - Set `ENV=production`
   - Konfigurasi database production

2. **Database**:
   - Setup PostgreSQL production
   - Jalankan migration
   - Setup backup strategy

3. **Security**:
   - Enable HTTPS
   - Setup firewall
   - Monitor logs

### Docker Deployment

Buat `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./main"]
```

## 🤝 Contributing

1. Fork repository
2. Buat feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push ke branch (`git push origin feature/AmazingFeature`)
5. Buat Pull Request

## 📝 License

Project ini menggunakan MIT License. Lihat file `LICENSE` untuk detail.

## 📞 Support

Jika Anda mengalami masalah atau memiliki pertanyaan:

1. Buka issue di GitHub repository
2. Periksa dokumentasi API
3. Pastikan environment sudah dikonfigurasi dengan benar

## 🔄 Changelog

### v1.0.0 (2024-08-24)
- Initial release
- JWT authentication
- Employee CRUD operations
- Advanced filtering dan search
- Pagination support
- API documentation
- Testing files (HTTP & Postman)

---

**Happy Coding! 🚀**