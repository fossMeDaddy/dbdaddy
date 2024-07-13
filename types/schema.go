package types

type Column struct {
	Name         string
	Default      string
	Nullable     bool
	DataType     string
	IsPrimaryKey bool
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
}
