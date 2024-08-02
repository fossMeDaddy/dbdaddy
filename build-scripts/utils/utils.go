package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func GetCurrentVersion() string {
	r, err := os.ReadFile("version")
	if err != nil {
		panic("couldn't read file 'version'")
	}

	return string(r)
}

func GetOutFilePath() string {
	v := GetCurrentVersion()
	fileFriendlyV := strings.ReplaceAll(v, ".", "-")

	outFile := path.Join("bin", fmt.Sprintf("dbdaddy-%s", fileFriendlyV))

	return outFile
}

func Build(outFile string) {
	cmd := exec.Command("go", "build", "-o", outFile, ".")
	cmdErr := cmd.Run()
	if cmdErr != nil {
		fmt.Println(cmdErr)
		panic("error occured while running 'go build'")
	}
}

func Release(version string, file string) {
	ghCmd := exec.Command("gh", "release", "create", version, file)
	ghCmd.Stdout = os.Stdout
	ghCmd.Stdin = os.Stdin
	ghCmd.Stderr = os.Stderr

	ghCmdErr := ghCmd.Run()
	if ghCmdErr != nil {
		panic("error occured while running 'gh release command'")
	}
}
