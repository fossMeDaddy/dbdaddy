package main

import "github.com/fossmedaddy/dbdaddy/build-scripts/utils"

func main() {
	v := utils.GetCurrentVersion()
	utils.CreateTag(v)
}
