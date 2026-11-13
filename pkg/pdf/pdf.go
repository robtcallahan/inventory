package pdf

import (
	"fmt"
	"inventory/pkg/models"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/template"
)

func PDF(fileName string, items []models.Item) error {
	var rows = make([][]string, len(items))

	for i, t := range items {
		purchaseDate := ""
		if t.PurchaseDate != nil {
			purchaseDate = t.PurchaseDate.Format("01/02/2006")
		}
		rows = append(rows, []string{
			strconv.Itoa(i + 1),
			strconv.Itoa(t.Quantity),
			t.Name,
			makeBlank(t.Manufacturer),
			makeBlank(t.Model),
			formatAccounting(float64(t.Value*t.Quantity), 8),
			makeBlank(t.Category),
			makeBlank(purchaseDate),
			//makeBlank(t.SerialNumber),
			//t.Description,
		})
	}
	reportData := template.ReportData{
		Title:    "Callahan Home Inventory",
		Subtitle: "November 12, 2026",
		Author:   "Robert Callahan",
		Date:     "414 Diehl St, Winchester, VA 22601",
		Sections: []template.ReportSection{
			{
				Title: "Items",
				Table: &template.ReportTable{
					//Header:       []string{"Qty", "Name", "Manufacturer", "Model", "Cost", "Category", "Pur Date", "Serial"},
					//ColumnWidths: []float64{5, 20, 15, 5, 10, 20, 10, 10},
					Header:       []string{" ", "Qty", "Name", "Manufacturer", "Model", "Cost", "Category", "Pur Date"},
					ColumnWidths: []float64{2, 3, 40, 10, 15, 10, 10, 10},
					Rows:         rows,
				},
			},
		},
	}
	doc := template.Report(reportData,
		template.WithPageSize(document.Letter.Landscape()),
		template.WithDefaultFont("Arial", 8),
	)

	data, err := doc.Generate()
	if err != nil {
		panic(err)
	}
	return os.WriteFile(fileName, data, 0644)
}

// formatAccounting formats a float64 in accounting style without decimals,
// aligning the '$' to the left and the number to the right within the given width.
func formatAccounting(amount float64, width int) string {
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

func makeBlank(cell string) string {
	if cell == "" {
		return " "
	}
	return cell
}
