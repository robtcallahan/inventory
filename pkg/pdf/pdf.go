package pdf

import (
	"inventory/pkg/models"
	"os"
	"strconv"

	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/template"
)

func PDF(fileName string, items []models.OutputItem) error {
	var rows = make([][]string, 0, len(items))

	for i, t := range items {
		purchaseDate := ""
		if t.PurchaseDate != nil {
			purchaseDate = t.PurchaseDate.Format("01/02/2006")
		}
		rows = append(rows, []string{
			strconv.Itoa(i + 1),
			strconv.Itoa(t.Quantity),
			t.Name,
			t.Description,
			makeBlank(t.Manufacturer),
			makeBlank(t.Model),
			models.FormatAccounting(float64(t.TotalCost), 8),
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
					Header:       []string{" ", "Qty", "Name", "Descr", "Man", "Model", "Cost", "Category", "Pur Date"},
					ColumnWidths: []float64{2, 3, 10, 20, 15, 20, 10, 10, 10},
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
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

func makeBlank(cell string) string {
	if cell == "" {
		return " "
	}
	return cell
}
