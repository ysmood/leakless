package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ysmood/leakless"
	"github.com/ysmood/leakless/pkg/utils"
)

func main() {
	p := filepath.FromSlash
	var cmd *exec.Cmd
	if os.Args[1] == "on" {
		l := leakless.New()
		cmd = l.Command(p("dist/zombie"))
		go func() {
			pid := <-l.Pid()
			utils.E(utils.OutputFile(filepath.FromSlash("tmp/sub-pid"), utils.MustToJSON(pid), nil))
		}()
	} else {
		cmd = exec.Command(p("dist/zombie"))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	utils.E(cmd.Start())
	utils.E(cmd.Process.Release())
	utils.Sleep(2)
}
