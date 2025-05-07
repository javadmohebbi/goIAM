# goIAM CLI

A command-line interface for interacting with the [goIAM](https://github.com/javadmohebbi/goIAM) Identity & Access Management microservice.

---

## Features

- Register new users with full details (email, phone, address, etc.)
- Login with optional TOTP-based 2FA support
- Interactive and non-interactive login flows (piped input or terminal prompts)
- Setup Two-Factor Authentication (2FA) with secret and QR generation
- Verify 2FA codes manually or during login flow
- Disable 2FA with TOTP confirmation
- Regenerate one-time backup codes for account recovery
- Uses standard Go + Cobra structure
- QR terminal output using `qrencode` (optional)

---

## Installation

### Prerequisites

- Go 1.20+
- (Optional) `qrencode` for QR code output:
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

## Global Flags

- `--api`: Base URL of goIAM API (default: `http://localhost:8080`)
- `--token`: JWT token for authenticated routes

---

## Commands

### Register a new user

Register a user using one of the following options:

#### ➤ With existing organization ID:
```bash
go run main.go register --username alice --email alice@example.com --organization-id 1
```

#### ➤ Creating a new organization automatically:
```bash
go run main.go register --username bob --email bob@example.com --organization-name "Bob Corp"
```

- `--organization-slug` is optional; if omitted, it will be derived from the name.
- If no `--organization-id` is provided, a new org will be created using name/slug or fallback to `goIAM-org-{uuid}`.

You will be prompted for password securely.

### Login with 2FA support

```bash
go run main.go login --username john
```

You will be prompted for password and 2FA code (if enabled).

Piped password input also works:
```bash
echo "mypassword" | go run main.go login -u john
```

### Setup 2FA (QR + Secret)

```bash
go run main.go --token=$JWT 2fa-setup
```

### Verify 2FA code

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

## Folder Structure

```
cmd/cli/
├── main.go             # CLI entry point
└── cmds/
    ├── registry.go     # Registers CLI subcommands
    ├── register.go     # User registration
    ├── login.go        # Login and 2FA prompt
    ├── 2fa_setup.go    # Setup TOTP 2FA
    ├── 2fa_verify.go   # Verify 2FA code
    ├── 2fa_disable.go  # Disable 2FA
    └── backup_codes.go # Regenerate backup codes
```

---

## License

[Javad Mohebi](https://github.com/javadmohebbi)