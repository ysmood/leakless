package main

import (
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/ysmood/byframe"
	"github.com/ysmood/leakless/lib"
)

func main() {
	if len(os.Args) < 4 {
		panic("wrong args, usage: leakless <uid> <tcp-addr> <cmd> [args...]")
	}

	if os.Args[1] == "--version" {
		_, _ = os.Stdout.WriteString(lib.Version + "\n")
		return
	}

	uid := os.Args[1]
	addr := os.Args[2]

	cmd := exec.Command(os.Args[3], os.Args[4:]...)

	go guard(uid, addr, cmd)
	go handleSignal(cmd)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			os.Exit(exitErr.ExitCode())
			return
		}
		panic("[leakless] exec error: " + err.Error())
	}
}

func guard(uid, addr string, cmd *exec.Cmd) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		kill(cmd)
	}

	s := byframe.NewScanner(conn).Limit(100)
	s.Scan()
	if string(s.Frame()) != uid {
		kill(cmd)
	}

	buf := make([]byte, 1)
	_, _ = conn.Read(buf)

	kill(cmd)
}

func handleSignal(cmd *exec.Cmd) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	for sig := range c {
		err := lib.Signal(cmd, sig)
		if err != nil {
			kill(cmd)
		}
	}
}

func kill(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
}
