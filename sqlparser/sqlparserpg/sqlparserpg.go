package sqlparserpg

import (
	pgParser "github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/sql/sem/tree"
	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/types"
)

// checks if the given statement is for a table, if so, returns (<schema>, <tablename>, error)
func getTableFromAST(ast tree.Statement) (types.Table, bool) {
	table := types.Table{}

	switch ast := ast.(type) {
	case *tree.CreateTable:
		table.Schema = ast.Table.Schema()
		table.Name = ast.Table.Table()
		table.Type = constants.TableTypeBaseTable

	case *tree.AlterTable:
		astTableName := ast.Table.ToTableName()

		table.Schema = astTableName.SchemaName.String()
		table.Name = astTableName.Table()
		table.Type = constants.TableTypeBaseTable

	case *tree.CreateView:
		table.Schema = ast.Name.Schema()
		table.Name = ast.Name.Table()
		table.Type = constants.TableTypeView

	default:
		return table, false
	}

	return table, true
}

func constructTableSchemaFromAST(ast tree.Statement, dbSchema *types.DbSchema) {
	table, _ := getTableFromAST(ast)

	tableSchema := types.TableSchema{
		Schema: table.Schema,
		Name:   table.Name,
	}

	switch ast := ast.(type) {
	case *tree.CreateTable:
	}
}

func GetDbSchemaFromSQL(sqlStr string) (types.DbSchema, error) {
	dbSchema := types.DbSchema{}

	stmts, err := pgParser.Parse(sqlStr)
	if err != nil {
		return dbSchema, err
	}

	for _, stmt := range stmts {

	}

	// sqlStmts := libUtils.GetSQLStmts(sqlStr)
	// for _, sqlStmt := range sqlStmts {
	// 	fmt.Println(sqlStmt)
	// }

	return dbSchema, nil
}
