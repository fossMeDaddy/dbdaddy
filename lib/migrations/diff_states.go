package migrationsLib

import (
	constants "dbdaddy/const"
	"dbdaddy/types"
	"fmt"
	"slices"
)

const (
	currentStateTag = "CS"
	prevStateTag    = "PS"
)

func getKeysFromState(state types.DbSchema, tag string) ([][]string, [][]string) {
	tableKeys := [][]string{}
	colKeys := [][]string{}

	for _, tableSchema := range state {
		tableKey := []string{tag, tableSchema.Db, tableSchema.Schema, tableSchema.Name}
		tableKeys = append(tableKeys, tableKey)

		for _, col := range tableSchema.Columns {
			colKey := append(tableKey, col.Name)
			colKeys = append(colKeys, colKey)
		}
	}

	return tableKeys, colKeys
}

/*
give changes to be done on 'prevState' in order to move from 'prevState' to 'currentState'

v0.1 - very simple, CREATE OR DELETE (DEFINITELY NOT FOR PRODUCTION DATABASES)
*/
func DiffDbSchema(currentState, prevState types.DbSchema) {
	changes := []types.MigAction{}

	tableKeysCS, _ := getKeysFromState(currentState, currentStateTag)
	tableKeysPS, _ := getKeysFromState(prevState, prevStateTag)

	tableStateKeysConcat := slices.Concat(tableKeysCS, tableKeysPS)
	slices.SortFunc(tableStateKeysConcat, slices.Compare)

	// TABLE KEYS
	for _, tableKey := range tableStateKeysConcat {
		if tableKey[0] == currentStateTag {
			psTableKey := slices.Concat([]string{prevStateTag}, tableKey[1:])
			_, found := slices.BinarySearchFunc(tableStateKeysConcat, psTableKey, slices.Compare)
			if !found {
				action := types.MigAction{
					Type:     constants.MigActionCreateTable,
					EntityId: tableKey[1:],
				}

				changes = append(changes, action)
			}
		} else if tableKey[0] == prevStateTag {
			csTableKey := slices.Concat([]string{currentStateTag}, tableKey[1:])
			_, found := slices.BinarySearchFunc(tableStateKeysConcat, csTableKey, slices.Compare)
			if !found {
				action := types.MigAction{
					Type:     constants.MigActionDropTable,
					EntityId: tableKey[1:],
				}

				changes = append(changes, action)
			}
		}
	}

	fmt.Println(changes)
}
