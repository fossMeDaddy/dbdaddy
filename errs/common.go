package errs

import (
	"fmt"
)

var ErrUnsupportedDriver = fmt.Errorf("driver not supported")

func UnsupportedDriverMsg(currentDriver string, supportedDrivers []string) error {
	return fmt.Errorf("'%s' driver is not supported, as of now, the supported drivers are: %v", currentDriver, supportedDrivers)
}
