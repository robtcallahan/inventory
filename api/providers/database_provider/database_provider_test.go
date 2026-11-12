package database_provider

import (
	"fmt"
	"testing"
)

type mockProvider struct {
	items []InventoryItem
}

func (mp *mockProvider) Save(item InventoryItem) error {
	mp.items = append(mp.items, item)
	return nil
}
func (mp *mockProvider) List() ([]InventoryItem, error) {
	return mp.items, nil
}

func Test_SaveAndList(t *testing.T) {
	mp := &mockProvider{}
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

	err := mp.Save(item)
	if err != nil {
		fmt.Errorf("got error %v, want nil", err)
	}
	items, err := mp.List()
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
