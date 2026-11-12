package cmd

import (
	stdcsv "encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"inventory/pkg/repository"
	"inventory/pkg/repository/mysql"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export items to a CSV file called inventory.csv",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return export()
	},
}

const (
	ProjectDir  = "/Users/rob/ws/go/src/inventory"
	CSVFileName = "inventory.csv"
)

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().BoolVar(&options.DryRun, "dry-run", false, "Print-only run: no DB or spreadsheet updates performed")
	exportCmd.Flags().StringVar(&options.ExportDirectory, "dir", ".", "directory for inventory.csv export file")
}

func export() error {
	conn, err := repository.Connect(&config.Database)
	if err != nil {
		return err
	}
	qHandler := mysql.NewMySQLQueryRepo(conn)
	items, err := qHandler.SortedListForCSV()
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return fmt.Errorf("no items found")
	}

	//csvPath := filepath.Join(ProjectDir, options.ExportDirectory, CSVFileName)
	csvPath := filepath.Join(options.ExportDirectory, CSVFileName)
	file, err := os.Create(csvPath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", csvPath, err)
	}
	defer file.Close()

	w := stdcsv.NewWriter(file)
	// Quantity	Name	Manufacturer	Model	Value	Category	Purchase Date	Serial Number
	record := []string{"Qty", "Name", "Manufacturer", "Model", "Value", "Category", "Purchase Date", "Serial Num"}
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
