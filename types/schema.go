package types

type DbType struct {
	Schema string
	Name   string
}

type DbSequence struct {
	Schema      string
	Name        string
	DataType    string
	IncrementBy int
	MinValue    int
	MaxValue    int
	StartValue  int
	CacheSize   int
	Cycle       bool
}

type DbIndex struct {
	Schema    string
	TableName string
	Name      string

	NAttributes      int
	IsUnique         bool
	NullsNotDistinct bool
	KeyCols          []int

	Syntax string
}

type DbConstraint struct {
	// on table
	TableSchema string
	TableName   string
	ColName     string // UNUSED

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

	MultiColConstraint bool
}

type Column struct {
	TableSchema      string
	TableName        string
	Name             string
	Default          string
	Nullable         bool
	DataType         string
	CharMaxLen       int
	NumericPrecision int
	NumericScale     int
}

type Schema struct {
	Name string
}

type TableSchema struct {
	Schema        string
	Name          string
	ViewDefSyntax string // populated for only views
	Columns       []Column
	Constraints   []*DbConstraint
	Indexes       []DbIndex
}

type Table struct {
	Schema string
	Name   string

	// IMP: possible values to be put in from constants 'TableType*'
	Type string
}

// the string key is of format "schema.table"
type DbSchema struct {
	DbName    string
	Tables    map[string]*TableSchema
	Views     map[string]*TableSchema
	Types     []DbType
	Schemas   []Schema
	Sequences []DbSequence
}
