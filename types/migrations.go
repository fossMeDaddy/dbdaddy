package types

type MigAction struct {
	// actions defined in constants
	Type string

	/*
		if table: `["CS | PS", "dbname", "schemaname", "tablename"]`

		else if column: `["CS | PS", "dbname", "schemaname", "tablename", "columnname"]`

		else if type: ["CS | PS", "dbname", "schemaname", "typename"]

		else if constraint: `["CS | PS", "dbname", "schemaname", "tablename", "constraint-name", "constraint-type"]`

		array makes things sortable...
	*/
	EntityId []string
}
