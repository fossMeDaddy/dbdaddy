package main

import (
	"dbdaddy/build-scripts/utils"
)

func main() {
	outFile := utils.GetOutFilePath()
	utils.Build(outFile)
}
