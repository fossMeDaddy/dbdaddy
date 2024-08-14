package types

type DbType struct {
	Schema string
	Name   string
}

type Column struct {
	Name               string
	Default            string
	Nullable           bool
	DataType           string
	CharMaxLen         int
	NumericPrecision   int
	NumericScale       int
	IsPrimaryKey       bool
	IsRelation         bool
	ForeignTableSchema string
	ForeignTableName   string
	ForeignColumnName  string
}

type TableSchema struct {
	Db      string
	Schema  string
	Name    string
	Columns []Column
}

type Table struct {
	Name   string
	Schema string
	Type   string // IMP: possible values to be put in from constants 'TableType*'
}

// the string key is of format "schema.table"
type DbSchema struct {
	Tables map[string]*TableSchema
	Types  []DbType
}
