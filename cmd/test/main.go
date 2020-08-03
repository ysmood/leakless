package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ysmood/leakless"
	"github.com/ysmood/leakless/lib"
)

func main() {
	p := filepath.FromSlash
	var cmd *exec.Cmd
	if os.Args[1] == "on" {
		l := leakless.New()
		cmd = l.Command(p("dist/zombie"))
		go func() {
			pid := <-l.Pid()
			lib.E(lib.OutputFile(filepath.FromSlash("tmp/sub-pid"), lib.MustToJSON(pid), nil))
		}()
	} else {
		cmd = exec.Command(p("dist/zombie"))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	lib.E(cmd.Start())
	lib.E(cmd.Process.Release())
	lib.Sleep(2)
}
