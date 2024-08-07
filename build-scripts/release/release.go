package main

import "dbdaddy/build-scripts/utils"

func main() {
	v := utils.GetCurrentVersion()
	utils.CreateTag(v)
}
