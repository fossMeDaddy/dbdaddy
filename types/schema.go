package types

type DbType struct {
	Schema string
	Name   string
}

type DbConstraint struct {
	ConName          string
	ConSchema        string
	Type             string // can be 'p', 'f', 'c' or 'u'
	UpdateActionType string
	DeleteActionType string
	Syntax           string

	// on table
	TableSchema string
	TableName   string
	ColName     string

	// foreign table
	FTableSchema string
	FTableName   string
	FColName     string
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
	Columns     []Column
	Constraints []*DbConstraint
}

type Table struct {
	Name   string
	Schema string

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
