// +build !windows

package main

import (
	"os"
)

func kill(p *os.Process) {
	_ = p.Kill()
}
