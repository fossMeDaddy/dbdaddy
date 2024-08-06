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

	return filepath.Join(homeDir, SelfGlobalDirName, "dbdaddy.config.json")
}

func GetLocalConfigPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(cwd, SelfGlobalDirName, "dbdaddy.config.json")
}

func GetGlobalDirPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homeDir, SelfGlobalDirName)
}
