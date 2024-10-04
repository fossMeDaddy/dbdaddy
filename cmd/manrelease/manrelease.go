package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fossmedaddy/dbdaddy/cmd/utils"
)

func main() {
	args := os.Args[1:]

	saveVersion := false
	dryRun := false

	if len(args) == 0 {
		fmt.Println("please the first argument, possible values for first arg: 'major' | 'minor' | 'patch'")
		panic("")
	}

	currentV := utils.GetCurrentVersion()

	version, vParseErr := utils.NewVersion(currentV)
	if vParseErr != nil {
		fmt.Println("unexpected error occured")
		fmt.Println(vParseErr)
		return
	}

	for _, arg := range args {
		if strings.HasPrefix(arg, "--tag=") {
			arg = strings.ReplaceAll(arg, "--tag=", "")
			arg = strings.Trim(arg, "'")
			arg = strings.Trim(arg, "\"")

			if !utils.VersionTagRegexp.MatchString(arg) {
				fmt.Println("invalid --tag value provided expected format: --tag=[0-9A-Za-z-]+")
				panic("")
			} else {
				version.Tag = arg
			}
		}

		if arg == "--save" {
			saveVersion = true
		}
		if arg == "--dry-run" {
			dryRun = true
		}
	}

	vBumpSpec := args[0]
	switch vBumpSpec {
	case "major":
		version.Major++
		version.Minor = 0
		version.Patch = 0

	case "minor":
		version.Minor++
		version.Patch = 0

	case "patch":
		version.Patch++

	default:
		panic("invalid first arg provided, possible values: 'major' | 'minor' | 'patch'")
	}

	if dryRun {
		fmt.Print(version.String())
		return
	}

	if saveVersion {
		if err := os.WriteFile("version", []byte(version.String()), 0777); err != nil {
			fmt.Println(err)
			panic("writing to version file failed")
		}
	}

	files, fileErr := filepath.Glob("bin/*")
	if fileErr != nil {
		fmt.Println(fileErr)
		panic("unexpected error occured")
	}

	ghArgs := []string{
		"release",
		"create",
	}

	ghArgs = append(ghArgs, "--draft")
	if len(version.Tag) > 0 {
		ghArgs = append(ghArgs, "--prerelease")
	} else {
		ghArgs = append(ghArgs, "--latest")
	}

	ghArgs = append(ghArgs, version.String())
	ghArgs = append(ghArgs, files...)

	ghCmd := exec.Command("gh", ghArgs...)
	ghCmd.Stdout = os.Stdout
	ghCmd.Stdin = os.Stdin
	ghCmd.Stderr = os.Stderr

	if err := ghCmd.Run(); err != nil {
		fmt.Println(err)
		panic("unexpected error occured while running 'gh' command")
	}
}
