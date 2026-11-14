package mysql

import (
	"inventory/pkg/config"
	"inventory/pkg/repository"
	"os"
	"testing"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	cfg := &config.DBConfig{
		Type: "mysql", Host: "localhost", Port: "3306",
		Name:     "inventory_test",
		Username: "inventory", Password: os.Getenv("INVENTORY_TEST_MYSQL_PASSWORD"),
	}
	conn, err := repository.Connect(cfg)
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
