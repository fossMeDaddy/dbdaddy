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

func GetOutDir() string {
	return path.Join("bin")
}

func GetOutFilePath(goos string, goarch string) string {
	outFile := path.Join("bin", fmt.Sprintf("dbdaddy-%s-%s", goos, goarch))

	return outFile
}

func Build(goos string, goarch string) {
	fmt.Println("Starting build for", goos, goarch)

	GOOS := strings.ToLower(goos)
	GOARCH := strings.ToLower(goarch)

	outFile := path.Join("bin", fmt.Sprintf("dbdaddy-%s-%s", GOOS, GOARCH))

	cmd := exec.Command("go", "build", "-o", outFile, ".")
	cmd.Env = append(os.Environ(), "GOOS="+GOOS, "GOARCH="+GOARCH)
	cmdErr := cmd.Run()
	if cmdErr != nil {
		fmt.Println("ERR:", cmdErr)
		panic("error occured while running 'go build'")
	}

	fmt.Println("Built binary for", goos, goarch)
	fmt.Println()
}

func Release(version string) {
	binDirEntry, err := os.ReadDir(GetOutDir())
	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprint("error occured while reading ", GetOutDir()))
	}

	binFiles := []string{}
	for _, dirEntry := range binDirEntry {
		binFiles = append(binFiles, path.Join(GetOutDir(), dirEntry.Name()))
	}

	// git tag vX.Y.Z
	gitTagCmd := exec.Command("git", "tag", version)
	tagErr := gitTagCmd.Run()
	if tagErr != nil {
		fmt.Println("error creating tag in git", tagErr)
	}

	// git push origin vX.Y.Z
	gitPushCmd := exec.Command("git", "push", "origin", version)
	pushErr := gitPushCmd.Run()
	if pushErr != nil {
		fmt.Println("error pushing to tag branch in git", pushErr)
	}

	args := []string{"release", "create", version, "--generate-notes"}
	args = append(args, binFiles...)

	ghCmd := exec.Command("gh", args...)
	ghCmd.Stdout = os.Stdout
	ghCmd.Stdin = os.Stdin
	ghCmd.Stderr = os.Stderr

	ghCmdErr := ghCmd.Run()
	if ghCmdErr != nil {
		fmt.Println(ghCmdErr)
		panic("error occured while running 'gh release command'")
	}
}
