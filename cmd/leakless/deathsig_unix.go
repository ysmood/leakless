// +build !windows

package main

import "os"

func deathsig(p *os.Process) (err error) {
	// do nothing on *nix for now
	return
}
