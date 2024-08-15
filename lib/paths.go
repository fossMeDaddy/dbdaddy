package lib

import (
	constants "dbdaddy/const"
	"dbdaddy/libUtils"
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

func FindConfigDirPath() (string, error) {
	configFile, err := FindConfigFilePath()
	if err != nil {
		return "", err
	}

	configDir, _ := path.Split(configFile)

	return configDir, nil
}

func FindTmpDirPath() (string, error) {
	configDir, err := FindConfigDirPath()
	if err != nil {
		return "", err
	}

	tmpDirPath := path.Join(configDir, constants.TmpDir)
	_, osDirErr := libUtils.DirExistsCreate(tmpDirPath)
	if osDirErr != nil {
		return "", err
	}

	return tmpDirPath, nil
}
