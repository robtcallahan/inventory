package repository

import (
	"inventory/pkg/models"
)

// QueryRepo represent the repositories
type QueryRepo interface {
	Save(item *models.Item) (*models.Item, error)
	Get(id int) (models.Item, error)
	SaveItems(items []models.Item) ([]models.Item, error)
	Delete(item models.Item) error
	List() ([]models.Item, error)
	ListByCategory(category string) ([]models.Item, error)
	TotalByCategory(category string) (int, error)
	// SortedListForOutput includes photo_count for each item (zero when no photos exist).
	SortedListForOutput() ([]models.OutputItem, error)
	SavePhoto(itemPhoto *models.ItemPhoto) (*models.ItemPhoto, error)
	ListPhotos(itemID int) ([]models.ItemPhoto, error)
	DeletePhoto(photoID int) error
}
