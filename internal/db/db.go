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

var DB *gorm.DB

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
		log.Fatalf("failed to connect to %s: %v", engine, err)
	}

	if err := DB.AutoMigrate(&User{}); err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}
}
