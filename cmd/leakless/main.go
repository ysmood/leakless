package main

import (
	"net"
	"os"
	"os/exec"
	"strconv"

	"github.com/ysmood/byframe"
	"github.com/ysmood/leakless/lib"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		_, _ = os.Stdout.WriteString(lib.Version + "\n")
		return
	}

	if len(os.Args) < 4 {
		panic("wrong args, usage: leakless <uid> <tcp-addr> <cmd> [args...]")
	}

	uid := os.Args[1]
	addr := os.Args[2]

	cmd := exec.Command(os.Args[3], os.Args[4:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		panic("[leakless] " + err.Error())
	}

	go guard(uid, addr, cmd.Process.Pid, cmd)

	err = cmd.Wait()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			os.Exit(exitErr.ExitCode())
			return
		}
		panic("[leakless] exec error: " + err.Error())
	}
}

func guard(uid, addr string, pid int, cmd *exec.Cmd) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		kill(cmd)
	}

	s := byframe.NewScanner(conn).Limit(100)
	s.Scan()
	if string(s.Frame()) != uid {
		kill(cmd)
	}

	_, err = conn.Write(byframe.Encode([]byte(strconv.Itoa(pid))))
	if err != nil {
		kill(cmd)
	}

	s.Scan()

	kill(cmd)
}

func kill(cmd *exec.Cmd) {
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
}
