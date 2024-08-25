package migrationsLib

import (
	"dbdaddy/types"
	"fmt"
	"reflect"
	"slices"

	"golang.org/x/exp/maps"
)

// WILL PANIC IF "entity" is a pointer!!!
func getDiffKey(entity interface{}, entityType types.EntityType, stateTag types.StateTag) types.DiffKey {
	diffKey := types.DiffKey{
		EntityType: entityType,
		StateTag:   stateTag,
	}

	entityReflect := reflect.ValueOf(entity)
	for i := 0; i < entityReflect.NumField(); i++ {
		diffKey.EntityId = append(diffKey.EntityId, entityReflect.Field(i).String())
	}

	return diffKey
}

func diffKeyCompareFunc(a, b types.DiffKey) int {
	aSlice := slices.Concat([]string{fmt.Sprint(len(a.DepIds)), string(a.StateTag)}, a.EntityId)
	bSlice := slices.Concat([]string{fmt.Sprint(len(b.DepIds)), string(b.StateTag)}, b.EntityId)

	return slices.Compare(aSlice, bSlice)
}

func getKeysFromStates(currentState, prevState *types.DbSchema) []types.DiffKey {
	// TYPE KEYS
	typeDiffKeysConcat := []types.DiffKey{}
	for _, dbType := range currentState.Types {
		typeDiffKeysConcat = append(typeDiffKeysConcat, getDiffKey(dbType, types.EntityTypeType, types.StateTagCS))
	}
	for _, dbType := range prevState.Types {
		typeDiffKeysConcat = append(typeDiffKeysConcat, getDiffKey(dbType, types.EntityTypeType, types.StateTagPS))
	}

	// TABLE, VIEW, CONSTRAINT & COL KEYS
	conDiffKeysConcat := []types.DiffKey{}
	tableDiffKeysConcat := []types.DiffKey{}
	viewDiffKeysConcat := []types.DiffKey{}
	colDiffKeysConcat := []types.DiffKey{}
	for _, dbTableMapKey := range slices.Concat(
		maps.Keys(currentState.Tables), maps.Keys(prevState.Tables),
		maps.Keys(currentState.Views), maps.Keys(prevState.Views),
	) {
		var dbTable *types.TableSchema
		var tableStateTag types.StateTag
		var entityType types.EntityType
		if currentState.Tables[dbTableMapKey] != nil {
			// CS TABLE
			dbTable = currentState.Tables[dbTableMapKey]
			tableStateTag = types.StateTagCS
			entityType = types.EntityTypeTable
		} else if prevState.Tables[dbTableMapKey] != nil {
			// PS TABLE
			dbTable = prevState.Tables[dbTableMapKey]
			tableStateTag = types.StateTagPS
			entityType = types.EntityTypeTable
		} else if currentState.Views[dbTableMapKey] != nil {
			// CS VIEW
			dbTable = currentState.Views[dbTableMapKey]
			tableStateTag = types.StateTagCS
			entityType = types.EntityTypeView
		} else {
			// PS VIEW
			dbTable = prevState.Views[dbTableMapKey]
			tableStateTag = types.StateTagPS
			entityType = types.EntityTypeView
		}

		if entityType == types.EntityTypeTable {
			tableDiffKeysConcat = append(tableDiffKeysConcat, getDiffKey(*dbTable, entityType, tableStateTag))
		} else {
			viewDiffKeysConcat = append(viewDiffKeysConcat, getDiffKey(*dbTable, entityType, tableStateTag))
			continue // DONT NEED COL TRACKING ON VIEWS
		}

		for _, col := range dbTable.Columns {
			colDiffKeysConcat = append(colDiffKeysConcat,
				getDiffKey(struct {
					Schemaname string
					Tablename  string
					types.Column
				}{
					Schemaname: dbTable.Schema,
					Tablename:  dbTable.Name,
					Column:     col,
				}, types.EntityTypeColumn, tableStateTag),
			)
		}

		for _, con := range dbTable.Constraints {
			conDiffKeysConcat = append(conDiffKeysConcat, getDiffKey(*con, types.EntityTypeConstraint, tableStateTag))
		}
	}

	diffKeys := slices.Concat(
		typeDiffKeysConcat,
		tableDiffKeysConcat,
		viewDiffKeysConcat,
		colDiffKeysConcat,
		conDiffKeysConcat,
	)

	slices.SortFunc(diffKeys, diffKeyCompareFunc)

	return diffKeys
}

func DiffDBSchema(currentState, prevState *types.DbSchema) []types.MigAction {
	diffKeys := getKeysFromStates(currentState, prevState)
	fmt.Println(diffKeys)

	return []types.MigAction{}
}
