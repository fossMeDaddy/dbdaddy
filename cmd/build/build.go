package main

import (
	"dbdaddy/cmd/utils"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(8)

	go (func() {
		defer wg.Done()
		utils.Build("linux", "arm64")
	})()
	go (func() {
		defer wg.Done()
		utils.Build("linux", "amd64")
	})()
	go (func() {
		defer wg.Done()
		utils.Build("linux", "386")
	})()

	go (func() {
		defer wg.Done()
		utils.Build("darwin", "arm64")
	})()
	go (func() {
		defer wg.Done()
		utils.Build("darwin", "amd64")
	})()

	go (func() {
		defer wg.Done()
		utils.Build("windows", "amd64")
	})()
	go (func() {
		defer wg.Done()
		utils.Build("windows", "386")
	})()

	go (func() {
		defer wg.Done()
		utils.Build("freebsd", "amd64")
	})()

	wg.Wait()
}
