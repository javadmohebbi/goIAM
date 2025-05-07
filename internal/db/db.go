// Package db provides database initialization and access using GORM,
// supporting multiple SQL database engines.
package db

import (
	"log"

	"gorm.io/gorm"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
)

// DB is a global variable representing the initialized GORM database instance.
// It is set by the Init function.
var DB *gorm.DB

// Init sets up the global DB connection using the provided database engine and DSN.
//
// Supported engines include: "sqlite", "postgres", "mysql", "sqlserver", "clickhouse".
// It uses the GORM library to establish the connection and automatically migrates
// the defined models (Organization, User, Group, Role, Policy, BackupCode).
//
// Parameters:
//   - engine: name of the database engine (e.g., "sqlite", "postgres")
//   - dsn: data source name (connection string)
//
// Logs fatal errors if the engine is unsupported, the connection fails,
// or if model migration fails.
func Init(engine, dsn string) {
	var dialector gorm.Dialector

	switch engine {
	case "sqlite":
		dialector = sqlite.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	case "mysql":
		dialector = mysql.Open(dsn)
	case "sqlserver":
		dialector = sqlserver.Open(dsn)
	case "clickhouse":
		dialector = clickhouse.Open(dsn)
	default:
		log.Fatalf("unsupported database engine: %s", engine)
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to %s: %v | %s", engine, err, dsn)
	}

	// Automatically migrate database schemas for the core models
	if err := DB.AutoMigrate(
		&Organization{},
		&User{},
		&Group{},
		&Role{},
		&Policy{},
		&BackupCode{},
	); err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}

	// Create default organization if none exist
	var count int64
	if err := DB.Model(&Organization{}).Count(&count).Error; err != nil {
		log.Fatalf("failed to count organizations: %v", err)
	}
	if count == 0 {
		if err := DB.Create(&Organization{Name: "goIAM"}).Error; err != nil {
			log.Fatalf("failed to create default organization: %v", err)
		}
		log.Println("default organization 'goIAM' created")
	}
}
