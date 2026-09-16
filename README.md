# IHB Transport - Logistics Management System

## 📋 Project Overview

**IHB Transport** is a delivery and logistics management system built with Go. It provides a backend API for managing delivery requests, tracking shipments, and handling the complete lifecycle of logistics operations from request to delivery.

### What is this project about?

This project is a **transportation and delivery management platform** that enables:

- **Delivery Request Management**: Clients can submit delivery requests with pickup/dropoff addresses and item details
- **Pricing System**: Admins can set prices for delivery requests
- **Status Tracking**: Complete tracking of delivery status from request to completion
- **Email Notifications**: Automated email notifications for different delivery stages
- **Admin Dashboard**: Secure admin authentication and management capabilities
- **Audit Logging**: Complete audit trails for status changes and email communications

## 🎯 Key Features

1. **Delivery Request Lifecycle Management**
   - Create delivery requests with client information
   - Track multiple status stages: REQUESTED → PRICED → ACCEPTED → IN_PROGRESS → DELIVERED
   - Update delivery status with automatic logging

2. **Admin Authentication & Authorization**
   - Secure JWT-based authentication
   - Bcrypt password hashing
   - Protected admin endpoints

3. **Comprehensive Logging**
   - Status change logs (who changed what and when)
   - Email delivery logs (track all communications)
   - Complete audit trail for compliance

4. **Email Notification System**
   - Automated notifications at each delivery stage
   - Email delivery tracking
   - Failed email error logging

## 🛠️ Technology Stack

### Backend Framework & Language
- **Go 1.24.0** - Primary programming language
- **Gin** - High-performance HTTP web framework

### Database
- **PostgreSQL** - Relational database
- **GORM** - ORM for database operations
- **pgx** - PostgreSQL driver

### Security & Authentication
- **JWT (golang-jwt/jwt/v5)** - Token-based authentication
- **bcrypt (golang.org/x/crypto)** - Password hashing

### Utilities
- **UUID (google/uuid)** - Unique identifier generation
- **godotenv** - Environment variable management

## 📁 Project Structure

```
IHB-transport/
├── cmd/
│   ├── server/          # Main API server
│   ├── migrate/         # Database migration tool
│   └── seed/            # Database seeding tool (creates default admin)
├── internal/
│   ├── auth/            # Authentication & JWT logic
│   ├── database/        # Database connection setup
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Database models & schemas
│   ├── repository/      # Data access layer
│   └── services/        # Business logic layer
├── migrations/          # Database migration definitions
├── utils/               # Utility functions (validation, etc.)
├── .env                 # Environment configuration
├── go.mod               # Go module dependencies
└── go.sum               # Dependency checksums
```

## 🗄️ Database Schema

### Core Models

**Admin**
- Manages administrative users with secure authentication
- Fields: ID, Email, PasswordHash, Timestamps

**DeliveryRequest**
- Central entity for all delivery operations
- Fields: ID, ClientName, ClientEmail, PickupAddress, DropoffAddress, ItemDescription, Weight, Status, Price, Timestamps
- Status Flow: REQUESTED → PRICED → ACCEPTED → IN_PROGRESS → DELIVERED

**StatusLog**
- Audit trail for all status changes
- Fields: ID, DeliveryID, OldStatus, NewStatus, ChangedBy, ChangedAt

**EmailLog**
- Tracks all email communications
- Fields: ID, DeliveryID, RecipientEmail, EmailType, Status, SentAt, ErrorMessage

## 🚀 Getting Started

### Prerequisites

- Go 1.24.0 or higher
- PostgreSQL database
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/kwabsntim/IHB-transport.git
   cd IHB-transport
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   
   Create or update the `.env` file with your configuration:
   ```env
   POSTGRES_HOST=localhost
   POSTGRES_PORT=5432
   POSTGRES_USER=your_username
   POSTGRES_PASSWORD=your_password
   POSTGRES_DB=logisticsDB
   ADMIN_EMAIL=admin@example.com
   ADMIN_PASSWORD=your_secure_password
   JWT_SECRET_KEY=your_jwt_secret
   ```

4. **Run database migrations**
   ```bash
   go run cmd/migrate/main.go
   ```

5. **Seed the database (create default admin)**
   ```bash
   go run cmd/seed/main.go
   ```

6. **Start the server**
   ```bash
   go run cmd/server/main.go
   ```

   The server will start on `http://localhost:8080`

## 📡 API Endpoints

### Public Endpoints

- **POST /login** - Admin login
  - Request body: `{"email": "string", "password": "string"}`
  - Response: `{"token": "jwt_token", "message": "Login successful"}`

- **GET /ping** - Health check
  - Response: `{"message": "pong"}`

### Protected Endpoints (Require Authentication)

*Note: Additional endpoints for delivery management are being developed*

## 🔒 Security Features

- **Password Hashing**: All passwords are hashed using bcrypt
- **JWT Authentication**: Secure token-based authentication
- **Environment Variables**: Sensitive data stored in `.env` file (excluded from version control)
- **Input Validation**: Comprehensive validation for all user inputs
  - Email format validation
  - Required field validation
  - Positive number validation for weights and prices
  - Maximum length validation for descriptions

## 🏗️ Architecture

The application follows a **layered architecture** pattern:

1. **Handlers Layer** (`internal/handlers`)
   - Receives HTTP requests
   - Validates input
   - Returns HTTP responses

2. **Services Layer** (`internal/services`)
   - Contains business logic
   - Orchestrates operations across repositories
   - Handles validation and data transformation

3. **Repository Layer** (`internal/repository`)
   - Direct database operations
   - CRUD operations for all models
   - Query abstractions

4. **Models Layer** (`internal/models`)
   - Database schema definitions
   - Data structures and constants

## 📊 Business Workflow

1. **Client submits delivery request**
   - Status: REQUESTED
   - System creates status log

2. **Admin reviews and sets price**
   - Status: REQUESTED → PRICED
   - Email sent to client with price

3. **Client accepts delivery**
   - Status: PRICED → ACCEPTED
   - Delivery assigned to driver

4. **Driver begins delivery**
   - Status: ACCEPTED → IN_PROGRESS
   - Client notified that driver is on the way

5. **Delivery completed**
   - Status: IN_PROGRESS → DELIVERED
   - Completion notification sent

## 🧰 Available Commands

### Development Commands

```bash
# Run the API server
go run cmd/server/main.go

# Run database migrations
go run cmd/migrate/main.go

# Seed the database with default admin
go run cmd/seed/main.go
```

### Build Commands

```bash
# Build the server binary
go build -o bin/server cmd/server/main.go

# Build migration binary
go build -o bin/migrate cmd/migrate/main.go

# Build seed binary
go build -o bin/seed cmd/seed/main.go
```

## 🔧 Development

### Code Organization

- **Separation of Concerns**: Clear separation between layers
- **Interface-Based Design**: Repositories and services use interfaces for testability
- **Dependency Injection**: Services receive their dependencies via constructors
- **GORM Conventions**: Follows GORM best practices for ORM operations

### Validation

Input validation is centralized in the `utils` package and includes:
- Required field validation
- Email format validation
- Positive number validation
- Maximum length validation

## 📝 Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `POSTGRES_HOST` | PostgreSQL host | `localhost` |
| `POSTGRES_PORT` | PostgreSQL port | `5432` |
| `POSTGRES_USER` | Database user | `postgres` |
| `POSTGRES_PASSWORD` | Database password | `password` |
| `POSTGRES_DB` | Database name | `logisticsDB` |
| `ADMIN_EMAIL` | Default admin email | `admin@example.com` |
| `ADMIN_PASSWORD` | Default admin password | `SecurePass123!` |
| `JWT_SECRET_KEY` | JWT signing secret | `your-secret-key` |

## 🤝 Contributing

This project follows standard Go coding conventions and best practices:
- Use `go fmt` to format code
- Follow the existing project structure
- Write clear commit messages
- Test your changes before submitting

## 📄 License

*License information not specified*

## 👥 Project Information

- **Repository**: [kwabsntim/IHB-transport](https://github.com/kwabsntim/IHB-transport)
- **Owner**: kwabsntim

## 🔮 Future Enhancements

Potential areas for expansion:
- REST API endpoints for delivery CRUD operations
- Real-time tracking integration
- SMS notifications
- Mobile application support
- Analytics dashboard
- Payment processing integration
- Multi-language support

---

**Note**: This is an active project under development. Some features mentioned in the code may still be in implementation phase.
