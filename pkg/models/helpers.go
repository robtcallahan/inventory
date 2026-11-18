package models

import (
	stdcsv "encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func OutputItemsToHTML(items []OutputItem) ([]HTMLItem, error) {
	htmlItems := make([]HTMLItem, 0, len(items))
	for i, item := range items {
		purchaseDate := ""
		if item.PurchaseDate != nil {
			purchaseDate = item.PurchaseDate.Format("01/02/2006")
		}
		htmlItem := HTMLItem{
			ID:           item.ID,
			Number:       strconv.Itoa(i + 1),
			Quantity:     strconv.Itoa(item.Quantity),
			Name:         item.Name,
			Description:  item.Description,
			Value:        item.Value,
			TotalCost:    item.TotalCost,
			Manufacturer: item.Manufacturer,
			Model:        item.Model,
			Category:     item.Category,
			PurchaseDate: purchaseDate,
			SerialNumber: item.SerialNumber,
			Details:      item.Details,
		}
		htmlItems = append(htmlItems, htmlItem)
	}
	return htmlItems, nil
}

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

// FormatAccounting formats a float64 in accounting style without decimals,
// aligning the '$' to the left and the number to the right within the given width.
func FormatAccounting(amount float64, width int) string {
	// Round to nearest whole number
	rounded := int64(math.Round(amount))
	isNegative := rounded < 0
	if isNegative {
		rounded = -rounded
	}

	// Format zero as a dash
	var numStr string
	if rounded == 0 {
		numStr = "-"
	} else {
		numStr = addCommas(rounded)
	}

	// Wrap negative numbers in parentheses
	if isNegative {
		numStr = fmt.Sprintf("(%s)", numStr)
	} else {
		// Add a trailing space so positive numbers align with the ')' of negative numbers
		numStr = numStr + " "
	}

	// Calculate padding between '$' and the number
	// Total width minus 1 for '$' and len(numStr)
	padding := width - 1 - len(numStr)
	if padding < 1 {
		padding = 1 // Ensure at least one space after '$'
	}

	return fmt.Sprintf("$%s%s", strings.Repeat(" ", padding), numStr)
}

// addCommas inserts thousand separators into a positive integer.
func addCommas(n int64) string {
	s := fmt.Sprintf("%d", n)
	var result []string
	for i := len(s); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		result = append([]string{s[start:i]}, result...)
	}
	return strings.Join(result, ",")
}
