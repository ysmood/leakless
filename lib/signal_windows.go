// +build windows

package lib

import (
	"os"
	"os/exec"
	"syscall"
)

// Signal to process
func Signal(cmd *exec.Cmd, sig os.Signal) error {
	d, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}
	p, err := d.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return err
	}
	r, _, err := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(cmd.Process.Pid))
	if r == 0 {
		return err
	}
	return nil
}
