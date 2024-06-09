package db_int

import (
	"dbdaddy/db"
	"dbdaddy/types"
	"fmt"
	"strings"
)

// i dont like this
func GetRows(queryStr string) (types.DbRows, error) {
	q, err := db.DB.Query(queryStr)
	if err != nil {
		return types.DbRows{}, err
	}

	cols, err := q.Columns()
	if err != nil {
		return types.DbRows{}, err
	}

	colTypes, err := q.ColumnTypes()
	if err != nil {
		return types.DbRows{}, err
	}

	data := types.DbRows{}

	vals := make([]interface{}, len(cols))
	for i := range cols {
		var ii interface{}
		vals[i] = &ii

		// initialize db rows
		data[cols[i]] = []types.DbRow{}
	}
	for q.Next() {
		err := q.Scan(vals...)
		if err != nil {
			return types.DbRows{}, err
		}

		for i := range vals {
			valPtr := vals[i].(*interface{})
			val := *valPtr

			if val == nil {
				data[cols[i]] = append(data[cols[i]], types.DbRow{
					DataType: colTypes[i].DatabaseTypeName(),
					StrValue: "",
				})
			} else if strings.Contains(strings.ToLower(colTypes[i].DatabaseTypeName()), "json") {
				jsonBytesVal := val.([]byte)
				dbRow := types.DbRow{
					DataType: colTypes[i].DatabaseTypeName(),
					StrValue: string(jsonBytesVal),
				}
				data[cols[i]] = append(data[cols[i]], dbRow)
			} else {
				dbRow := types.DbRow{
					DataType: colTypes[i].DatabaseTypeName(),
					StrValue: fmt.Sprint(val),
				}
				data[cols[i]] = append(data[cols[i]], dbRow)
			}
		}
	}

	return data, nil
}
