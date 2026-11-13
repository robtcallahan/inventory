package models

import (
	stdcsv "encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

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

func ItemsExportToCSVFile(fileName string, items []Item) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", fileName, err)
	}
	defer file.Close()

	w := stdcsv.NewWriter(file)
	record := []string{"Qty", "Name", "Manufacturer", "Model", "Value", "Category", "Purchase Date", "Serial Num", "Description"}
	if err = w.Write(record); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	for _, t := range items {
		purchaseDate := ""
		if t.PurchaseDate != nil {
			purchaseDate = t.PurchaseDate.Format("01/02/2006")
		}
		record = []string{
			strconv.Itoa(t.Quantity),
			t.Name,
			t.Manufacturer,
			t.Model,
			strconv.Itoa(t.Value * t.Quantity),
			t.Category,
			purchaseDate,
			t.SerialNumber,
			t.Description,
		}
		if err := w.Write(record); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
