# goIAM

**goIAM** is a modern, pluggable Identity and Access Management (IAM) microservice written in Go. It supports multiple authentication providers (local DB, LDAP, Firebase, Auth0), role-based access control (RBAC), SSO, and 2FA — all configurable via YAML and exposed through a REST API and CLI.

---

## 🚀 Features

- ✅ Modular authentication (Local, LDAP/Active Directory, Firebase, Auth0, ...)
- 🔐 Role-based access control (RBAC)
- 🔁 SSO-ready structure
- 🔢 TOTP-based 2FA support (Google Authenticator, etc.)
- ⚙️ Configurable via `config.yaml`
- 🌐 REST API (Fiber v3) + CLI with flags
- 🔐 JWT-based authentication
- 🧪 Easy to run with `go run .`

---

## 📁 Project Structure

```bash
goIAM/
├── cmd/server/           # Main API entry point (real logic)
├── internal/
│   ├── api/              # Fiber routes & handlers
│   ├── auth/             # Local, LDAP, Firebase, Auth0 backends
│   ├── config/           # YAML config loader
│   ├── db/               # Models and database logic
│   ├── policy/           # Policy enforcement (RBAC)
│   ├── sso/              # SSO logic (future)
│   └── totp/             # 2FA logic (TOTP)
├── pkg/                  # Public libraries/utilities (optional)
├── main.go               # Thin wrapper for go run .
├── config.yaml           # Example configuration file
├── go.mod
└── go.sum
```

---

## 📦 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/javadmohebbi/goIAM.git
cd goIAM
```

### 2. Run the API

```bash
go run .
```

### 3. CLI Flags

| Flag         | Description                          | Default       |
|--------------|--------------------------------------|---------------|
| `--config`   | Path to YAML config file             | `config.yaml` |
| `--port`     | HTTP server port                     | `8080`        |
| `--debug`    | Enable debug logging (true/false)    | `false`       |

Example:

```bash
go run . --port 9090 --debug --config=./config.yaml
```

---

## ⚙️ Sample `config.yaml`

```yaml
port: 8080
debug: true
database_dsn: "sqlite://./data.db"
auth_provider: "local" # options: local, ldap, auth0, firebase
```

---

## 📌 TODO

- [ ] Add OAuth2 and SAML integrations
- [ ] Admin panel for user/group/policy management
- [ ] JWT refresh tokens
- [ ] OpenID Connect support
- [ ] Audit logs and admin CLI

---

## 📄 License

MIT License © [Javad Mohebi](https://github.com/javadmohebbi)
