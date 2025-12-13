package file

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/xuri/excelize/v2"
)

func ReadExcelRows(data []byte, sheetIndex, minRow int) ([][]string, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to read excel file: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(sheetIndex)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}
	if len(rows) < minRow {
		return nil, errors.New("excel file has no data rows")
	}

	return rows, nil
}
