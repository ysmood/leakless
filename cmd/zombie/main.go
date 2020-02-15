package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ysmood/kit"
)

type stamp struct {
	PID  int
	Time string
}

func main() {
	go func() {
		kit.Sleep(10)
		os.Exit(1)
	}()

	id := os.Getpid()

	for {
		now := time.Now().Format(time.RFC3339Nano)
		s := stamp{
			PID:  id,
			Time: now,
		}
		kit.E(kit.OutputFile(filepath.FromSlash("tmp/pid"), s, nil))
		kit.Sleep(0.3)
	}
}
