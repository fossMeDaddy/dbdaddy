package globals

import "database/sql"

// global DB related vars (please dont fuck with them)
var (
	// contains full version string e.g. "vX.Y.Z-ABC"
	Version string

	DB *sql.DB

	ConnDbName string
)
