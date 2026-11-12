package cmd

import (
	"inventory/pkg/csv"
	"inventory/pkg/models"
	"inventory/pkg/repository"
	"inventory/pkg/repository/mysql"
	"time"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import items from a CSV file into the database.",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return importCSV()
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().BoolVar(&options.DryRun, "dry-run", false, "Print-only run: no DB or spreadsheet updates performed")
	importCmd.Flags().StringVar(&options.CSVFileName, "csv", ".", "directory for inventory.csv export file")
}

func importCSV() error {
	var purchasedAt *time.Time

	csvItems, err := csv.ImportCSV(options.CSVFileName)
	if err != nil {
		return err
	}

	conn, err := repository.Connect(&config.Database)
	if err != nil {
		return err
	}
	qHandler := mysql.NewMySQLQueryRepo(conn)

	for _, i := range csvItems {
		purchasedAt = nil
		if i.PurchaseDate != "" {
			if t, err := time.Parse("01/02/2006", i.PurchaseDate); err == nil {
				purchasedAt = &t
			}
		}
		item := models.Item{
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
		if err := qHandler.Save(item); err != nil {
			return err
		}
	}
	return nil
}
