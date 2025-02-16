package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)

func GetCurrentVersion() string {
	r, err := os.ReadFile("version")
	if err != nil {
		panic("couldn't read file 'version'")
	}

	v := string(r)
	v = strings.ToLower(v)
	v = strings.Trim(v, " ")
	v = strings.Trim(v, fmt.Sprintln())

	return v
}

func GetOutDir() string {
	return path.Join("bin")
}

func GetOutFilePath(goos string, goarch string) string {
	outFile := path.Join("bin", fmt.Sprintf("dbdaddy-%s-%s", goos, goarch))

	return outFile
}

func Build(goos string, goarch string, version string) {
	fmt.Println("Starting build for", goos, goarch)

	GOOS := strings.ToLower(goos)
	GOARCH := strings.ToLower(goarch)

	outFile := path.Join("bin", fmt.Sprintf("dbdaddy-%s-%s", GOOS, GOARCH))

	cmd := exec.Command(
		"go",
		"build",
		"-ldflags",
		fmt.Sprintf("-X 'github.com/fossmedaddy/dbdaddy/globals.Version=%s'", version),
		"-o",
		outFile,
		"cmd/main/main.go",
	)
	cmd.Env = append(os.Environ(), "GOOS="+GOOS, "GOARCH="+GOARCH)
	cmd.Stdout = os.Stdout
	cmdErr := cmd.Run()
	if cmdErr != nil {
		fmt.Println("ERR:", cmdErr)
		panic("error occured while running 'go build'")
	}

	fmt.Println("Built binary for", goos, goarch)
	fmt.Println()
}

func BuildAllTargets(version string) {
	var wg sync.WaitGroup

	wg.Add(8)

	go (func() {
		defer wg.Done()
		Build("linux", "arm64", version)
	})()
	go (func() {
		defer wg.Done()
		Build("linux", "amd64", version)
	})()
	go (func() {
		defer wg.Done()
		Build("linux", "386", version)
	})()

	go (func() {
		defer wg.Done()
		Build("darwin", "arm64", version)
	})()
	go (func() {
		defer wg.Done()
		Build("darwin", "amd64", version)
	})()

	go (func() {
		defer wg.Done()
		Build("windows", "amd64", version)
	})()
	go (func() {
		defer wg.Done()
		Build("windows", "386", version)
	})()

	go (func() {
		defer wg.Done()
		Build("freebsd", "amd64", version)
	})()

	wg.Wait()
}

func CreateTag(version string) {
	// create git tag
	tagCmd := exec.Command("git", "tag", version)
	tagErr := tagCmd.Run()
	if tagErr != nil {
		fmt.Println("ERROR creating new tag:", tagErr)
		return
	}

	// push git tag
	pushCmd := exec.Command("git", "push", "origin", version)
	pushErr := pushCmd.Run()
	if pushErr != nil {
		fmt.Println("ERROR pushing tag:", pushErr)
		return
	}

	// after pushing tag, GH action will build & publish the release for the tag
}
