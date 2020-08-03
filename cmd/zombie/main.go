package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ysmood/leakless/lib"
)

type stamp struct {
	PID  int
	Time string
}

func main() {
	go func() {
		lib.Sleep(10)
		os.Exit(1)
	}()

	id := os.Getpid()

	for {
		now := time.Now().Format(time.RFC3339Nano)
		s := stamp{
			PID:  id,
			Time: now,
		}
		lib.E(lib.OutputFile(filepath.FromSlash("tmp/pid"), s, nil))
		lib.Sleep(0.3)
	}
}
