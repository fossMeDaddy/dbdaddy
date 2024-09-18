package db_int

import (
	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db/pg"
	"github.com/fossmedaddy/dbdaddy/errs"
	"github.com/fossmedaddy/dbdaddy/types"
)

func DumpDb(outputFilePath string, connConfig types.ConnConfig, onlySchema bool) error {
	switch connConfig.Driver {
	case constants.DbDriverPostgres:
		return pg.DumpDb(outputFilePath, connConfig, onlySchema)
	default:
		return errs.ErrUnsupportedDriver
	}
}

func RestoreDb(connConfig types.ConnConfig, dumpFilePath string, override bool) error {
	switch connConfig.Driver {
	case constants.DbDriverPostgres:
		return pg.RestoreDb(connConfig, dumpFilePath, override)
	default:
		return errs.ErrUnsupportedDriver
	}
}
