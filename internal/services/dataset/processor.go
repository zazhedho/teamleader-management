package servicedataset

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"teamleader-management/pkg/file"
	"time"

	domaindataset "teamleader-management/internal/domain/dataset"
	domainmetric "teamleader-management/internal/domain/metric"
	interfacedataset "teamleader-management/internal/interfaces/dataset"
	interfacemetric "teamleader-management/internal/interfaces/metric"
	interfaceperson "teamleader-management/internal/interfaces/person"
	"teamleader-management/utils"
)

type Processor struct {
	DatasetRepo interfacedataset.RepoDatasetInterface
	MetricRepo  interfacemetric.RepoMetricInterface
	PersonRepo  interfaceperson.RepoPersonInterface
}

func NewProcessor(datasetRepo interfacedataset.RepoDatasetInterface, metricRepo interfacemetric.RepoMetricInterface, personRepo interfaceperson.RepoPersonInterface) *Processor {
	return &Processor{
		DatasetRepo: datasetRepo,
		MetricRepo:  metricRepo,
		PersonRepo:  personRepo,
	}
}

func parseDoneStatus(val string, rowIndex int) (bool, error) {
	switch val {
	case "DONE", "YES", "Y", "1", "TRUE":
		return true, nil
	case "NOT DONE", "NO", "N", "0", "FALSE":
		return false, nil
	default:
		return false, fmt.Errorf("row %d invalid status value: %s", rowIndex+1, val)
	}
}

// ProcessStream executes import flow from in-memory file data.
func (p *Processor) ProcessStream(ds *domaindataset.DashboardDataset, data []byte, actorId string) (domaindataset.DashboardDataset, error) {
	if ds.Status != utils.DatasetStatusProcessing && ds.Status != utils.DatasetStatusUploaded && ds.Status != utils.DatasetStatusFailed {
		return *ds, errors.New("dataset is not in a processable state")
	}

	if err := p.markStatus(ds, utils.DatasetStatusProcessing, actorId); err != nil {
		return *ds, err
	}

	var processErr error
	switch strings.ToUpper(ds.Type) {
	case utils.DatasetQuiz:
		processErr = p.processQuiz(ds, data, actorId)
	case utils.DatasetLoginApple:
		processErr = p.processAppleLogin(ds, data, actorId)
	case utils.DatasetSalesFLP:
		processErr = p.processSalesFLP(ds, data, actorId)
	case utils.DatasetPointApple:
		processErr = p.processApplePoints(ds, data, actorId)
	case utils.DatasetPointMyHero:
		processErr = p.processMyHeroPoints(ds, data, actorId)
	case utils.DatasetTotalProspect:
		processErr = p.processProspects(ds, data, actorId)
	default:
		processErr = fmt.Errorf("dataset type %s not supported yet", ds.Type)
	}

	if processErr != nil {
		_ = p.markStatus(ds, utils.DatasetStatusFailed, actorId)
		return *ds, processErr
	}

	if err := p.markStatus(ds, utils.DatasetStatusDone, actorId); err != nil {
		return *ds, err
	}

	return *ds, nil
}

func (p *Processor) markStatus(ds *domaindataset.DashboardDataset, status, actorId string) error {
	now := time.Now()
	ds.Status = status
	ds.UpdatedAt = now
	ds.UpdatedBy = actorId
	return p.DatasetRepo.Update(*ds)
}

func (p *Processor) processQuiz(ds *domaindataset.DashboardDataset, data []byte, actorId string) error {
	rows, err := file.ReadExcelRows(data, 0, 2)
	if err != nil {
		return err
	}

	results := make([]domainmetric.QuizResult, 0, len(rows)-1)
	now := time.Now()
	for i, row := range rows {
		if i == 0 {
			continue // header
		}
		if len(row) < 8 { // No, Honda ID, Nama, Jabatan, Role, Kode Dealer, Nilai, Lulus/Tidak, (Status optional)
			return fmt.Errorf("row %d has insufficient columns", i+1)
		}
		hondaId := strings.TrimSpace(row[1])
		if hondaId == "" {
			return fmt.Errorf("row %d missing Honda ID", i+1)
		}
		person, err := p.PersonRepo.GetByHondaID(hondaId)
		if err != nil {
			return fmt.Errorf("row %d person not found for honda_id %s", i+1, hondaId)
		}

		var dealerCode *string
		if len(row) >= 6 && strings.TrimSpace(row[5]) != "" {
			val := strings.TrimSpace(row[5])
			dealerCode = &val
		}

		var score *float64
		if len(row) >= 7 && strings.TrimSpace(row[6]) != "" {
			val, convErr := strconv.ParseFloat(strings.TrimSpace(row[6]), 64)
			if convErr != nil {
				return fmt.Errorf("row %d invalid score: %v", i+1, convErr)
			}
			score = &val
		}

		var passStatus *string
		if len(row) >= 8 && strings.TrimSpace(row[7]) != "" {
			val := strings.TrimSpace(row[7])
			passStatus = &val
		}

		result := domainmetric.QuizResult{
			Id:         utils.CreateUUID(),
			DatasetId:  ds.Id,
			PeriodDate: ds.PeriodDate,
			PersonId:   person.Id,
			HondaId:    hondaId,
			DealerCode: dealerCode,
			Score:      score,
			PassStatus: passStatus,
			CreatedAt:  now,
			CreatedBy:  actorId,
			UpdatedAt:  now,
			UpdatedBy:  actorId,
		}
		results = append(results, result)
	}

	if err := p.MetricRepo.SaveQuizResults(results); err != nil {
		return err
	}

	return nil
}

func (p *Processor) processAppleLogin(ds *domaindataset.DashboardDataset, data []byte, actorId string) error {
	rows, err := file.ReadExcelRows(data, 0, 2)
	if err != nil {
		return err
	}

	logins := make([]domainmetric.AppleLogin, 0, len(rows)-1)
	now := time.Now()
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 7 { // No, Honda ID, Nama Lengkap, Jabatan, Kode Dealer, Frequent PAGI, Frequent SORE
			return fmt.Errorf("row %d has insufficient columns", i+1)
		}
		hondaId := strings.TrimSpace(row[1])
		if hondaId == "" {
			return fmt.Errorf("row %d missing Honda ID", i+1)
		}
		person, err := p.PersonRepo.GetByHondaID(hondaId)
		if err != nil {
			return fmt.Errorf("row %d person not found for honda_id %s", i+1, hondaId)
		}

		var dealerCode *string
		if len(row) >= 5 && strings.TrimSpace(row[4]) != "" {
			val := strings.TrimSpace(row[4])
			dealerCode = &val
		}

		morningVal := strings.ToUpper(strings.TrimSpace(row[5]))
		eveningVal := strings.ToUpper(strings.TrimSpace(row[6]))

		morningDone, err := parseDoneStatus(morningVal, i)
		if err != nil {
			return err
		}
		eveningDone, err := parseDoneStatus(eveningVal, i)
		if err != nil {
			return err
		}

		login := domainmetric.AppleLogin{
			Id:          utils.CreateUUID(),
			DatasetId:   ds.Id,
			PeriodDate:  ds.PeriodDate,
			PersonId:    person.Id,
			HondaId:     hondaId,
			DealerCode:  dealerCode,
			LoginDate:   ds.PeriodDate,
			MorningDone: morningDone,
			EveningDone: eveningDone,
			CreatedAt:   now,
			CreatedBy:   actorId,
			UpdatedAt:   now,
			UpdatedBy:   actorId,
		}
		logins = append(logins, login)
	}

	if err := p.MetricRepo.SaveAppleLogins(logins); err != nil {
		return err
	}

	return nil
}

func (p *Processor) processSalesFLP(ds *domaindataset.DashboardDataset, data []byte, actorId string) error {
	rows, err := file.ReadExcelRows(data, 0, 2)
	if err != nil {
		return err
	}

	entries := make([]domainmetric.SalesFLP, 0, len(rows)-1)
	now := time.Now()
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 5 { // No, Honda ID, Kode Dealer, Nama Sales People, Sales
			return fmt.Errorf("row %d has insufficient columns", i+1)
		}

		hondaId := strings.TrimSpace(row[1])
		if hondaId == "" {
			return fmt.Errorf("row %d missing Honda ID", i+1)
		}
		person, err := p.PersonRepo.GetByHondaID(hondaId)
		if err != nil {
			return fmt.Errorf("row %d person not found for honda_id %s", i+1, hondaId)
		}

		var dealerCode *string
		if val := strings.TrimSpace(row[2]); val != "" {
			dealerCode = &val
		}

		amountStr := strings.TrimSpace(row[4])
		if amountStr == "" {
			return fmt.Errorf("row %d sales amount is required", i+1)
		}
		amount, err := utils.ParseInt(amountStr)
		if err != nil {
			return fmt.Errorf("row %d invalid %s value: Sales", i+1, err)
		}

		entry := domainmetric.SalesFLP{
			Id:         utils.CreateUUID(),
			DatasetId:  ds.Id,
			PeriodDate: ds.PeriodDate,
			PersonId:   person.Id,
			HondaId:    hondaId,
			DealerCode: dealerCode,
			Amount:     amount,
			CreatedAt:  now,
			CreatedBy:  actorId,
			UpdatedAt:  now,
			UpdatedBy:  actorId,
		}
		entries = append(entries, entry)
	}

	if err := p.MetricRepo.SaveSalesFLP(entries); err != nil {
		return err
	}

	return nil
}

func (p *Processor) processMyHeroPoints(ds *domaindataset.DashboardDataset, data []byte, actorId string) error {
	rows, err := file.ReadExcelRows(data, 0, 2)
	if err != nil {
		return err
	}

	entries := make([]domainmetric.MyHeroPoint, 0, len(rows)-1)
	now := time.Now()

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 5 { // No, Honda Id, Kode Dealer, Nama, Jumlah Poin
			return fmt.Errorf("row %d has insufficient columns", i+1)
		}

		hondaId := strings.TrimSpace(row[1])
		if hondaId == "" {
			return fmt.Errorf("row %d missing Honda ID", i+1)
		}

		person, err := p.PersonRepo.GetByHondaID(hondaId)
		if err != nil {
			return fmt.Errorf("row %d person not found for honda_id %s", i+1, hondaId)
		}

		var dealerCode *string
		if val := strings.TrimSpace(row[2]); val != "" {
			dealerCode = &val
		}

		pointsStr := strings.TrimSpace(row[4])
		if pointsStr == "" {
			return fmt.Errorf("row %d jumlah poin is required", i+1)
		}
		points, err := utils.ParseInt(pointsStr)
		if err != nil {
			return fmt.Errorf("row %d invalid %s value: %v", i+1, "Jumlah Poin", err)
		}

		entry := domainmetric.MyHeroPoint{
			Id:         utils.CreateUUID(),
			DatasetId:  ds.Id,
			PeriodDate: ds.PeriodDate,
			PersonId:   person.Id,
			HondaId:    hondaId,
			DealerCode: dealerCode,
			Points:     points,
			CreatedAt:  now,
			CreatedBy:  actorId,
			UpdatedAt:  now,
			UpdatedBy:  actorId,
		}
		entries = append(entries, entry)
	}

	if err := p.MetricRepo.SaveMyHeroPoints(entries); err != nil {
		return err
	}

	return nil
}

func (p *Processor) processProspects(ds *domaindataset.DashboardDataset, data []byte, actorId string) error {
	rows, err := file.ReadExcelRows(data, 0, 2)
	if err != nil {
		return err
	}

	entries := make([]domainmetric.Prospect, 0, len(rows)-1)
	now := time.Now()

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 4 { // No, Honda ID, Nama Sales, Prospek
			return fmt.Errorf("row %d has insufficient columns", i+1)
		}

		hondaId := strings.TrimSpace(row[1])
		if hondaId == "" {
			return fmt.Errorf("row %d missing Honda ID", i+1)
		}
		person, err := p.PersonRepo.GetByHondaID(hondaId)
		if err != nil {
			return fmt.Errorf("row %d person not found for honda_id %s", i+1, hondaId)
		}

		prospectStr := strings.TrimSpace(row[3])
		if prospectStr == "" {
			return fmt.Errorf("row %d prospek is required", i+1)
		}
		prospectCount, err := utils.ParseInt(prospectStr)
		if err != nil {
			return fmt.Errorf("row %d invalid Prospek value: %v", i+1, err)
		}

		entry := domainmetric.Prospect{
			Id:            utils.CreateUUID(),
			DatasetId:     ds.Id,
			PeriodDate:    ds.PeriodDate,
			PersonId:      person.Id,
			HondaId:       hondaId,
			ProspectCount: prospectCount,
			CreatedAt:     now,
			CreatedBy:     actorId,
			UpdatedAt:     now,
			UpdatedBy:     actorId,
		}
		entries = append(entries, entry)
	}

	if err := p.MetricRepo.SaveProspects(entries); err != nil {
		return err
	}

	return nil
}

func (p *Processor) processApplePoints(ds *domaindataset.DashboardDataset, data []byte, actorId string) error {
	rows, err := file.ReadExcelRows(data, 0, 1)
	if err != nil {
		return err
	}

	entries := make([]domainmetric.ApplePoint, 0, len(rows)-1)
	now := time.Now()

	for i, row := range rows {
		if i == 0 {
			continue // header
		}
		if len(row) < 3 { // Honda Id, Nama Sales, Point Apple
			return fmt.Errorf("row %d has insufficient columns", i+1)
		}

		hondaId := strings.TrimSpace(row[0])
		if hondaId == "" {
			return fmt.Errorf("row %d missing Honda ID", i+1)
		}

		person, err := p.PersonRepo.GetByHondaID(hondaId)
		if err != nil {
			return fmt.Errorf("row %d person not found for honda_id %s", i+1, hondaId)
		}

		pointsStr := strings.TrimSpace(row[2])
		if pointsStr == "" {
			return fmt.Errorf("row %d point apple is required", i+1)
		}
		points, err := utils.ParseInt(pointsStr)
		if err != nil {
			return fmt.Errorf("row %d invalid %s value: %v", i+1, "Point Apple", err)
		}

		entry := domainmetric.ApplePoint{
			Id:         utils.CreateUUID(),
			DatasetId:  ds.Id,
			PeriodDate: ds.PeriodDate,
			PersonId:   person.Id,
			HondaId:    hondaId,
			Points:     points,
			CreatedAt:  now,
			CreatedBy:  actorId,
			UpdatedAt:  now,
			UpdatedBy:  actorId,
		}
		entries = append(entries, entry)
	}

	if err := p.MetricRepo.SaveApplePoints(entries); err != nil {
		return err
	}

	return nil
}
