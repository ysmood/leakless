package main

import (
	"encoding/json"
	"net"
	"os"
	"os/exec"

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

	conn, err := net.Dial("tcp", addr)
	panicErr(err)

	cmd := exec.Command(os.Args[3], os.Args[4:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		send(conn, 0, err.Error())
		panicErr(err)
	}

	err = deathsig(cmd.Process)
	if err != nil {
		send(conn, cmd.Process.Pid, err.Error())
		panicErr(err)
	}

	send(conn, cmd.Process.Pid, "")

	go guard(conn, uid, cmd)

	err = cmd.Wait()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			os.Exit(exitErr.ExitCode())
			return
		}
		send(conn, 0, err.Error())
	}
}

func guard(conn net.Conn, uid string, cmd *exec.Cmd) {
	defer kill(cmd)

	dec := json.NewDecoder(conn)

	var msg lib.Message
	err := dec.Decode(&msg)
	if err != nil {
		return
	}
	if msg.UID != uid {
		return
	}

	_ = dec.Decode(&msg)
}

func panicErr(err error) {
	if err == nil {
		return
	}
	panic("[leakless] " + err.Error())
}

func send(conn net.Conn, pid int, errMessage string) {
	enc := json.NewEncoder(conn)
	err := enc.Encode(lib.Message{PID: pid, Error: errMessage})
	panicErr(err)
}

func kill(cmd *exec.Cmd) {
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
}
