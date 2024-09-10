package main

import "github.com/fossmedaddy/dbdaddy/cmd/utils"

func main() {
	v := utils.GetCurrentVersion()
	utils.CreateTag(v)
}
