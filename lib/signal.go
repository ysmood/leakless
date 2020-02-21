// +build !windows

package lib

import (
	"os"
	"os/exec"
)

// Signal to process
func Signal(cmd *exec.Cmd, sig os.Signal) error {
	return cmd.Process.Signal(sig)
}
