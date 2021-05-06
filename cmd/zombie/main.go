package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ysmood/leakless/pkg/utils"
)

type stamp struct {
	PID  int
	Time string
}

func main() {
	go func() {
		utils.Sleep(10)
		os.Exit(1)
	}()

	id := os.Getpid()

	for {
		now := time.Now().Format(time.RFC3339Nano)
		s := stamp{
			PID:  id,
			Time: now,
		}
		utils.E(utils.OutputFile(filepath.FromSlash("tmp/pid"), s, nil))
		utils.Sleep(0.3)
	}
}
