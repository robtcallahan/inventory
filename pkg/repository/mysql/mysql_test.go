package mysql

import (
	cfg "inventory/pkg/config"
	"inventory/pkg/repository"
	"testing"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

const ConfigFile = "../../../config/config.json"

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	var err error
	config, err := cfg.ReadConfig(ConfigFile)
	if err != nil {
		t.Fatalf("could not read config file: %s, err: %s", ConfigFile, err.Error())
	}

	conn, err := repository.Connect(&config.TestDatabase)
	if err != nil {
		t.Fatalf("connect test db: %v", err)
	}
	sqlDB, err := conn.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	if err := goose.SetDialect("mysql"); err != nil {
		t.Fatalf("set dialect: %v", err)
	}
	// tests run in the package dir, so the path is relative to pkg/repository/mysql
	if err := goose.Up(sqlDB, "../../../migrations"); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}
	if err := conn.Exec("TRUNCATE TABLE items").Error; err != nil {
		t.Fatalf("truncate items: %v", err)
	}
	return conn
}

func Test_Save(t *testing.T) {
	setupTestDB(t)
}
