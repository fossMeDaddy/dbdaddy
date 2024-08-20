package migrationsLib

import (
	constants "dbdaddy/const"
	"dbdaddy/libUtils"
	"dbdaddy/types"
	"slices"
	"strings"
	"sync"

	"golang.org/x/exp/maps"
)

type (
	keyType [][]string

	groupedActionsType map[string][]types.MigAction
)

var (
	changeTypePrefixSortingOrder = map[string]int{
		"DROP":   0,
		"CREATE": 1,
	}
)

func getKeysFromState(state *types.DbSchema, tag string) (keyType, keyType, keyType, keyType, keyType) {
	keysWg := sync.WaitGroup{}

	tableKeys := keyType{}
	viewKeys := keyType{}
	colKeys := keyType{}
	typeKeys := keyType{}
	conKeys := keyType{}

	(func() {
		keysWg.Add(3)
		defer keysWg.Wait()

		go (func() {
			defer keysWg.Done()

			for _, typeKey := range state.Types {
				typeKeys = append(typeKeys, []string{tag, typeKey.Schema, typeKey.Name})
			}
		})()

		go (func() {
			defer keysWg.Done()

			for mapKey := range state.Tables {
				tableSchema := state.Tables[mapKey]

				tableKey := []string{tag, state.DbName, tableSchema.Schema, tableSchema.Name}
				tableKeys = append(tableKeys, tableKey)

				for _, col := range tableSchema.Columns {
					colKey := append(tableKey, col.Name)
					colKeys = append(colKeys, colKey)
				}

				for _, con := range tableSchema.Constraints {
					conKey := append(tableKey, con.ConName, con.Type, con.Syntax)
					conKeys = append(conKeys, conKey)
				}
			}
		})()

		go (func() {
			defer keysWg.Done()

			for viewKey := range state.Views {
				viewSchema := state.Views[viewKey]

				viewKey := []string{tag, state.DbName, viewSchema.Schema, viewSchema.Name}
				viewKeys = append(viewKeys, viewKey)
			}
		})()
	})()

	return tableKeys, viewKeys, colKeys, typeKeys, conKeys
}

// takes in CS & PS concatenated sorted table keys.
// returns changes & mapping for what tables were created/deleted
func getTableStateChanges(tableStateKeysConcat keyType) ([]types.MigAction, map[string]string) {
	tableChanges := []types.MigAction{}

	// tableid is mapped to a string containing same values as "MigAction*" Strings
	tableChangesMapping := map[string]string{}

	for _, tableKey := range tableStateKeysConcat {
		if tableKey[0] == currentStateTag {
			psTableKey := slices.Concat([]string{prevStateTag}, tableKey[1:])
			_, found := slices.BinarySearchFunc(tableStateKeysConcat, psTableKey, slices.Compare)
			if !found {
				// CREATE
				action := types.MigAction{
					Type:     constants.MigActionCreateTable,
					EntityId: tableKey,
				}

				tableChanges = append(tableChanges, action)
				tableChangesMapping[libUtils.GetTableId(tableKey[2], tableKey[3])] = action.Type
			}
		} else if tableKey[0] == prevStateTag {
			csTableKey := slices.Concat([]string{currentStateTag}, tableKey[1:])
			_, found := slices.BinarySearchFunc(tableStateKeysConcat, csTableKey, slices.Compare)
			if !found {
				// DROP
				action := types.MigAction{
					Type:     constants.MigActionDropTable,
					EntityId: tableKey,
				}

				tableChanges = append(tableChanges, action)
				tableChangesMapping[libUtils.GetTableId(tableKey[2], tableKey[3])] = action.Type
			}
		}
	}

	return tableChanges, tableChangesMapping
}

func getViewStateChanges(viewStateKeysConcat keyType, tableChangesMapping map[string]string) []types.MigAction {
	viewChanges := []types.MigAction{}

	for _, viewKey := range viewStateKeysConcat {
		if viewKey[0] == currentStateTag {
			psViewKey := slices.Concat([]string{prevStateTag}, viewKey[1:])
			_, found := slices.BinarySearchFunc(viewStateKeysConcat, psViewKey, slices.Compare)
			if !found {
				// CREATE
				action := types.MigAction{
					Type:     constants.MigActionCreateView,
					EntityId: viewKey,
				}

				viewChanges = append(viewChanges, action)
			}
		} else if viewKey[0] == prevStateTag {
			csViewKey := slices.Concat([]string{currentStateTag}, viewKey[1:])
			_, found := slices.BinarySearchFunc(viewStateKeysConcat, csViewKey, slices.Compare)
			if !found {
				// if table was removed, with "CASCADE" added (which we're doing currently)
				// dependent views also will be removed
				tableid := libUtils.GetTableId(viewKey[2], viewKey[3])
				if getChangeTypePrefix(tableChangesMapping[tableid]) == "DROP" {
					continue
				}

				// DROP
				action := types.MigAction{
					Type:     constants.MigActionDropView,
					EntityId: viewKey,
				}

				viewChanges = append(viewChanges, action)
			}
		}
	}

	return viewChanges
}

func getColStateChanges(colStateKeysConcat keyType, tableChangesMapping map[string]string) []types.MigAction {
	localChanges := []types.MigAction{}

	for _, colKey := range colStateKeysConcat {
		// if table itself was changed (created/deleted)
		// there's no point of tracking the columns
		tableid := libUtils.GetTableId(colKey[2], colKey[3])
		if slices.Contains([]string{constants.MigActionCreateTable, constants.MigActionDropTable}, tableChangesMapping[tableid]) {
			continue
		}

		if colKey[0] == currentStateTag {
			psColKey := slices.Concat([]string{prevStateTag}, colKey[1:])
			_, found := slices.BinarySearchFunc(colStateKeysConcat, psColKey, slices.Compare)
			if !found {
				// CREATE
				action := types.MigAction{
					Type:     constants.MigActionCreateCol,
					EntityId: colKey,
				}

				localChanges = append(localChanges, action)
			}
		} else if colKey[0] == prevStateTag {
			csColKey := slices.Concat([]string{currentStateTag}, colKey[1:])
			_, found := slices.BinarySearchFunc(colStateKeysConcat, csColKey, slices.Compare)
			if !found {
				// DROP
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

func getTypeStateChanges(typeStateKeysConcat keyType) []types.MigAction {
	localChanges := []types.MigAction{}

	for _, typeKey := range typeStateKeysConcat {
		if typeKey[0] == currentStateTag {
			psTypeKey := slices.Concat([]string{prevStateTag}, typeKey[1:])
			_, found := slices.BinarySearchFunc(typeStateKeysConcat, psTypeKey, slices.Compare)
			if !found {
				// CREATE
				action := types.MigAction{
					Type:     constants.MigActionCreateType,
					EntityId: typeKey,
				}

				localChanges = append(localChanges, action)
			}
		} else if typeKey[0] == prevStateTag {
			csTypeKey := slices.Concat([]string{currentStateTag}, typeKey[1:])
			_, found := slices.BinarySearchFunc(typeStateKeysConcat, csTypeKey, slices.Compare)
			if !found {
				// DROP
				action := types.MigAction{
					Type:     constants.MigActionDropType,
					EntityId: typeKey,
				}

				localChanges = append(localChanges, action)
			}
		}
	}

	return localChanges
}

func getConStateChanges(conStateKeysConcat keyType, tableChangesMapping map[string]string) []types.MigAction {
	localChanges := []types.MigAction{}

	for _, conKey := range conStateKeysConcat {
		tableid := libUtils.GetTableId(conKey[2], conKey[3])
		if tableChangesMapping[tableid] == constants.MigActionDropTable {
			continue
		}

		if conKey[0] == currentStateTag {
			psConKey := slices.Concat([]string{prevStateTag}, conKey[1:])
			_, found := slices.BinarySearchFunc(conStateKeysConcat, psConKey, slices.Compare)
			if !found {
				// CREATE
				change := types.MigAction{
					Type:     constants.MigActionCreateConstraint,
					EntityId: conKey,
				}

				localChanges = append(localChanges, change)
			}
		} else if conKey[0] == prevStateTag {
			csConKey := slices.Concat([]string{currentStateTag}, conKey[1:])
			_, found := slices.BinarySearchFunc(conStateKeysConcat, csConKey, slices.Compare)
			if !found {
				// DROP
				change := types.MigAction{
					Type:     constants.MigActionDropConstraint,
					EntityId: conKey,
				}

				localChanges = append(localChanges, change)
			}
		}
	}

	return localChanges
}

func getChangeTypePrefix(changeType string) string {
	prefix := strings.Split(changeType, "_")[0]
	return prefix
}

func changeTypePrefixSortingCmp(a, b types.MigAction) int {
	aPrefix := getChangeTypePrefix(a.Type)
	bPrefix := getChangeTypePrefix(b.Type)

	return changeTypePrefixSortingOrder[aPrefix] - changeTypePrefixSortingOrder[bPrefix]
}

func groupActionsByType(changes []types.MigAction) groupedActionsType {
	groupedActions := groupedActionsType{}
	for _, action := range changes {
		groupedActions[action.Type] = append(groupedActions[action.Type], action)
	}

	return groupedActions
}

func pickActionGroups(keys []string, actions []types.MigAction) ([]types.MigAction, []types.MigAction) {
	groupedActions := groupActionsByType(actions)

	segment := []types.MigAction{}
	antiSegment := []types.MigAction{}

	allKeys := maps.Keys(groupedActions)
	keyPresenceMapping := map[string]bool{}
	for _, key := range keys {
		keyPresenceMapping[key] = true
	}

	for _, key := range allKeys {
		if keyPresenceMapping[key] {
			segment = slices.Concat(segment, groupedActions[key])
		} else {
			antiSegment = slices.Concat(antiSegment, groupedActions[key])
		}
	}

	return segment, antiSegment
}

/*
give changes to be done on 'prevState' in order to move from 'prevState' to 'currentState'

v0.1 - very simple, CREATE OR DELETE (DEFINITELY NOT FOR PRODUCTION DATABASES)
*/
func DiffDbSchema(currentState, prevState *types.DbSchema) []types.MigAction {
	var wg sync.WaitGroup

	changes := []types.MigAction{}

	var (
		tableKeysCS, tableKeysPS keyType
		viewKeysCS, viewKeysPS   keyType
		colKeysCS, colKeysPS     keyType
		conKeysCS, conKeysPS     keyType
		typeKeysCS, typeKeysPS   keyType

		tableStateKeysConcat keyType
		viewStateKeysConcat  keyType
		colStateKeysConcat   keyType
		typeStateKeysConcat  keyType
		conStateKeysConcat   keyType
	)

	// get keys from state
	(func() {
		wg.Add(2)
		defer wg.Wait()

		go (func() {
			defer wg.Done()
			tableKeysCS, viewKeysCS, colKeysCS, typeKeysCS, conKeysCS = getKeysFromState(currentState, currentStateTag)
		})()

		go (func() {
			defer wg.Done()
			tableKeysPS, viewKeysPS, colKeysPS, typeKeysPS, conKeysPS = getKeysFromState(prevState, prevStateTag)
		})()
	})()

	// concat keys & sort
	(func() {
		wg.Add(5)
		defer wg.Wait()

		go (func() {
			defer wg.Done()
			tableStateKeysConcat = slices.Concat(tableKeysCS, tableKeysPS)
			slices.SortFunc(tableStateKeysConcat, slices.Compare)
		})()

		go (func() {
			defer wg.Done()
			colStateKeysConcat = slices.Concat(colKeysCS, colKeysPS)
			slices.SortFunc(colStateKeysConcat, slices.Compare)
		})()

		go (func() {
			defer wg.Done()
			typeStateKeysConcat = slices.Concat(typeKeysCS, typeKeysPS)
			slices.SortFunc(typeStateKeysConcat, slices.Compare)
		})()

		go (func() {
			defer wg.Done()
			conStateKeysConcat = slices.Concat(conKeysCS, conKeysPS)
			slices.SortFunc(conStateKeysConcat, slices.Compare)
		})()

		go (func() {
			defer wg.Done()
			viewStateKeysConcat = slices.Concat(viewKeysCS, viewKeysPS)
			slices.SortFunc(viewStateKeysConcat, slices.Compare)
		})()
	})()

	// accumulate changes keys

	// TODO: TYPES SQL NOT SUPPORTED CURRENTLY
	_ = getTypeStateChanges(typeStateKeysConcat)

	tableChanges, tableChangesMapping := getTableStateChanges(tableStateKeysConcat)
	slices.SortFunc(tableChanges, changeTypePrefixSortingCmp)

	dropViewChanges, otherViewChanges := pickActionGroups(
		[]string{constants.MigActionDropView},
		getViewStateChanges(viewStateKeysConcat, tableChangesMapping),
	)

	colChanges := getColStateChanges(colStateKeysConcat, tableChangesMapping)
	slices.SortFunc(colChanges, changeTypePrefixSortingCmp)

	// sort constraints in an order that prevents SQL errors during migration run
	conChanges := getConStateChanges(conStateKeysConcat, tableChangesMapping)
	conSortingOrder := map[string]int{
		"p": 0,
		"u": 1,
		"c": 2,
		"f": 3,
	}
	slices.SortFunc(conChanges, func(a, b types.MigAction) int {
		aPrefix := getChangeTypePrefix(a.Type)
		bPrefix := getChangeTypePrefix(b.Type)

		if aPrefix == "CREATE" && bPrefix == "CREATE" {
			return conSortingOrder[a.EntityId[5]] - conSortingOrder[b.EntityId[5]]
		} else {
			return changeTypePrefixSortingCmp(a, b)
		}
	})

	changes = slices.Concat(
		dropViewChanges,
		tableChanges,
		colChanges,
		otherViewChanges,
		conChanges,
	)

	return changes
}
