package migrationsLib

import "dbdaddy/types"

func GetSQLFromDiffChanges(currentState, prevState *types.DbSchemaMapping, changes []types.MigAction) string {
	return ""
}
