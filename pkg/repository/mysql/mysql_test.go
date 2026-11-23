package mysql

import (
	cfg "inventory/pkg/config"
	"inventory/pkg/models"
	"inventory/pkg/repository"
	"reflect"
	"testing"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

const ConfigFile = "../../../config/config.json"

var want = []models.Item{
	{
		Quantity:    1,
		Name:        "name3",
		Description: "desc3",
		Value:       300,
		Category:    "category2",
	},
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
		t.Fatalf("migrate up on test db: %v", err)
	}
	return conn
}

func TeardownTestDB(t *testing.T, conn *gorm.DB) {
	sqlDB, err := conn.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	if err := goose.Down(sqlDB, "../../../migrations"); err != nil {
		t.Fatalf("migrate down on test db: %v", err)
	}
}

func Test_Save(t *testing.T) {
	conn := setupTestDB(t)
	repo := NewMySQLQueryRepo(conn)
	want := &models.Item{
		Name:        "name",
		Description: "desc",
		Value:       100,
	}

	item, err := repo.Save(want)
	if err != nil {
		t.Fatalf("could not save item: %v", err)
	}
	got, err := repo.Get(item.ID)
	if err != nil {
		t.Fatalf("could not get item with id: %d, %v", item.ID, err.Error())
	}
	if got.ID != 1 {
		t.Fatalf("expected id: 1, got %d", got.ID)
	}
	if got.Name != want.Name || got.Description != want.Description || got.Value != want.
		Value {
		t.Errorf("expected: %v, got: %v", want, got)
	}
	TeardownTestDB(t, conn)
}

func Test_Get(t *testing.T) {
	conn := setupTestDB(t)
	repo := NewMySQLQueryRepo(conn)
	want := &models.Item{
		Name:        "name",
		Description: "desc",
		Value:       300,
	}
	item, err := repo.Save(want)
	if err != nil {
		t.Fatalf("could not save item: %v", err)
	}
	got, err := repo.Get(item.ID)
	if err != nil {
		t.Fatalf("could not get item with id: %d, %v", item.ID, err.Error())
	}
	if got.ID != 1 {
		t.Fatalf("expected id: 1, got %d", got.ID)
	}
	if got.Name != want.Name || got.Description != want.Description {
		t.Errorf("expected: %v, got: %v", want, got)
	}
	TeardownTestDB(t, conn)
}

func Test_Delete(t *testing.T) {
	conn := setupTestDB(t)
	repo := NewMySQLQueryRepo(conn)

	item := &models.Item{
		Name:        "name",
		Description: "desc",
		Value:       100,
	}
	_, err := repo.Save(item)
	if err != nil {
		t.Fatalf("could not save item: %v", err)
	}
	err = repo.Delete(*item)
	if err != nil {
		t.Errorf("could not delete item: %v", err)
	}
	TeardownTestDB(t, conn)
}

func Test_SaveAndListItems(t *testing.T) {
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
	wantI := []int{1, 2, 0, 3}
	for i, item := range got {
		if item.Quantity != want[wantI[i]].Quantity {
			t.Errorf("expected want.Quantity: %d, got %d", want[wantI[i]].Quantity, item.Quantity)
		}
		if item.Name != want[wantI[i]].Name {
			t.Errorf("expected want.Name: %s, got %s", want[wantI[i]].Name, item.Name)
		}
		if item.Description != want[wantI[i]].Description {
			t.Errorf("expected want.Description: %s, got %s", want[wantI[i]].Description, item.Description)
		}
		if item.Value != want[wantI[i]].Value {
			t.Errorf("expected want.Value: %d, got %d", want[wantI[i]].Value, item.Value)
		}
		if item.Category != want[wantI[i]].Category {
			t.Errorf("expected want.Category: %s, got %s", want[wantI[i]].Category, item.Category)
		}
	}
	TeardownTestDB(t, conn)
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
	wantI := []int{2, 0}
	for i, item := range got {
		if item.Quantity != want[wantI[i]].Quantity {
			t.Errorf("expected want.Quantity: %d, got %d", want[wantI[i]].Quantity, item.Quantity)
		}
		if item.Name != want[wantI[i]].Name {
			t.Errorf("expected want.Name: %s, got %s", want[wantI[i]].Name, item.Name)
		}
		if item.Description != want[wantI[i]].Description {
			t.Errorf("expected want.Description: %s, got %s", want[wantI[i]].Description, item.Description)
		}
		if item.Value != want[wantI[i]].Value {
			t.Errorf("expected want.Value: %d, got %d", want[wantI[i]].Value, item.Value)
		}
		if item.Category != want[wantI[i]].Category {
			t.Errorf("expected want.Category: %s, got %s", want[wantI[i]].Category, item.Category)
		}
	}
	TeardownTestDB(t, conn)
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
	TeardownTestDB(t, conn)
}

func Test_SortedListForOutput(t *testing.T) {
	conn := setupTestDB(t)
	repo := NewMySQLQueryRepo(conn)
	_, err := repo.SaveItems(want)
	if err != nil {
		t.Fatalf("could not save items: %v", err)
	}
	got, err := repo.SortedListForOutput()
	if err != nil {
		t.Fatalf("could not get sorted list for output: %v", err)
	}
	if len(got) != 4 {
		t.Fatalf("expected 4 items, got %d", len(got))
	}
	// name2, name4, name3, name1. Desc2 before desc4
	wantI := []string{"name2", "name4", "name3", "name1"}
	for i, item := range got {
		if item.Name != wantI[i] {
			t.Errorf("expected want.Name: %s, got %s", wantI[i], item.Name)
		}
	}
	TeardownTestDB(t, conn)
}

func Test_SkippingDuplicates(t *testing.T) {
	conn := setupTestDB(t)
	repo := NewMySQLQueryRepo(conn)
	skippedFirstSave, err := repo.SaveItems(want)
	if err != nil {
		t.Fatalf("could not save items: %v", err)
	}
	if len(skippedFirstSave) != 0 {
		t.Fatalf("expected 0 items on first save, got %d", len(skippedFirstSave))
	}
	skippedSecondSave, err := repo.SaveItems(want)
	if err != nil {
		t.Fatalf("could not save items: %v", err)
	}
	if len(skippedSecondSave) != 4 {
		t.Fatalf("expected 4 items on second save, got %d", len(skippedSecondSave))
	}
	TeardownTestDB(t, conn)
}

func TestMysqlQueryRepo_SavePhoto(t *testing.T) {
	conn := setupTestDB(t)
	repo := NewMySQLQueryRepo(conn)
	wantItem := &models.Item{
		Name:        "name",
		Description: "desc",
		Value:       100,
	}

	item, err := repo.Save(wantItem)
	if err != nil {
		t.Fatalf("could not save item: %v", err)
	}

	wantPhoto := models.ItemPhoto{
		ItemID: item.ID,
	}
	photo, err := repo.Get(item.ID)
	if err != nil {
		t.Fatalf("could not save photo: %v", err)
	}
	if !reflect.DeepEqual(photo, wantPhoto) {
		t.Errorf("expected photo %v, got %v", wantPhoto, photo)
	}
	TeardownTestDB(t, conn)
}
