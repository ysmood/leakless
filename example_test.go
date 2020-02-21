package leakless_test

import (
	"fmt"
	"time"

	"github.com/ysmood/leakless"
)

func ExampleNew() {
	ll := leakless.New()

	// just like using the exec.Command
	_ = ll.Command("sleep", "3").Start()

	// get the pid of the guarded sub-process
	select {
	case <-time.Tick(3 * time.Second):
		panic("timeout")
	case pid := <-ll.Pid():
		fmt.Println(pid)
	}
}
