package migrationsLib

import (
	constants "dbdaddy/const"
	"dbdaddy/lib"
	"dbdaddy/types"
	"slices"
	"sync"
)

var (
	globalWg sync.WaitGroup
)

func getKeysFromState(state types.DbSchema, tag string) ([][]string, [][]string) {
	tableKeys := [][]string{}
	colKeys := [][]string{}

	for mapKey := range state.Tables {
		tableSchema := state.Tables[mapKey]

		tableKey := []string{tag, tableSchema.Db, tableSchema.Schema, tableSchema.Name}
		tableKeys = append(tableKeys, tableKey)

		for _, col := range tableSchema.Columns {
			colKey := append(tableKey, col.Name)
			colKeys = append(colKeys, colKey)
		}
	}

	return tableKeys, colKeys
}

// takes in CS & PS concatenated sorted table keys.
// returns changes & mapping for what tables were created/deleted
func getTableStateChanges(tableStateKeysConcat [][]string) ([]types.MigAction, map[string]bool) {
	tableChanges := []types.MigAction{}
	tableChangesMapping := map[string]bool{}

	for _, tableKey := range tableStateKeysConcat {
		if tableKey[0] == currentStateTag {
			// CREATE
			psTableKey := slices.Concat([]string{prevStateTag}, tableKey[1:])
			_, found := slices.BinarySearchFunc(tableStateKeysConcat, psTableKey, slices.Compare)
			if !found {
				action := types.MigAction{
					Type:     constants.MigActionCreateTable,
					EntityId: tableKey,
				}

				tableChanges = append(tableChanges, action)
				tableChangesMapping[lib.GetTableId(tableKey[2], tableKey[3])] = true
			}
		} else if tableKey[0] == prevStateTag {
			// DROP
			csTableKey := slices.Concat([]string{currentStateTag}, tableKey[1:])
			_, found := slices.BinarySearchFunc(tableStateKeysConcat, csTableKey, slices.Compare)
			if !found {
				action := types.MigAction{
					Type:     constants.MigActionDropTable,
					EntityId: tableKey,
				}

				tableChanges = append(tableChanges, action)
				tableChangesMapping[lib.GetTableId(tableKey[2], tableKey[3])] = true
			}
		}
	}

	return tableChanges, tableChangesMapping
}

func getColStateChanges(colStateKeysConcat [][]string, tableChangesMapping map[string]bool) []types.MigAction {
	localChanges := []types.MigAction{}

	for _, colKey := range colStateKeysConcat {
		// if table itself was changed (created/deleted)
		// there's no point of tracking the columns
		tableid := lib.GetTableId(colKey[2], colKey[3])
		if tableChangesMapping[tableid] {
			continue
		}

		if colKey[0] == currentStateTag {
			// CREATE
			psColKey := slices.Concat([]string{prevStateTag}, colKey[1:])
			_, found := slices.BinarySearchFunc(colStateKeysConcat, psColKey, slices.Compare)
			if !found {
				action := types.MigAction{
					Type:     constants.MigActionCreateCol,
					EntityId: colKey,
				}

				localChanges = append(localChanges, action)
			}
		} else if colKey[0] == prevStateTag {
			// DROP
			csColKey := slices.Concat([]string{currentStateTag}, colKey[1:])
			_, found := slices.BinarySearchFunc(colStateKeysConcat, csColKey, slices.Compare)
			if !found {
				action := types.MigAction{
					Type:     constants.MigActionDropCol,
					EntityId: colKey,
				}

				localChanges = append(localChanges, action)
			}
		}
	}

	return localChanges
}

/*
give changes to be done on 'prevState' in order to move from 'prevState' to 'currentState'

v0.1 - very simple, CREATE OR DELETE (DEFINITELY NOT FOR PRODUCTION DATABASES)
*/
func DiffDbSchema(currentState, prevState types.DbSchema) []types.MigAction {
	changes := []types.MigAction{}

	var (
		tableKeysCS, colKeysCS [][]string
		tableKeysPS, colKeysPS [][]string

		tableStateKeysConcat [][]string
		colStateKeysConcat   [][]string
	)

	// get keys from state
	(func() {
		globalWg.Add(2)
		defer globalWg.Wait()

		go (func() {
			defer globalWg.Done()
			tableKeysCS, colKeysCS = getKeysFromState(currentState, currentStateTag)
		})()

		go (func() {
			defer globalWg.Done()
			tableKeysPS, colKeysPS = getKeysFromState(prevState, prevStateTag)
		})()
	})()

	// concat keys & sort
	(func() {
		globalWg.Add(2)
		defer globalWg.Wait()

		go (func() {
			defer globalWg.Done()
			tableStateKeysConcat = slices.Concat(tableKeysCS, tableKeysPS)
			slices.SortFunc(tableStateKeysConcat, slices.Compare)
		})()

		go (func() {
			defer globalWg.Done()
			colStateKeysConcat = slices.Concat(colKeysCS, colKeysPS)
			slices.SortFunc(colStateKeysConcat, slices.Compare)
		})()
	})()

	// accumulate changes keys
	tableChanges, tableChangesMapping := getTableStateChanges(tableStateKeysConcat)
	colChanges := getColStateChanges(colStateKeysConcat, tableChangesMapping)
	changes = slices.Concat(changes, tableChanges, colChanges)

	return changes
}
