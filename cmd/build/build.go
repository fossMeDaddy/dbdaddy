package main

import (
	"os"

	"github.com/fossmedaddy/dbdaddy/cmd/utils"
)

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		panic("invalid number of args given, expected: <=1")
	}

	currentV := utils.GetCurrentVersion()
	if len(args) == 1 {
		if v, err := utils.NewVersion(args[0]); err != nil {
			panic("error occured in parsing version argument. " + err.Error())
		} else {
			currentV = v.String()
		}
	}

	utils.BuildAllTargets(currentV)
}
