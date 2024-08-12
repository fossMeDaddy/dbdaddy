package migrationsLib

import (
	constants "dbdaddy/const"
	"dbdaddy/types"
	"slices"
	"sync"
)

const (
	currentStateTag = "CS"
	prevStateTag    = "PS"
)

var (
	globalWg sync.WaitGroup
)

func getKeysFromState(state types.DbSchemaMapping, tag string) ([][]string, [][]string) {
	tableKeys := [][]string{}
	colKeys := [][]string{}

	for mapKey := range state {
		tableSchema := state[mapKey]

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
func DiffDbSchema(currentState, prevState types.DbSchemaMapping) []types.MigAction {
	safeChanges := types.NewSafeVar([]types.MigAction{})

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
	(func() {
		globalWg.Add(2)
		defer globalWg.Wait()

		go (func() {
			defer globalWg.Done()

			localChanges := []types.MigAction{}

			for _, tableKey := range tableStateKeysConcat {
				if tableKey[0] == currentStateTag {
					// CREATE
					psTableKey := slices.Concat([]string{prevStateTag}, tableKey[1:])
					_, found := slices.BinarySearchFunc(tableStateKeysConcat, psTableKey, slices.Compare)
					if !found {
						action := types.MigAction{
							Type:     constants.MigActionCreateTable,
							EntityId: tableKey[1:],
						}

						localChanges = append(localChanges, action)
					}
				} else if tableKey[0] == prevStateTag {
					// DROP
					csTableKey := slices.Concat([]string{currentStateTag}, tableKey[1:])
					_, found := slices.BinarySearchFunc(tableStateKeysConcat, csTableKey, slices.Compare)
					if !found {
						action := types.MigAction{
							Type:     constants.MigActionDropTable,
							EntityId: tableKey[1:],
						}

						localChanges = append(localChanges, action)
					}
				}
			}

			safeChanges.Set(localChanges)
		})()

		go (func() {
			defer globalWg.Done()

			localChanges := []types.MigAction{}

			for _, colKey := range colStateKeysConcat {
				if colKey[0] == currentStateTag {
					// CREATE
					psColKey := slices.Concat([]string{prevStateTag}, colKey[1:])
					_, found := slices.BinarySearchFunc(colStateKeysConcat, psColKey, slices.Compare)
					if !found {
						action := types.MigAction{
							Type:     constants.MigActionCreateCol,
							EntityId: colKey[1:],
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
							EntityId: colKey[1:],
						}

						localChanges = append(localChanges, action)
					}
				}
			}

			safeChanges.Set(localChanges)
		})()
	})()

	return safeChanges.Get()
}
