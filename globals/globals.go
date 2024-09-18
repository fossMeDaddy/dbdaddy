package globals

import (
	"database/sql"

	"github.com/fossmedaddy/dbdaddy/types"
)

// global DB related vars (please dont fuck with them)
var (
	// contains full version string e.g. "vX.Y.Z-ABC"
	Version string

	DB *sql.DB

	CurrentConnConfig types.ConnConfig
)
