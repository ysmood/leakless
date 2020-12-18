package main

import (
	"fmt"
	"os/exec"

	"github.com/ysmood/leakless"
)

func main() {
	path := leakless.GetLeaklessBin()
	out, err := exec.Command("go", "build", "-o", path, "./cmd/leakless").CombinedOutput()
	fmt.Println(path, string(out), err)
}
