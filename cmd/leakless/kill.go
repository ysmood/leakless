package main

import (
	"os"
	"syscall"
)

// kill process and all its children process
func killTree(pid int) error {
	group, _ := os.FindProcess(-1 * pid)

	return group.Signal(syscall.SIGTERM)
}
