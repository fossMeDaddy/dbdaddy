package lib

import (
	constants "dbdaddy/const"
	"fmt"
	"os"
	"path"
)

func FindConfigFilePath() (string, error) {
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

func FindConfigDir() (string, error) {
	configFile, err := FindConfigFilePath()
	if err != nil {
		return "", err
	}

	configDir, _ := path.Split(configFile)

	return configDir, nil
}

func GetAbsolutePathFor(relativePath string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(cwd, relativePath), err
}

func GetTmpDirPath() (string, error) {
	configDir := constants.GetGlobalDirPath()

	tmpDirPath := path.Join(configDir, constants.TmpDir)
	err := DirExistsCreate(tmpDirPath)
	if err != nil {
		return "", err
	}

	return tmpDirPath, nil
}
