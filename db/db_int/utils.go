package db_int

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/errs"
	"github.com/fossmedaddy/dbdaddy/globals"
	"github.com/fossmedaddy/dbdaddy/types"

	"github.com/spf13/viper"
)

// i dont like this
func GetRows(queryStr string) (types.QueryResult, error) {
	queryResult := types.QueryResult{
		Data: types.DbRows{},
	}

	q, err := globals.DB.Query(queryStr)
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

	defer q.Close()
	for q.Next() {
		err := q.Scan(vals...)
		if err != nil {
			return queryResult, err
		}
		for i := range vals {
			valPtr := vals[i].(*interface{})
			val := *valPtr
			if b, ok := val.([]byte); ok {
				s := string(b)
				// Optional: You can check if the string seems to be Base64 encoded (length, valid characters)
				if isValidBase64(s) {
					decoded, err := base64.StdEncoding.DecodeString(s)
					if err != nil {
						fmt.Println("Failed to decode:", s, "Error:", err)
						return queryResult, err
					}
					fmt.Println(string(decoded))
					val = decoded
				} else {
					// If it's not a valid Base64, print it as is
					fmt.Println(s)
					val = s
				}
			} else {
				fmt.Println(val)
			}
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

// takes in a sql string & differentiates between multiple executable statements
// by a ";" followed by a newline.
// EXPECTS AN INPUT GENERATED FROM 'libUtils.GetSqlStmts()' func
func ExecuteStatementsTx(stmts []string) error {
	tx, err := globals.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, stmt := range stmts {
		if _, err := tx.Exec(stmt); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// no transactions, rawdawg them queries right in the DB
func ExecuteStatements(stmts []string) error {
	for _, stmt := range stmts {
		if _, err := globals.DB.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}

func isValidBase64(s string) bool {
	// Basic checks: length should be a multiple of 4 and contain only valid Base64 characters
	if len(s)%4 != 0 {
		return false
	}
	for _, c := range s {
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '+' || c == '/' || c == '=' {
			continue
		}
		return false
	}
	return true
}
