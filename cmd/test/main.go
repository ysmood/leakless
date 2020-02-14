package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ysmood/kit"
	"github.com/ysmood/leakless"
)

func main() {
	p := filepath.FromSlash
	var cmd *exec.Cmd
	if os.Args[1] == "on" {
		cmd = leakless.New().Bin(p("dist/leakless")).Command(p("dist/zombie"))
	} else {
		cmd = exec.Command(p("dist/zombie"))
	}
	kit.E(cmd.Start())
	kit.Sleep(1)
}
