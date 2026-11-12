package database_provider

import (
	"fmt"
	"testing"
)

func Test_SaveAndList(t *testing.T) {
	provider := &DatabaseProvider{}
	item := InventoryItem{
		Quantity:     1,
		Name:         "Item 1",
		Manufacturer: "Manufacturer",
		Model:        "Model",
		Description:  "Description",
		SerialNumber: "Serial Number",
		Category:     "Category",
		PurchaseDate: "",
		Value:        100,
	}

	err := provider.Save(item)
	if err != nil {
		fmt.Errorf("got error %v, want nil", err)
	}
	items, err := provider.List()
	if err != nil {
		fmt.Errorf("got error %v, want nil", err)
	}
	if len(items) != 1 {
		t.Fatalf("got len(items) = %d, want 1", len(items))
	}
	if items[0] != item {
		t.Errorf("got item = %v, want %v", item, item)
	}
}

func Test_ListByCategory(t *testing.T) {
	provider := &DatabaseProvider{
		Repository: []InventoryItem{
			{Name: "Taylor 610CE", Category: "Guitars"},
			{Name: "Paradym Loud", Category: "Entertainment"},
		},
	}
	got, err := provider.ListByCategory("Guitars")
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if len(got) != 1 {
		t.Fatalf("got len(items) = %d, want 1", len(got))
	}
	if got[0] != provider.Repository[0] {
		t.Errorf("got item = %v, want %v", got, provider.Repository[0])
	}
}

func Test_TotalByCategory(t *testing.T) {
	provider := &DatabaseProvider{
		Repository: []InventoryItem{
			{Name: "KRK Studio Monitors", Value: 179, Quantity: 2, Category: "Music Gear"},
			{Name: "Macbook Pro", Value: 2070, Quantity: 1, Category: "Computers"},
		},
	}
	want := 358
	got, err := provider.TotalByCategory("Music Gear")
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if got != want {
		t.Errorf("got = %d, want %d", got, want)
	}
}
