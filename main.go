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

	"github.com/mholt/archiver/v3"
	"github.com/ysmood/byframe/v2"
	"github.com/ysmood/kit"
	"github.com/ysmood/leakless/lib"
)

// Launcher struct
type Launcher struct {
	pid chan int
	err string
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

// Err message from the guard process
func (l *Launcher) Err() string {
	return l.err
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

		data, err := byframe.Encode(lib.Message{UID: uid})
		kit.E(err)
		_, err = conn.Write(data)
		kit.E(err)

		s := byframe.NewScanner(conn).Limit(1000)
		for s.Scan() {
			var msg lib.Message
			err = s.Decode(&msg)
			kit.E(err)

			l.err = msg.Error
			l.pid <- msg.PID
		}
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
