package cmd

import (
	"fmt"
	"inventory/api/providers/database_provider"
	"inventory/pkg/csv"
	"inventory/pkg/models"

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
	importCmd.Flags().StringVar(&options.ImportFileName, "file", "inventory.csv", "import CSV file name")
}

func importCSV() error {
	csvItems, err := csv.ImportCSV(options.ImportFileName)
	if err != nil {
		return err
	}

	items, err := models.CSVToItem(csvItems)
	if err != nil {
		return err
	}

	provider, err := database_provider.NewMySQLProvider(&config.Database)
	if err != nil {
		return err
	}
	skipped, err := provider.SaveItems(items)
	if err != nil {
		return err
	}

	// print out any duplicate items that were skipped so they can be corrected in the spreadsheet
	if len(skipped) != 0 {
		fmt.Println("Item skipped due to duplicate name/description unique key:")
		fmt.Printf("Row %-30s Description\n", "Name")
		for _, skippedItem := range skipped {
			for i, item := range items {
				// find the skipped item in the items slice to determine the row which is i+2
				if item.Name == skippedItem.Name && item.Description == skippedItem.Description {
					fmt.Printf("%3d %-30s %s\n", i+2, item.Name, item.Description)
					break
				}
			}
		}
	}
	return nil
}
