package constants

import (
	"os"
	"path/filepath"
)

func GetGlobalConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homeDir, SelfConfigDirName, SelfConfigFileName)
}

func GetLocalConfigPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(cwd, SelfConfigDirName, SelfConfigFileName)
}

func GetGlobalDirPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homeDir, SelfConfigDirName)
}
