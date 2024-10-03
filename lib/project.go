package lib

import (
	"os"
	"path"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db/db_int"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/fossmedaddy/dbdaddy/lib/migrationsLib"
	"github.com/fossmedaddy/dbdaddy/types"
)

func LoadDbSchemaIntoSchemaDir(schemaDirPath string) error {
	dbSchema, schemaErr := db_int.GetDbSchema()
	if schemaErr != nil {
		return schemaErr
	}

	changes := migrationsLib.DiffDBSchema(dbSchema, &types.DbSchema{})
	schemaSql, sqlErr := migrationsLib.GetSQLFromDiffChanges(changes)
	if sqlErr != nil {
		return sqlErr
	}

	if err := os.RemoveAll(schemaDirPath); err != nil {
		return err
	}

	if err := os.Mkdir(schemaDirPath, 0777); err != nil {
		return err
	}

	if err := os.WriteFile(path.Join(schemaDirPath, "schema.sql"), []byte(schemaSql), 0777); err != nil {
		return err
	}

	return nil
}

func EnsureProjectDirsExist() error {
	_, migErr := libUtils.EnsureDirExists(path.Join(constants.MigDirName))
	_, schemaErr := libUtils.EnsureDirExists(path.Join(constants.SchemaDirName))
	_, scriptsErr := libUtils.EnsureDirExists(path.Join(constants.ScriptsDirName))

	if migErr != nil {
		return migErr
	}
	if schemaErr != nil {
		return schemaErr
	}
	if scriptsErr != nil {
		return scriptsErr
	}

	return nil
}
