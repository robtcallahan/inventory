package cmd

import (
	"fmt"
	"inventory/api/providers/database_provider"
	"inventory/pkg/models"
	"inventory/pkg/pdf"

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

	exportCmd.Flags().BoolVar(&options.DryRun, "dry-run", false, "Print-only run: no DB or spreadsheet updates performed")
	exportCmd.Flags().StringVar(&options.OutputFormat, "format", "pdf", "export format (pdf or csv")
	exportCmd.Flags().StringVar(&options.ExportFileName, "file", "inventory.pdf", "file name for export")
}

func export() error {
	provider, err := database_provider.NewMySQLProvider(&config.Database)
	if err != nil {
		return err
	}

	items, err := provider.SortedListForOutput()
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return fmt.Errorf("no items found")
	}

	switch options.OutputFormat {
	case "pdf":
		return pdf.PDF(options.ExportFileName, items)
	case "csv":
		return models.ItemsExportToCSVFile(options.ExportFileName, items)
	default:
		return fmt.Errorf("invalid output format: %s, 'pdf' or 'csv' allowed", options.OutputFormat)
	}
}
