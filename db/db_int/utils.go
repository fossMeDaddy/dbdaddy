package db_int

import (
	constants "dbdaddy/const"
	"dbdaddy/db"
	"dbdaddy/errs"
	"dbdaddy/types"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

// i dont like this
func GetRows(queryStr string) (types.QueryResult, error) {
	queryResult := types.QueryResult{
		Data: types.DbRows{},
	}

	q, err := db.DB.Query(queryStr)
	if err != nil {
		return queryResult, err
	}

	cols, err := q.Columns()
	if err != nil {
		return queryResult, err
	}
	queryResult.Columns = cols

	colTypes, err := q.ColumnTypes()
	if err != nil {
		return queryResult, err
	}

	vals := make([]interface{}, len(cols))
	for i := range cols {
		var ii interface{}
		vals[i] = &ii

		// initialize db rows
		queryResult.Data[cols[i]] = []types.DbRow{}
	}

	for q.Next() {
		err := q.Scan(vals...)
		if err != nil {
			return queryResult, err
		}

		for i := range vals {
			valPtr := vals[i].(*interface{})
			val := *valPtr

			if val == nil {
				queryResult.Data[cols[i]] = append(queryResult.Data[cols[i]], types.DbRow{
					DataType: colTypes[i].DatabaseTypeName(),
					Value:    val,
				})
			} else if strings.Contains(strings.ToLower(colTypes[i].DatabaseTypeName()), "json") {
				jsonBytesVal := val.([]byte)
				dbRow := types.DbRow{
					DataType: colTypes[i].DatabaseTypeName(),
					Value:    string(jsonBytesVal),
				}
				queryResult.Data[cols[i]] = append(queryResult.Data[cols[i]], dbRow)
			} else {
				dbRow := types.DbRow{
					DataType: colTypes[i].DatabaseTypeName(),
					Value:    val,
				}
				queryResult.Data[cols[i]] = append(queryResult.Data[cols[i]], dbRow)
			}
		}

		queryResult.RowCount++
	}

	return queryResult, nil
}

func ExecuteStatements__DEPRECATED(dbname, sqlStr string) error {
	if _, err := exec.LookPath("psql"); err != nil {
		return errs.ErrPsqlCmdNotFound
	}

	cmd := exec.Command("psql", "--single-transaction")

	inPipe, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if _, err := inPipe.Write([]byte(sqlStr)); err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("PGHOST=%s", viper.GetString(constants.DbConfigHostKey)),
		fmt.Sprintf("PGPORT=%s", viper.GetString(constants.DbConfigPortKey)),
		fmt.Sprintf("PGUSER=%s", viper.GetString(constants.DbConfigUserKey)),
		fmt.Sprintf("PGPASSWORD=%s", viper.GetString(constants.DbConfigPassKey)),
		fmt.Sprintf("PGDATABASE=%s", dbname),
	)

	return cmd.Run()
}

func ExecuteStatements(sqlStr string) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	for _, stmt := range strings.Split(sqlStr, ";"+fmt.Sprintln()) {
		if _, err := tx.Exec(stmt); err != nil {
			return err
		}
	}

	return tx.Commit()
}
