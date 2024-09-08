package lib

import (
	"fmt"
	"strings"

	"github.com/fossmedaddy/dbdaddy/types"
)

func GetCsvString(queryResult types.QueryResult) string {
	csvFile := ""

	cols := []string{}
	for col := range queryResult.Data {
		cols = append(cols, col)
	}

	rowCount := len(queryResult.Data[cols[0]])

	// write header
	csvFile += fmt.Sprintln(strings.Join(cols, ","))

	// write body
	for i := range rowCount {
		rowStr := ""
		for _, col := range cols {
			rowStr += fmt.Sprint(queryResult.Data[col][i].Value) + ","
		}
		rowStr = strings.TrimRight(rowStr, ",")
		csvFile += fmt.Sprintln(rowStr)
	}

	return csvFile
}

func GetFormattedColumns(queryResult types.QueryResult) string {
	COL_CAP := 100

	text := ""

	// get padding data
	colPadding := map[string]int{}
	for _, col := range queryResult.Columns {
		padding := len(col)
		for _, item := range queryResult.Data[col] {
			padding = max(len(fmt.Sprint(item.Value)), padding)
		}
		padding = min(padding, COL_CAP)

		colPadding[col] = padding
	}

	// header
	for _, col := range queryResult.Columns {
		text += printFormattedColumn(col, colPadding[col], COL_CAP)
	}
	text += fmt.Sprintln()

	text += fmt.Sprintln(strings.Repeat("-", len(text)-1))

	for i := range queryResult.RowCount {
		for _, col := range queryResult.Columns {
			text += printFormattedColumn(fmt.Sprint(queryResult.Data[col][i].Value), colPadding[col], COL_CAP)
		}
		text += fmt.Sprintln()
	}

	return text
}

func printFormattedColumn(data string, padding int, cap int) string {
	if len(data) > cap {
		data = data[:cap-3] + "..."
	}
	return fmt.Sprintf(
		" %-*s |",
		padding, data,
	)
}
