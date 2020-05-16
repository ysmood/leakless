package leakless_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ysmood/kit"
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
	kit.E(os.RemoveAll(binDir))
	kit.E(kit.Mkdir("dist", nil))
	kit.Exec("go", "build", "../cmd/test").Dir("dist").MustDo()
	kit.Exec("go", "build", "../cmd/zombie").Dir("dist").MustDo()

	os.Exit(m.Run())
}

func TestBasic(t *testing.T) {
	kit.Exec(p("dist/test"), "on").MustDo()

	kit.Sleep(2)
	var s stamp
	var pid int
	_ = kit.ReadJSON(p("tmp/pid"), &s)
	_ = kit.ReadJSON(p("tmp/sub-pid"), &pid)

	assert.NotEmpty(t, s.Pid)
	assert.Equal(t, s.Pid, pid)
	assert.True(t, time.Since(s.Time) > time.Second)
}

func TestErr(t *testing.T) {
	l := leakless.New()
	kit.E(l.Command("not-exists").Start())

	pid := <-l.Pid()
	assert.Zero(t, pid)
	assert.Regexp(t, `executable file not found`, l.Err())
}

func TestZombie(t *testing.T) {
	cmd := exec.Command(p("dist/test"), "off")

	kit.E(cmd.Run())

	kit.Sleep(2)
	var s stamp
	_ = kit.ReadJSON(p("tmp/pid"), &s)

	assert.NotEmpty(t, s.Pid)
	assert.True(t, time.Since(s.Time) < time.Second)
}
