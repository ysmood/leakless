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
		l := leakless.New().Bin(p("dist/leakless"))
		cmd = l.Command(p("dist/zombie"))
		go func() {
			pid := <-l.Pid()
			kit.E(kit.OutputFile(filepath.FromSlash("tmp/sub-pid"), kit.MustToJSON(pid), nil))
		}()
	} else {
		cmd = exec.Command(p("dist/zombie"))
	}
	kit.E(cmd.Start())
	kit.E(cmd.Process.Release())
	kit.Sleep(2)
}
