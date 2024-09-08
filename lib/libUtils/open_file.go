package libUtils

import (
	"fmt"
	"os"
	"os/exec"
)

func OpenFileInEditor(filePath string) {
	vimOsCmd := exec.Command("vim", filePath)
	vimOsCmd.Stdin = os.Stdin
	vimOsCmd.Stdout = os.Stdout
	vimOsCmd.Stderr = os.Stderr

	vimErr := vimOsCmd.Run()
	if vimErr != nil {
		fmt.Println("Failed to open vim, trying nano...")

		nanoOsCmd := exec.Command("nano", filePath)
		nanoOsCmd.Stdin = os.Stdin
		nanoOsCmd.Stdout = os.Stdout
		nanoOsCmd.Stderr = os.Stderr

		nanoErr := nanoOsCmd.Run()
		if nanoErr != nil {
			fmt.Println("Holy shit bro?! wtf are you using for an OS? no vim, no nano, is this what, an alpine docker container?")
			fmt.Println("nano command gave the error:\n" + nanoErr.Error())
			fmt.Println("vim command gave the error:\n" + vimErr.Error())
			os.Exit(1)
		}
	}
}
