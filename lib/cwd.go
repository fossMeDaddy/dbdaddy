package lib

import (
	constants "dbdaddy/const"
	"fmt"
	"os"
)

func FindCwdDBCreds() {
	// will only run when initializing

	// dirEntries, err := os.ReadDir(".")
	// if err != nil {
	// 	panic("Unexpected error occured!\n" + err.Error())
	// }

	// for _, dirEntry := range dirEntries {
	// 	fmt.Println(dirEntry.Name())
	// }
}

func GetConfigFilePath() (string, error) {
	_, cwdFileErr := os.Stat(constants.GetLocalConfigPath())
	_, globalFileErr := os.Stat(constants.GetGlobalConfigPath())

	if os.IsNotExist(cwdFileErr) && os.IsNotExist(globalFileErr) {
		return "", fmt.Errorf(globalFileErr.Error() + "\n" + cwdFileErr.Error())
	}

	if os.IsNotExist(cwdFileErr) {
		return constants.GetGlobalConfigPath(), nil
	} else {
		return constants.GetLocalConfigPath(), nil
	}
}
