package database_provider

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"inventory/pkg/models"
)

type fakeRepository struct {
	fakeID           int
	fakeItems        []models.Item
	fakeItem         *models.Item
	fakeOutputItems  []models.OutputItem
	fakeError        error
	fakeSavePhoto    *models.ItemPhoto
	fakeSavePhotoErr error
	fakePhotos       []models.ItemPhoto
	fakePhotosErr    error
}

func (p *fakeRepository) Save(item *models.Item) (*models.Item, error) {
	p.fakeItem = item
	p.fakeItems = append(p.fakeItems, *item)
	return nil, p.fakeError
}
func (p *fakeRepository) Get(id int) (models.Item, error) {
	p.fakeID = id
	return models.Item{}, nil
}
func (p *fakeRepository) SaveItems(items []models.Item) ([]models.Item, error) {
	p.fakeItems = append(p.fakeItems, items...)
	return nil, nil
}
func (p *fakeRepository) Delete(item models.Item) error {
	p.fakeItem = &item
	return p.fakeError
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

func (p *fakeRepository) SavePhoto(itemPhoto *models.ItemPhoto) (*models.ItemPhoto, error) {
	p.fakeSavePhoto = itemPhoto
	return itemPhoto, p.fakeSavePhotoErr
}
func (p *fakeRepository) ListPhotos(itemID int) ([]models.ItemPhoto, error) {
	p.fakeID = itemID
	return p.fakePhotos, p.fakePhotosErr
}
func (p *fakeRepository) DeletePhoto(photoID int) error {
	p.fakeID = photoID
	return p.fakeError
}

func Test_Save(t *testing.T) {
	fakeRepo := &fakeRepository{fakeItems: []models.Item{}}
	provider := &DatabaseProvider{Repository: fakeRepo}

	now := time.Now()
	item := &models.Item{
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
			want:    item,
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
			_, err := provider.Save(item)
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

func Test_Delete(t *testing.T) {
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
		wantErr bool
		fakeErr error
	}{
		{
			name:    "no error",
			wantErr: false,
		},
		{
			name:    "error",
			wantErr: true,
			fakeErr: errors.New("could not delete item"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeRepo.fakeError = tt.fakeErr
			err := provider.Delete(item)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
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
				TotalValue: 2070, Category: "Computers"},
			{Name: "KRK Studio Monitors", Description: "KRK Studio Monitors, yellow", Value: 179, Quantity: 2,
				TotalValue: 358, Category: "Music Gear"},
			{Name: "Humidifier", Description: "TommyBoy white humidifier", Value: 250, Quantity: 1, TotalValue: 250, Category: "Appliances"},
			{Name: "Chef's knife", Description: "Flextone 10\" Chef's knife", Value: 50, Quantity: 2, TotalValue: 100,
				Category: "Kitchen"},
			{Name: "Chef's knife", Description: "Harbinger 8\" Chef's knife", Value: 50, Quantity: 2, TotalValue: 100,
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

func Test_SavePhoto(t *testing.T) {
	tests := []struct {
		name    string
		fakeErr error
	}{
		{name: "success"},
		{name: "repository error", fakeErr: errors.New("could not save photo")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeRepo := &fakeRepository{
				fakeSavePhotoErr: tt.fakeErr,
			}
			provider := &DatabaseProvider{Repository: fakeRepo}

			want := &models.ItemPhoto{
				ItemID:   42,
				Filename: "photo1.jpg",
			}

			got, err := provider.SavePhoto(want)

			if !errors.Is(err, tt.fakeErr) {
				t.Fatalf("SavePhoto() error = %v, want %v", err, tt.fakeErr)
			}

			if fakeRepo.fakeSavePhoto != want {
				t.Errorf("SavePhoto() passed %p, want %p",
					fakeRepo.fakeSavePhoto, want)
			}

			if tt.fakeErr == nil && !reflect.DeepEqual(got, want) {
				t.Errorf("SavePhoto() = %+v, want %+v", got, want)
			}
		})
	}
}

func Test_ListPhotos(t *testing.T) {
	tests := []struct {
		name    string
		photos  []models.ItemPhoto
		fakeErr error
	}{
		{
			name: "multiple photos",
			photos: []models.ItemPhoto{
				{ID: 1, ItemID: 42, Filename: "photo1.jpg"},
				{ID: 2, ItemID: 42, Filename: "photo2.jpg"},
			},
		},
		{
			name:   "no photos",
			photos: []models.ItemPhoto{},
		},
		{
			name:    "repository error",
			fakeErr: errors.New("could not list photos"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeRepo := &fakeRepository{
				fakePhotos:    tt.photos,
				fakePhotosErr: tt.fakeErr,
			}
			provider := &DatabaseProvider{Repository: fakeRepo}

			const itemID = 42
			got, err := provider.ListPhotos(itemID)

			if !errors.Is(err, tt.fakeErr) {
				t.Fatalf("ListPhotos() error = %v, want %v", err, tt.fakeErr)
			}

			if fakeRepo.fakeID != itemID {
				t.Errorf("ListPhotos() itemID = %d, want %d",
					fakeRepo.fakeID, itemID)
			}

			if tt.fakeErr == nil && !reflect.DeepEqual(got, tt.photos) {
				t.Errorf("ListPhotos() = %+v, want %+v", got, tt.photos)
			}
		})
	}
}

func Test_DeletePhoto(t *testing.T) {
	tests := []struct {
		name    string
		fakeErr error
	}{
		{name: "success"},
		{name: "repository error", fakeErr: errors.New("could not delete photo")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeRepo := &fakeRepository{
				fakeError: tt.fakeErr,
			}
			provider := &DatabaseProvider{Repository: fakeRepo}

			const photoID = 17
			err := provider.DeletePhoto(photoID)

			if !errors.Is(err, tt.fakeErr) {
				t.Fatalf("DeletePhoto() error = %v, want %v", err, tt.fakeErr)
			}

			if fakeRepo.fakeID != photoID {
				t.Errorf("DeletePhoto() photoID = %d, want %d",
					fakeRepo.fakeID, photoID)
			}
		})
	}
}
