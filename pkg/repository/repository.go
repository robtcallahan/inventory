package repository

import (
	"inventory/pkg/models"
)

// QueryRepo represent the repositories
type QueryRepo interface {
	Save(item models.Item) error
	List() ([]models.Item, error)
	ListByCategory(category string) ([]models.Item, error)
	TotalByCategory(category string) (int, error)
	SortedListForCSV() ([]models.Item, error)
}
