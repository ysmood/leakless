//go:generate go run ./cmd/pack

package leakless

import (
	"bytes"
	"encoding/base64"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/mholt/archiver"
	"github.com/ysmood/byframe"
	"github.com/ysmood/kit"
	"github.com/ysmood/leakless/lib"
)

// Launcher struct
type Launcher struct {
	pid chan int
}

// New leakless instance
func New() *Launcher {
	return &Launcher{
		pid: make(chan int),
	}
}

// Command will try to download the leakless bin and prefix the exec.Cmd with the leakless options.
func (l *Launcher) Command(name string, arg ...string) *exec.Cmd {
	bin := l.getLeaklessBin()

	uid := kit.RandString(16)
	addr := l.serve(uid)

	arg = append([]string{uid, addr, name}, arg...)
	return exec.Command(bin, arg...)
}

// Pid signals the pid of the guarded sub-process. The channel may never receive the pid.
func (l *Launcher) Pid() chan int {
	return l.pid
}

func (l *Launcher) serve(uid string) string {
	srv, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic("[leakless] serve error: " + err.Error())
	}

	go func() {
		conn, err := srv.Accept()
		if err != nil {
			return
		}

		_, err = conn.Write(byframe.Encode([]byte(uid)))
		kit.E(err)

		s := byframe.NewScanner(conn).Limit(100)
		s.Scan()
		pid, err := strconv.ParseInt(string(s.Frame()), 10, 64)
		kit.E(err)

		l.pid <- int(pid)

		s.Scan()
	}()

	return srv.Addr().String()
}

func (l *Launcher) getLeaklessBin() string {
	dir := filepath.Join(os.TempDir(), "leakless-"+lib.Version)
	bin := filepath.Join(dir, "leakless")

	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	if !kit.FileExists(bin) {
		gz := &archiver.Gz{CompressionLevel: 9}
		decompressed := bytes.NewBuffer(nil)
		data, err := base64.StdEncoding.DecodeString(leaklessBin)
		kit.E(err)
		kit.E(gz.Decompress(bytes.NewReader(data), decompressed))

		err = kit.OutputFile(bin, decompressed.Bytes(), nil)
		kit.E(err)
		kit.E(kit.Chmod(bin, 0755))
	}

	return bin
}
