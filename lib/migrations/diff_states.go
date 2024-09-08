package migrationsLib

import (
	"dbdaddy/lib/libUtils"
	"dbdaddy/types"
	"slices"

	"golang.org/x/exp/maps"
)

var (
	actionTypeSortingOrder = map[types.ActionType]int{
		types.ActionTypeDrop:   1,
		types.ActionTypeCreate: 2,
	}

	// dropping order
	entityTypeSortingOrder = map[types.EntityType]int{
		types.EntityTypeView:       1,
		types.EntityTypeConstraint: 2,
		types.EntityTypeColumn:     3,
		types.EntityTypeTable:      4,
		types.EntityTypeSequence:   5,
		types.EntityTypeSchema:     6,
	}

	conSortingOrder = map[string]int{
		"p": 1,
		"u": 2,
		"c": 3,
		"f": 4,
	}
)

// entityPtr, is a pointer!!
func getDiffKey(entityPtr interface{}, entityType types.EntityType, stateTag types.StateTag) types.DiffKey {
	diffKey := types.DiffKey{
		Entity:   types.NewEntity(entityType, entityPtr),
		StateTag: stateTag,
	}

	return diffKey
}

func diffKeyCompareFunc(a, b types.DiffKey) int {
	aSlice := slices.Concat([]string{string(a.StateTag)}, a.Entity.Id)
	bSlice := slices.Concat([]string{string(b.StateTag)}, b.Entity.Id)

	return slices.Compare(aSlice, bSlice)
}

func getDiffKeysFromStates(currentState, prevState *types.DbSchema) []types.DiffKey {
	// SCHEMA KEYS
	schemaDiffKeysConcat := []types.DiffKey{}
	for _, schema := range currentState.Schemas {
		schemaDiffKeysConcat = append(schemaDiffKeysConcat, getDiffKey(&schema, types.EntityTypeSchema, types.StateTagCS))
	}
	for _, schema := range prevState.Schemas {
		schemaDiffKeysConcat = append(schemaDiffKeysConcat, getDiffKey(&schema, types.EntityTypeSchema, types.StateTagPS))
	}

	// SEQUENCE KEYS
	seqDiffKeysConcat := []types.DiffKey{}
	for _, seq := range currentState.Sequences {
		seqDiffKeysConcat = append(seqDiffKeysConcat, getDiffKey(&seq, types.EntityTypeSequence, types.StateTagCS))
	}
	for _, seq := range prevState.Sequences {
		seqDiffKeysConcat = append(seqDiffKeysConcat, getDiffKey(&seq, types.EntityTypeSequence, types.StateTagPS))
	}

	// TABLE, VIEW, CONSTRAINT & COL KEYS
	conDiffKeysConcat := []types.DiffKey{}
	tableDiffKeysConcat := []types.DiffKey{}
	viewDiffKeysConcat := []types.DiffKey{}
	colDiffKeysConcat := []types.DiffKey{}

	currentStateMaxI := len(maps.Keys(currentState.Tables)) + len(maps.Keys(currentState.Views))
	for i, dbTableMapKey := range slices.Concat(
		maps.Keys(currentState.Tables), maps.Keys(currentState.Views),
		maps.Keys(prevState.Tables), maps.Keys(prevState.Views),
	) {
		var dbTable *types.TableSchema
		var tableStateTag types.StateTag
		var entityType types.EntityType

		if i < currentStateMaxI {
			tableStateTag = types.StateTagCS
		} else {
			tableStateTag = types.StateTagPS
		}

		if currentState.Tables[dbTableMapKey] != nil {
			dbTable = currentState.Tables[dbTableMapKey]
			entityType = types.EntityTypeTable
		} else if prevState.Tables[dbTableMapKey] != nil {
			dbTable = prevState.Tables[dbTableMapKey]
			entityType = types.EntityTypeTable
		} else if currentState.Views[dbTableMapKey] != nil {
			dbTable = currentState.Views[dbTableMapKey]
			entityType = types.EntityTypeView
		} else {
			dbTable = prevState.Views[dbTableMapKey]
			tableStateTag = types.StateTagPS
			entityType = types.EntityTypeView
		}

		if entityType == types.EntityTypeTable {
			tableDiffKeysConcat = append(tableDiffKeysConcat, getDiffKey(dbTable, entityType, tableStateTag))
		} else {
			viewDiffKeysConcat = append(viewDiffKeysConcat, getDiffKey(dbTable, entityType, tableStateTag))
			continue // DONT NEED COL TRACKING ON VIEWS
		}

		for _, col := range dbTable.Columns {
			colDiffKeysConcat = append(colDiffKeysConcat,
				getDiffKey(&col, types.EntityTypeColumn, tableStateTag),
			)
		}

		for _, con := range dbTable.Constraints {
			conDiffKeysConcat = append(conDiffKeysConcat, getDiffKey(con, types.EntityTypeConstraint, tableStateTag))
		}
	}

	diffKeys := slices.Concat(
		schemaDiffKeysConcat,
		seqDiffKeysConcat,
		tableDiffKeysConcat,
		viewDiffKeysConcat,
		colDiffKeysConcat,
		conDiffKeysConcat,
	)

	// TODO: not using binary search, use it in future
	slices.SortFunc(diffKeys, diffKeyCompareFunc)

	return diffKeys
}

func getChangesSortCompareKey(action types.MigAction) []int {
	entityTypeSorting := 1
	if action.ActionType == types.ActionTypeCreate {
		entityTypeSorting = -1
	}
	sortingKey := []int{
		actionTypeSortingOrder[action.ActionType],
		entityTypeSortingOrder[action.Entity.Type] * entityTypeSorting,
	}
	if action.Entity.Type == types.EntityTypeConstraint {
		con := action.Entity.Ptr.(*types.DbConstraint)
		sortingKey = append(sortingKey, conSortingOrder[con.Type])
	}

	return sortingKey
}

func changesSortCompareFunc(a, b types.MigAction) int {
	return slices.Compare(
		getChangesSortCompareKey(a),
		getChangesSortCompareKey(b),
	)
}

func DiffDBSchema(currentState, prevState *types.DbSchema) []types.MigAction {
	changes := []types.MigAction{}
	diffKeys := getDiffKeysFromStates(currentState, prevState)

	// action entity type bucket
	entityTypeGroupedActions := map[types.EntityType][]types.MigAction{}

	// changes mapping for avoiding creating dependent changes (e.g dropping a column from table if table is already dropped)
	schemaChangesMapping := map[string]*types.MigAction{}
	tableChangesMapping := map[string]*types.MigAction{}

	for _, diffKey := range diffKeys {
		targetDiffKey := diffKey
		if diffKey.StateTag == types.StateTagCS {
			targetDiffKey.StateTag = types.StateTagPS
		} else {
			targetDiffKey.StateTag = types.StateTagCS
		}

		_, found := slices.BinarySearchFunc(diffKeys, targetDiffKey, diffKeyCompareFunc)
		if !found {
			action := types.MigAction{
				StateTag: diffKey.StateTag,
				Entity:   diffKey.Entity,
			}
			if targetDiffKey.StateTag == types.StateTagCS {
				// could not be found in CS (ACTION: DROP)
				action.ActionType = types.ActionTypeDrop
			} else {
				// could not be found in PS (ACTION: CREATE)
				action.ActionType = types.ActionTypeCreate
			}

			// changes book-keeping
			switch action.Entity.Type {
			case types.EntityTypeSchema:
				schema := action.Entity.Ptr.(*types.Schema)
				schemaChangesMapping[schema.Name] = &action

			case types.EntityTypeTable:
				table := action.Entity.Ptr.(*types.TableSchema)
				tableChangesMapping[libUtils.GetTableId(table.Schema, table.Name)] = &action
			}

			entityTypeGroupedActions[action.Entity.Type] = append(
				entityTypeGroupedActions[action.Entity.Type],
				action,
			)
		} else {
			// MODIFICATION GOES HERE
		}
	}

	for entityType, actions := range entityTypeGroupedActions {
		newActions := []types.MigAction{}

		switch entityType {
		case types.EntityTypeTable:
			for _, action := range actions {
				tableSchema := action.Entity.Ptr.(*types.TableSchema)
				schemaChange := schemaChangesMapping[libUtils.GetTableId(tableSchema.Schema, tableSchema.Name)]
				if schemaChange != nil && schemaChange.ActionType == types.ActionTypeDrop {
					continue
				}

				newActions = append(newActions, action)
			}
		case types.EntityTypeColumn:
			for _, action := range actions {
				col := action.Entity.Ptr.(*types.Column)
				tableChange := tableChangesMapping[libUtils.GetTableId(col.TableSchema, col.TableName)]
				if tableChange != nil {
					continue
				}

				newActions = append(newActions, action)
			}
		case types.EntityTypeConstraint:
			for _, action := range actions {
				con := action.Entity.Ptr.(*types.DbConstraint)
				tableChange := tableChangesMapping[libUtils.GetTableId(con.TableSchema, con.TableName)]
				if tableChange != nil && tableChange.ActionType == types.ActionTypeDrop {
					continue
				}

				newActions = append(newActions, action)
			}
		default:
			newActions = actions
		}

		changes = slices.Concat(changes, newActions)
	}

	slices.SortFunc(changes, changesSortCompareFunc)

	return changes
}
