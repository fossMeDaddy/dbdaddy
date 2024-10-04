package main

import (
	"github.com/fossmedaddy/dbdaddy/cmd/utils"
)

func main() {
	currentV := utils.GetCurrentVersion()
	utils.BuildAllTargets(currentV)
}
