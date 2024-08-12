package types

type Column struct {
	Name               string
	Default            string
	Nullable           bool
	DataType           string
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

type DbSchema []TableSchema

// the string key is of format "schema.table"
type DbSchemaMapping map[string]*TableSchema
