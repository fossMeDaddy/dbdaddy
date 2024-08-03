package main

import "dbdaddy/build-scripts/utils"

func main() {
	v := utils.GetCurrentVersion()
	utils.Release(v, "bin")
}
