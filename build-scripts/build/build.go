package main

import (
	"dbdaddy/build-scripts/utils"
)

func main() {
	utils.Build("linux", "arm64")
	utils.Build("linux", "amd64")
	utils.Build("linux", "386")

	utils.Build("darwin", "arm64")
	utils.Build("darwin", "amd64")

	utils.Build("windows", "amd64")
	utils.Build("windows", "386")

	utils.Build("freebsd", "amd64")
}
