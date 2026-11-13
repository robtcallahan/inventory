package database_provider

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"inventory/pkg/models"
)

type fakeRepository struct {
	fakeItems       []models.Item
	fakeItem        *models.Item
	fakeOutputItems []models.OutputItem
	fakeError       error
}

func (p *fakeRepository) Save(item models.Item) error {
	p.fakeItem = &item
	p.fakeItems = append(p.fakeItems, item)
	return p.fakeError
}
func (p *fakeRepository) SaveItems(items []models.Item) error {
	p.fakeItems = append(p.fakeItems, items...)
	return nil
}
func (p *fakeRepository) List() ([]models.Item, error) {
	return p.fakeItems, p.fakeError
}
func (p *fakeRepository) ListByCategory(category string) ([]models.Item, error) {
	var matched []models.Item
	for _, item := range p.fakeItems {
		if item.Category == category {
			matched = append(matched, item)
		}
	}
	return matched, p.fakeError
}
func (p *fakeRepository) TotalByCategory(category string) (int, error) {
	totalValue := 0
	for _, item := range p.fakeItems {
		if item.Category == category {
			totalValue = totalValue + (item.Value * item.Quantity)
		}
	}
	return totalValue, p.fakeError
}
func (p *fakeRepository) SortedListForOutput() ([]models.OutputItem, error) {
	return p.fakeOutputItems, p.fakeError
}

func Test_Save(t *testing.T) {
	fakeRepo := &fakeRepository{fakeItems: []models.Item{}}
	provider := &DatabaseProvider{Repository: fakeRepo}

	now := time.Now()
	item := models.Item{
		Quantity:     1,
		Name:         "Item 1",
		Manufacturer: "Manufacturer",
		Model:        "Model",
		Description:  "Description",
		SerialNumber: "Serial Number",
		Category:     "Category",
		PurchaseDate: &now,
		Value:        100,
	}

	tests := []struct {
		name    string
		want    *models.Item
		wantErr bool
		fakeErr error
	}{
		{
			name:    "no error",
			want:    &item,
			wantErr: false,
		},
		{
			name:    "error",
			want:    nil,
			wantErr: true,
			fakeErr: errors.New("could not save item"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeRepo.fakeError = tt.fakeErr
			err := provider.Save(item)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(fakeRepo.fakeItem, tt.want) {
				t.Errorf("Save() = %v, want %v", fakeRepo.fakeItem, tt.want)
			}
		})
	}

}

func Test_List(t *testing.T) {
	fakeRepo := &fakeRepository{
		fakeItems: []models.Item{
			{Name: "Taylor 610CE", Category: "Guitars"},
			{Name: "Paradym Loud", Category: "Entertainment"},
		},
	}
	provider := &DatabaseProvider{Repository: fakeRepo}

	got, err := provider.List()
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if len(got) != 2 {
		t.Fatalf("got len(items) = %d, want 2", len(got))
	}
	if !reflect.DeepEqual(got, fakeRepo.fakeItems) {
		t.Errorf("got item = %v, want %v", got, fakeRepo.fakeItems)
	}
}

func Test_ListByCategory(t *testing.T) {
	fakeRepo := &fakeRepository{
		fakeItems: []models.Item{
			{Name: "Taylor 610CE", Category: "Musical Instruments"},
			{Name: "Paradym Loud", Category: "Entertainment"},
		},
	}
	provider := &DatabaseProvider{Repository: fakeRepo}

	got, err := provider.ListByCategory("Musical Instruments")
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if len(got) != 1 {
		t.Fatalf("got len(items) = %d, want 1", len(got))
	}
	if got[0] != fakeRepo.fakeItems[0] {
		t.Errorf("got item = %v, want %v", got, fakeRepo.fakeItems[0])
	}
}

func Test_TotalByCategory(t *testing.T) {
	fakeRepo := &fakeRepository{
		fakeItems: []models.Item{
			{Name: "KRK Studio Monitors", Value: 179, Quantity: 2, Category: "Music Gear"},
			{Name: "Macbook Pro", Value: 2070, Quantity: 1, Category: "Computers"},
		},
	}
	provider := &DatabaseProvider{Repository: fakeRepo}

	want := 358
	got, err := provider.TotalByCategory("Music Gear")
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if got != want {
		t.Errorf("got = %d, want %d", got, want)
	}
}

func Test_SortedListForOutput(t *testing.T) {
	fakeRepo := &fakeRepository{
		fakeOutputItems: []models.OutputItem{
			{Name: "MacBook Pro", Description: "2019 MacBook Pro, 16GB, 512 SSD", Value: 2070, Quantity: 1,
				TotalCost: 2070, Category: "Computers"},
			{Name: "KRK Studio Monitors", Description: "KRK Studio Monitors, yellow", Value: 179, Quantity: 2,
				TotalCost: 358, Category: "Music Gear"},
			{Name: "Humidifier", Description: "TommyBoy white humidifier", Value: 250, Quantity: 1, TotalCost: 250, Category: "Appliances"},
			{Name: "Chef's knife", Description: "Flextone 10\" Chef's knife", Value: 50, Quantity: 2, TotalCost: 100,
				Category: "Kitchen"},
			{Name: "Chef's knife", Description: "Harbinger 8\" Chef's knife", Value: 50, Quantity: 2, TotalCost: 100,
				Category: "Appliances"},
		},
	}
	provider := &DatabaseProvider{Repository: fakeRepo}

	got, err := provider.SortedListForOutput()
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if len(got) != 5 {
		t.Fatalf("got len(items) = %d, want 5", len(got))
	}
	if got[0].Name != fakeRepo.fakeOutputItems[0].Name {
		t.Errorf("got item[0].Name = %v, want %v", got[0].Name, fakeRepo.fakeOutputItems[0].Name)
	}
	if got[4].Description != fakeRepo.fakeOutputItems[4].Description {
		t.Errorf("got item[4].Name = %v, want %v", got[4].Description, fakeRepo.fakeOutputItems[4].Description)
	}
}
