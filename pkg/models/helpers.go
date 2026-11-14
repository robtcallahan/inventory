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
			Quantity:     i.Quantity,
			Name:         i.Name,
			Value:        int(i.Value),
			Description:  i.Description,
			Manufacturer: i.Manufacturer,
			Model:        i.Model,
			Category:     i.Category,
			PurchaseDate: purchasedAt,
			SerialNumber: i.SerialNumber,
			Details:      i.Details,
		}
		items = append(items, item)
	}
	return items, nil
}

func ItemsExportToCSVFile(fileName string, items []OutputItem) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", fileName, err)
	}
	defer file.Close()

	w := stdcsv.NewWriter(file)
	record := []string{"Qty", "Name", "Value", "Total Cost", "Description", "Manufacturer", "Model", "Category",
		"Purchase Date", "Serial Num", "Details"}
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
			strconv.Itoa(t.Value),
			strconv.Itoa(t.TotalCost),
			t.Description,
			t.Manufacturer,
			t.Model,
			t.Category,
			purchaseDate,
			t.SerialNumber,
			t.Details,
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
