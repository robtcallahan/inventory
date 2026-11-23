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

func OutputItemsToHTML(items []OutputItem) []HTMLItem {
	htmlItems := make([]HTMLItem, 0, len(items))
	for i, item := range items {
		purchaseDate := ""
		if item.PurchaseDate != nil {
			purchaseDate = item.PurchaseDate.Format("01/02/2006")
		}
		htmlItem := HTMLItem{
			ID:           item.ID,
			Number:       i + 1,
			Quantity:     item.Quantity,
			Name:         item.Name,
			ShortDescription:  item.ShortDescription,
			ReplacementValue:        item.ReplacementValue,
			ResaleValue:  item.ResaleValue,
			TotalValue:   item.TotalValue,
			Manufacturer: item.Manufacturer,
			Model:        item.Model,
			Category:     item.Category,
			PurchaseDate: purchaseDate,
			SerialNumber: item.SerialNumber,
			Room: item.Room,
			Condition: item.Condition,
			PurchaseSource: item.PurchaseSource,
			PurchasePrice: item.PurchasePrice,
			WebsiteURL: item.WebsiteURL,
			ManualURL: item.ManualURL,
			Details:      item.Details,
			PhotoCount:   item.PhotoCount,
		}
		htmlItems = append(htmlItems, htmlItem)
	}
	return htmlItems
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
			ReplacementValue:        int(i.ReplacementValue),
			ResaleValue:  int(i.ResaleValue),
			ShortDescription:  i.ShortDescription,
			Manufacturer: i.Manufacturer,
			Model:        i.Model,
			Category:     i.Category,
			PurchaseDate: purchasedAt,
			SerialNumber: i.SerialNumber,
			Room: i.Room,
			Condition: i.Condition,
			PurchaseSource: i.PurchaseSource,
			PurchasePrice: int(i.PurchasePrice),
			WebsiteURL: i.WebsiteURL,
			ManualURL: i.ManualURL,
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
	record := []string{"Qty", "Name", "Replacement Value", "Resale Value", "Total Cost", "Short Description", "Manufacturer", "Model", "Category",
		"Purchase Date", "Serial Num", "Details", "Room", "Condition", "Purchase Source", "Purchase Price", "Website URL", "Manual URL"}
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
			strconv.Itoa(t.ReplacementValue),
			strconv.Itoa(t.ResaleValue),
			strconv.Itoa(t.TotalValue),
			t.ShortDescription,
			t.Manufacturer,
			t.Model,
			t.Category,
			purchaseDate,
			t.SerialNumber,
			t.Details,
			t.Room,
			t.Condition,
			t.PurchaseSource,
			strconv.Itoa(t.PurchasePrice),
			t.WebsiteURL,
			t.ManualURL,
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
