package mysql

import (
	cfg "inventory/pkg/config"
	"inventory/pkg/models"
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
	conn := setupTestDB(t)
	repo := NewMySQLQueryRepo(conn)
	want := models.Item{
		Name:        "name",
		Description: "desc",
		Value:       100,
	}
	err := repo.Save(want)
	if err != nil {
		t.Fatalf("could not save item: %v", err)
	}
	got, err := repo.List()
	if err != nil {
		t.Fatalf("could not list items: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}
	if got[0].Name != want.Name || got[0].Description != want.Description || got[0].Value != want.
		Value {
		t.Errorf("expected: %v, got: %v", want, got)
	}
}

func Test_SaveItems(t *testing.T) {
	conn := setupTestDB(t)
	repo := NewMySQLQueryRepo(conn)
	want := []models.Item{
		{
			Name:        "name1",
			Description: "desc1",
			Value:       100,
		},
		{
			Name:        "name2",
			Description: "desc2",
			Value:       200,
		},
	}

	got, err := repo.SaveItems(want)
	if err != nil {
		t.Fatalf("could not save items: %v", err)
	}
	got, err = repo.List()
	if err != nil {
		t.Fatalf("could not list items: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
	if got[0].Name != want[0].Name || got[0].Description != want[0].Description || got[0].Value != want[0].
		Value || got[1].Name != want[1].Name || got[1].Description != want[1].Description || got[1].Value != want[1].
		Value {
		t.Errorf("expected: %v, got: %v", want, got)
	}
}

func Test_ListByCategory(t *testing.T) {
	conn := setupTestDB(t)
	repo := NewMySQLQueryRepo(conn)
	want := []models.Item{
		{
			Name:        "name1",
			Description: "desc1",
			Value:       100,
			Category:    "category1",
		},
		{
			Name:        "name2",
			Description: "desc2",
			Value:       200,
			Category:    "category2",
		},
		{
			Name:        "name3",
			Description: "desc3",
			Value:       300,
			Category:    "category2",
		},
		{
			Name:        "name4",
			Description: "desc4",
			Value:       400,
			Category:    "category3",
		},
	}
	_, err := repo.SaveItems(want)
	if err != nil {
		t.Fatalf("could not save items: %v", err)
	}
	got, err := repo.ListByCategory("category2")
	if err != nil {
		t.Fatalf("could not list items: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
	if got[0].Name != want[1].Name || got[0].Description != want[1].Description || got[0].Value != want[1].
		Value || got[1].Name != want[2].Name || got[1].Description != want[2].Description || got[1].Value != want[2].
		Value {
		t.Errorf("expected: %v, got: %v", want, got)
	}
}
