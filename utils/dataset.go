package utils

var AllowedDatasetTypes = map[string]bool{
	DatasetQuiz:          true,
	DatasetLoginApple:    true,
	DatasetSalesFLP:      true,
	DatasetPointApple:    true,
	DatasetPointMyHero:   true,
	DatasetTotalProspect: true,
}

var AllowedPeriodFrequencies = map[string]bool{
	"DAILY":     true,
	"WEEKLY":    true,
	"MONTHLY":   true,
	"QUARTERLY": true,
	"YEARLY":    true,
}

const (
	DatasetStatusUploaded   = "UPLOADED"
	DatasetStatusProcessing = "PROCESSING"
	DatasetStatusDone       = "DONE"
	DatasetStatusFailed     = "FAILED"
)
