package libUtils

import (
	"errors"
	"fmt"
	"os"
)

// returns "(true, nil)" if dir was created, else returns "(false, ...)"
func DirExistsCreate(dirName string) (bool, error) {
	// Stat the directory to check if it exists
	if _, err := os.Stat(dirName); errors.Is(err, os.ErrNotExist) {
		// Directory does not exist, so create it
		err := os.MkdirAll(dirName, 0755) // 0755 permissions allowing read and execute for others
		if err != nil {
			return false, fmt.Errorf("error creating directory: %v", err)
		}

		return true, nil
	} else if err != nil {
		// Some other error occurred when trying to access the directory
		return false, fmt.Errorf("error checking directory: %v", err)
	}

	return false, nil
}

func Exists(path string) bool {
	info, _ := os.Stat(path)
	return info != nil
}
