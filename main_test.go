package leakless_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ysmood/leakless"
	"github.com/ysmood/leakless/lib"
)

var p = filepath.FromSlash

type stamp struct {
	Pid  int
	Time time.Time
}

func TestMain(m *testing.M) {
	binDir := filepath.Join(os.TempDir(), "leakless-"+lib.Version)
	lib.E(os.RemoveAll(binDir))
	lib.E(lib.Mkdir("dist", nil))
	lib.Exec("go", "dist", "build", "../cmd/test")
	lib.Exec("go", "dist", "build", "../cmd/zombie")

	os.Exit(m.Run())
}

func TestBasic(t *testing.T) {
	lib.Exec(p("dist/test"), "", "on")

	lib.Sleep(2)
	var s stamp
	var pid int
	_ = lib.ReadJSON(p("tmp/pid"), &s)
	_ = lib.ReadJSON(p("tmp/sub-pid"), &pid)

	assert.NotEmpty(t, s.Pid)
	assert.Equal(t, s.Pid, pid)
	assert.True(t, time.Since(s.Time) > time.Second)
}

func TestErr(t *testing.T) {
	l := leakless.New()
	lib.E(l.Command("not-exists").Start())

	pid := <-l.Pid()
	assert.Zero(t, pid)
	assert.Regexp(t, `executable file not found`, l.Err())
}

func TestZombie(t *testing.T) {
	cmd := exec.Command(p("dist/test"), "off")

	lib.E(cmd.Run())

	lib.Sleep(2)
	var s stamp
	_ = lib.ReadJSON(p("tmp/pid"), &s)

	assert.NotEmpty(t, s.Pid)
	assert.True(t, time.Since(s.Time) < time.Second)
}
