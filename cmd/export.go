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

func init() {
	rootCmd.AddCommand(exportCmd)
}

func export() error {
	conn, err := repository.Connect(&config.Database)
	if err != nil {
		return err
	}
	qHandler := mysql.NewMySQLQueryRepo(conn)
	items, err := qHandler.So

	csvPath := filepath.Join(ProjectDir, dirName, CSVFileName)
	file, err := os.Create(csvPath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", csvPath, err)
	}
	defer file.Close()

	w := stdcsv.NewWriter(file)
	// Quantity	Name	Manufacturer	Model	Value	Category	Purchase Date	Serial Number
	record := []string{"Qty", "Name", "Manufacturer", "Model", "Value", "Category", "Purchase Date", "Serial Num"}
	if err := w.Write(record); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	for _, t := range p.Repository {
		record = []string{
			strconv.Itoa(t.Quantity),
			t.Name,
			t.Manufacturer,
			t.Model,
			strconv.Itoa(t.Value * t.Quantity),
			t.Category,
			t.PurchaseDate,
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
