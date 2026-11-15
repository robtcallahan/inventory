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

var want = []models.Item{
	{
		Quantity:    1,
		Name:        "name1",
		Description: "desc1",
		Value:       100,
		Category:    "category1",
	},
	{
		Quantity:    2,
		Name:        "name2",
		Description: "desc2",
		Value:       200,
		Category:    "category2",
	},
	{
		Quantity:    1,
		Name:        "name3",
		Description: "desc3",
		Value:       300,
		Category:    "category2",
	},
	{
		Quantity:    1,
		Name:        "name4",
		Description: "desc4",
		Value:       400,
		Category:    "category3",
	},
}

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
	wantOne := models.Item{
		Name:        "name",
		Description: "desc",
		Value:       100,
	}

	err := repo.Save(wantOne)
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
	if got[0].Name != wantOne.Name || got[0].Description != wantOne.Description || got[0].Value != wantOne.
		Value {
		t.Errorf("expected: %v, got: %v", wantOne, got)
	}
}

func Test_SaveItems(t *testing.T) {
	conn := setupTestDB(t)
	repo := NewMySQLQueryRepo(conn)

	skipped, err := repo.SaveItems(want)
	if err != nil {
		t.Fatalf("could not save items: %v", err)
	}
	if len(skipped) != 0 {
		t.Fatalf("expected 0 skipped items, got %d", len(skipped))
	}
	got, err := repo.List()
	if err != nil {
		t.Fatalf("could not list items: %v", err)
	}
	if len(got) != 4 {
		t.Fatalf("expected 4 items, got %d", len(got))
	}
	for i, item := range got {
		if item.Quantity != want[i].Quantity {
			t.Errorf("expected want.Quantity: %d, got %d", want[i].Quantity, item.Quantity)
		}
		if item.Name != want[i].Name {
			t.Errorf("expected want.Name: %s, got %s", want[i].Name, item.Name)
		}
		if item.Description != want[i].Description {
			t.Errorf("expected want.Description: %s, got %s", want[i].Description, item.Description)
		}
		if item.Value != want[i].Value {
			t.Errorf("expected want.Value: %d, got %d", want[i].Value, item.Value)
		}
		if item.Category != want[i].Category {
			t.Errorf("expected want.Category: %s, got %s", want[i].Category, item.Category)
		}
	}
}

func Test_ListByCategory(t *testing.T) {
	conn := setupTestDB(t)
	repo := NewMySQLQueryRepo(conn)
	_, err := repo.SaveItems(want)
	if err != nil {
		t.Fatalf("could not save items: %v", err)
	}

	got, err := repo.ListByCategory("category2")
	if err != nil {
		t.Fatalf("could not get list by category: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
	for i, item := range got {
		if item.Quantity != want[i+1].Quantity {
			t.Errorf("expected want.Quantity: %d, got %d", want[i+1].Quantity, item.Quantity)
		}
		if item.Name != want[i+1].Name {
			t.Errorf("expected want.Name: %s, got %s", want[i+1].Name, item.Name)
		}
		if item.Description != want[i+1].Description {
			t.Errorf("expected want.Description: %s, got %s", want[i+1].Description, item.Description)
		}
		if item.Value != want[i+1].Value {
			t.Errorf("expected want.Value: %d, got %d", want[i+1].Value, item.Value)
		}
		if item.Category != want[i+1].Category {
			t.Errorf("expected want.Category: %s, got %s", want[i+1].Category, item.Category)
		}
	}
}

func Test_TotalByCategory(t *testing.T) {
	conn := setupTestDB(t)
	repo := NewMySQLQueryRepo(conn)
	_, err := repo.SaveItems(want)
	if err != nil {
		t.Fatalf("could not save items: %v", err)
	}

	got, err := repo.TotalByCategory("category2")
	if err != nil {
		t.Fatalf("could not total by category: %v", err)
	}
	if got != 700 {
		t.Fatalf("expected 700, got %d", got)
	}
}
