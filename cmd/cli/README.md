# goIAM CLI

A command-line interface for interacting with the [goIAM](https://github.com/javadmohebbi/goIAM) Identity & Access Management microservice.

---

## 🔧 Features

- Register new users
- Login and receive JWT
- Setup Two-Factor Authentication (2FA) with TOTP
- Verify TOTP code
- Disable 2FA
- Regenerate one-time backup codes
- Supports QR code output using `qrencode`

---

## 🚀 Installation

### Prerequisites

- Go 1.20+
- (Optional) `qrencode` for displaying QR codes in terminal:
  ```bash
  sudo apt install qrencode  # Debian/Ubuntu
  brew install qrencode      # macOS
  ```

### Clone and Run

```bash
git clone https://github.com/javadmohebbi/goIAM.git
cd goIAM/cmd/cli
go run main.go --help
```

Or build it:

```bash
go build -o goiam-cli main.go
./goiam-cli --help
```

---

## ⚙️ Global Flags

- `--api`: Base URL of goIAM API (default: `http://localhost:8080`)
- `--token`: JWT token for authenticated routes

---

## 🛠 Commands

### Register a new user

```bash
go run main.go register --username john --password secret123 --email john@example.com --first John --last Doe
```

### Login and get JWT

```bash
go run main.go login --username john --password secret123
```

### Setup 2FA (QR + Secret)

```bash
go run main.go --token=$JWT 2fa-setup
```

### Verify TOTP code

```bash
go run main.go --token=$JWT 2fa-verify --code=123456
```

### Disable 2FA

```bash
go run main.go --token=$JWT 2fa-disable --code=123456
```

### Regenerate backup codes

```bash
go run main.go --token=$JWT backup-codes
```

---

## 📦 Folder Structure

```bash
cli/
├── main.go             # Entry point
└── cmds/
    ├── registry.go     # Register commands
    ├── register.go     # `register` command
    ├── login.go        # `login` command
    ├── 2fa_setup.go    # `2fa-setup` command
    ├── 2fa_verify.go   # `2fa-verify` command
    ├── 2fa_disable.go  # `2fa-disable` command
    └── backup_codes.go # `backup-codes` command
```

---

## 📝 License

MIT © [Javad Mohebi](https://github.com/javadmohebbi)
