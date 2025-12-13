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
	PeriodDaily:     true,
	PeriodWeekly:    true,
	PeriodMonthly:   true,
	PeriodQuarterly: true,
	PeriodYearly:    true,
}

const (
	DatasetStatusUploaded   = "UPLOADED"
	DatasetStatusProcessing = "PROCESSING"
	DatasetStatusDone       = "DONE"
	DatasetStatusFailed     = "FAILED"
)
