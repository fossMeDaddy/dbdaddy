package main

import "dbdaddy/build-scripts/utils"

func main() {
	v := utils.GetCurrentVersion()
	outFile := utils.GetOutFilePath()
	utils.Release(v, outFile)
}
