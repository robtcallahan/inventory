package config

const (
	OutputFormatPDF = "pdf"
	OutputFormatCSV = "csv"
)

type Options struct {
	ExportFileName string
	ImportFileName string
	OutputFormat   string
	DryRun         bool
}
