package types

type MigAction struct {
	// actions defined in constants
	Type string

	/*
		if table: `["CS | PS", "tabletype", "dbname", "schemaname", "tablename"]`

		else if column: `["CS | PS", "tabletype", "dbname", "schemaname", "tablename", "columnname"]`

		else if type: ["CS | PS", "dbname", "schemaname", "typename"]

		else if constraint: `["CS | PS", "dbname", "tabletype", "schemaname", "tablename", "constraint-name", "constraint-type", "constraint-syntax"]`

		array makes things sortable...
	*/
	EntityId []string
}
