//go:generate go run ./cmd/pack

package leakless

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/ysmood/byframe/v2"
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

	uid := fmt.Sprintf("%x", lib.RandBytes(16))
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
		lib.E(err)
		_, err = conn.Write(data)
		lib.E(err)

		s := byframe.NewScanner(conn).Limit(1000)
		for s.Scan() {
			var msg lib.Message
			err = s.Decode(&msg)
			lib.E(err)

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

	if !lib.FileExists(bin) {
		raw, err := base64.StdEncoding.DecodeString(leaklessBin)
		lib.E(err)
		gr, err := gzip.NewReader(bytes.NewBuffer(raw))
		lib.E(err)
		data, err := ioutil.ReadAll(gr)
		lib.E(err)
		lib.E(gr.Close())

		err = lib.OutputFile(bin, data, nil)
		lib.E(err)
		lib.E(os.Chmod(bin, 0755))
	}

	return bin
}
