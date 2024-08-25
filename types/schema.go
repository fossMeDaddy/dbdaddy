package types

type DbType struct {
	Schema string
	Name   string
}

type DbConstraint struct {
	// on table
	TableSchema string
	TableName   string
	ColName     string

	ConName          string
	ConSchema        string
	Type             string // can be 'p', 'f', 'c' or 'u'
	UpdateActionType string
	DeleteActionType string

	// foreign table
	FTableSchema string
	FTableName   string
	FColName     string

	Syntax string
}

type Column struct {
	Name             string
	Default          string
	Nullable         bool
	DataType         string
	CharMaxLen       int
	NumericPrecision int
	NumericScale     int
}

type TableSchema struct {
	Schema      string
	Name        string
	DefSyntax   string // populated for only views
	Columns     []Column
	Constraints []*DbConstraint
}

type Table struct {
	Schema string
	Name   string

	// IMP: possible values to be put in from constants 'TableType*'
	Type string
}

// the string key is of format "schema.table"
type DbSchema struct {
	DbName string
	Tables map[string]*TableSchema
	Views  map[string]*TableSchema
	Types  []DbType
}
