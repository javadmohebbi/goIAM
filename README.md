# goIAM

**goIAM** is a modern, modular Identity and Access Management (IAM) microservice in Go. It supports local authentication with username/password, 2FA (TOTP + backup codes), and pluggable future support for LDAP, Firebase, and Auth0.

---

## 🚀 Features

- ✅ Local authentication with password hashing
- 🔐 TOTP-based 2FA (Google Authenticator, Authy, etc.)
- 🔁 One-time backup codes
- 🔐 JWT-secured routes
- 🧩 Groups, Roles, Policies for future access control
- 🌐 Fiber v3 HTTP API + CLI compatibility
- ⚙️ Configurable with `config.yaml`

---

## 🏗️ Project Structure

```bash
goIAM/
├── cmd/server/           # CLI entry point (Main)
├── internal/
│   ├── api/              # Fiber routes and logic
│   ├── auth/             # Password hashing, TOTP, backup code
│   ├── config/           # YAML config loader
│   ├── db/               # GORM models and DB logic
│   ├── middleware/       # JWT + 2FA verification
├── main.go               # Thin wrapper for go run .
├── config.yaml           # Configuration file
```

---

## 📦 Getting Started

### 1. Clone and build

```bash
git clone https://github.com/javadmohebbi/goIAM.git
cd goIAM
go run .
```


### 2. Example `config.yaml`

```yaml
port: 8080
debug: true
jwt_secret: "your-secret"
database: "sqlite"
database_dsn: "./data/iam.db"
auth_provider: "local"
```

### 3. Environment Variables

You can override configuration values using environment variables:

| Variable              | Description                                | Example              |
|-----------------------|--------------------------------------------|----------------------|
| IAM_PORT              | Port for HTTP server                       | `8080`               |
| IAM_DEBUG             | Enable debug logging                       | `true`               |
| IAM_CONFIG_PATH       | Path to YAML config file                   | `./config.yaml`      |
| IAM_DATABASE          | Database engine                            | `sqlite`, `postgres` |
| IAM_DATABASE_DSN      | Database connection string (DSN)           | `./data/iam.db`      |
| IAM_AUTH_PROVIDER     | Comma-separated auth providers             | `local,ldap`         |

---

## 🔐 API Endpoints (Tested with curl)

### Register

```bash
curl -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d '{
  "username": "john",
  "password": "secret123",
  "email": "john@example.com",
  "phone_number": "1234567890",
  "first_name": "John",
  "middle_name": "Q",
  "last_name": "Public",
  "address": "123 Main St"
}'
```

### Login

```bash
curl -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d '{
  "username": "john",
  "password": "secret123"
}'
```

### 2FA Setup (TOTP)

```bash
curl -X POST http://localhost:8080/secure/auth/2fa/setup -H "Authorization: Bearer $TOKEN"
```

### 2FA Verify

```bash
curl -X POST http://localhost:8080/secure/auth/2fa/verify -H "Authorization: Bearer $TOKEN" -d '{"code": "123456"}'
```

### Backup Codes

```bash
curl -X POST http://localhost:8080/secure/auth/backup-codes/regenerate -H "Authorization: Bearer $TOKEN"
```

### 2FA Disable

```bash
curl -X POST http://localhost:8080/secure/auth/2fa/disable -H "Authorization: Bearer $TOKEN" -d '{"code": "123456"}'
```

---

## ✅ Coming Soon

- LDAP, Firebase, Auth0 login strategies
- Admin interface for managing users, policies, and roles
- OAuth2 / OpenID Connect support

---

## 📄 License

© [Javad Mohebi](https://github.com/javadmohebbi)
