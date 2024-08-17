package errs

import (
	"fmt"
)

var ErrUnsupportedDriver = fmt.Errorf("driver not supported")
var ErrDbAlreadyExists = fmt.Errorf("database already exists")

func UnsupportedDriverMsg(currentDriver string, supportedDrivers []string) string {
	return fmt.Sprintf("'%s' driver is not supported, as of now, the supported drivers are: %v", currentDriver, supportedDrivers)
}

func DbAlreadyExistsMsg(dbname string) string {
	return fmt.Sprintf("database: '%s' already exists", dbname)
}