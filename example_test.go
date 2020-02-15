package leakless_test

import (
	"github.com/ysmood/leakless"
)

func ExampleNew() {
	// just like using the exec.Command
	cmd := leakless.New().Command("go", "version")

	_ = cmd.Run()
}
