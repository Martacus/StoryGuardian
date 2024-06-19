package utility

import (
	"fmt"
	"os/exec"
	"runtime"
)

func openExplorer(location string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", location)
	case "darwin":
		cmd = exec.Command("open", location)
	case "linux":
		cmd = exec.Command("xdg-open", location)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
