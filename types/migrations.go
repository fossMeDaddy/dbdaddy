package types

type MigAction struct {
	// actions defined in constants
	Type string

	// if table: `["dbname", "schemaname", "tablename"]`
	// else if column: `["dbname", "schemaname", "tablename", "columnname"]`
	EntityId []string
}
