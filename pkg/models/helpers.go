package models

import "time"

func CSVToItem(csvItems []CSVItem) ([]Item, error) {
	var items []Item
	var purchasedAt *time.Time

	for _, i := range csvItems {
		purchasedAt = nil
		if i.PurchaseDate != "" {
			if t, err := time.Parse("01/02/2006", i.PurchaseDate); err == nil {
				purchasedAt = &t
			}
		}
		item := Item{
			Category:     i.Category,
			PurchaseDate: purchasedAt,
			Value:        i.Value,
			Name:         i.Name,
			Model:        i.Model,
			Manufacturer: i.Manufacturer,
			Quantity:     i.Quantity,
			Description:  i.Description,
			SerialNumber: i.SerialNumber,
		}
		items = append(items, item)
	}
	return items, nil
}
