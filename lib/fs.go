package lib

import (
	"fmt"
	"os"
)

func DirExistsCreate(dirName string) error {
	// Stat the directory to check if it exists
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		// Directory does not exist, so create it
		err := os.MkdirAll(dirName, 0755) // 0755 permissions allowing read and execute for others
		if err != nil {
			return fmt.Errorf("error creating directory: %v", err)
		}
	} else if err != nil {
		// Some other error occurred when trying to access the directory
		return fmt.Errorf("error checking directory: %v", err)
	}
	return nil
}

func FileExists(path string) bool {
	info, _ := os.Stat(path)
	return info != nil
}
