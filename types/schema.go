package types

type DbType struct {
	Schema string
	Name   string
}

type DbConstraint struct {
	ConName          string
	ConSchema        string
	Type             string
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

	// OnlyFrontend (only for use by frontend-loving buttplug furries)
	IsPrimaryKey bool
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
	Type   string // IMP: possible values to be put in from constants 'TableType*'
}

// the string key is of format "schema.table"
type DbSchema struct {
	DbName string
	Tables map[string]*TableSchema
	Types  []DbType
}
