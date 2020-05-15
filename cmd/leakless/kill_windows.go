package main

import (
	"os/exec"
	"strconv"
)

// kill process and all its children process
func killTree(pid int) error {
	return exec.Command("taskkill", "/t", "/f", "/pid", strconv.Itoa(pid)).Run()
}
