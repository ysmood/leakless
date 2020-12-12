// +build windows

package main

import (
	"os"
	"os/exec"
	"strconv"
)

func kill(p *os.Process) {
	_ = exec.Command("taskkill", "/t", "/f", "/pid", strconv.Itoa(p.Pid)).Run()
}
