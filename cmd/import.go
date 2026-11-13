package cmd

import (
	"inventory/api/providers/database_provider"
	"inventory/pkg/csv"
	"inventory/pkg/models"
	"inventory/pkg/repository"
	"inventory/pkg/repository/mysql"

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
	importCmd.Flags().StringVar(&options.ImportFileName, "file", "./inventory.csv", "import CSV file name")
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

	conn, err := repository.Connect(&config.Database)
	if err != nil {
		return err
	}
	qHandler := mysql.NewMySQLQueryRepo(conn)
	provider := database_provider.DatabaseProvider{Repository: qHandler}
	return provider.SaveItems(items)
}
