package libUtils

import (
	"fmt"
	"os"
	"os/exec"
)

func OpenFileInEditor(filePath string) error {
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
			fmt.Println("Holy shit bro?! wtf are you using for an OS? no vim, no nano, where tf am i?!")
			fmt.Println("nano command gave the error:\n" + nanoErr.Error())
			fmt.Println("vim command gave the error:\n" + vimErr.Error())
			return fmt.Errorf("vim error: %s%snano error: %s", vimErr.Error(), fmt.Sprintln(), nanoErr.Error())
		}
	}

	return nil
}
