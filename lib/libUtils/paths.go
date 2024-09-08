package libUtils

import (
	"fmt"
	"os"
	"path"

	"github.com/fossmedaddy/dbdaddy/constants"
)

func GetAbsolutePathFor(relativePath string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(cwd, relativePath), err
}

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

	configDir := path.Dir(configFile)

	return configDir, nil
}

func FindTmpDirPath() (string, error) {
	configDir, err := FindConfigDirPath()
	if err != nil {
		return "", err
	}

	tmpDirPath := path.Join(configDir, constants.TmpDir)
	_, osDirErr := EnsureDirExists(tmpDirPath)
	if osDirErr != nil {
		return "", err
	}

	return tmpDirPath, nil
}
