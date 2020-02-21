package leakless_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/kit"
	"github.com/ysmood/leakless/lib"
)

var p = filepath.FromSlash

type stamp struct {
	PID  int
	Time string
}

func TestMain(m *testing.M) {
	kit.E(kit.Mkdir("dist", nil))
	kit.Exec("go", "build", "../cmd/leakless").Dir("dist").MustDo()
	kit.Exec("go", "build", "../cmd/test").Dir("dist").MustDo()
	kit.Exec("go", "build", "../cmd/zombie").Dir("dist").MustDo()

	os.Exit(m.Run())
}

func TestBasic(t *testing.T) {
	cmd := exec.Command(p("dist/test"), "on", "exit")

	kit.E(cmd.Run())

	done := false
	prev := ""
	for range make([]int, 5) {
		kit.Sleep(1)
		var s stamp
		_ = kit.ReadJSON(p("tmp/pid"), &s)
		assert.NotEmpty(t, s.Time)

		done = prev == s.Time
		prev = s.Time
	}
	assert.True(t, done)
}

func TestSignal(t *testing.T) {
	cmd := exec.Command(p("dist/test"), "on", "off")

	kit.E(cmd.Start())

	kit.Sleep(1)

	kit.E(lib.Signal(cmd, os.Interrupt))

	done := false
	prev := ""
	for range make([]int, 5) {
		kit.Sleep(1)
		var s stamp
		_ = kit.ReadJSON(p("tmp/pid"), &s)
		assert.NotEmpty(t, s.Time)

		done = prev == s.Time
		prev = s.Time
	}
	assert.True(t, done)
}

func TestZombie(t *testing.T) {
	cmd := exec.Command(p("dist/test"), "off", "exit")

	kit.E(cmd.Run())

	done := false
	prev := ""
	for range make([]int, 5) {
		kit.Sleep(1)
		var s stamp
		_ = kit.ReadJSON(p("tmp/pid"), &s)
		assert.NotEmpty(t, s.Time)

		done = prev == s.Time
		prev = s.Time
	}
	assert.False(t, done)
}
