package leakless

import (
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
	bin  string
	host string
	pid  chan int
}

// New leakless instance
func New() *Launcher {
	return &Launcher{
		pid: make(chan int),
	}
}

// Bin sets the leakless bin location
func (l *Launcher) Bin(path string) *Launcher {
	l.bin = path
	return l
}

// Host sets the host to download leakless bin
func (l *Launcher) Host(host string) *Launcher {
	l.host = host
	return l
}

// Command will try to download the leakless bin and prefix the exec.Cmd with the leakless options.
func (l *Launcher) Command(name string, arg ...string) *exec.Cmd {
	bin := l.bin
	if bin == "" {
		var err error
		bin, err = l.getLeaklessBin()
		if err != nil {
			panic(err)
		}
	}

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

func (l *Launcher) getLeaklessBin() (string, error) {
	dir := filepath.Join(os.TempDir(), "leakless-"+lib.Version)
	bin := filepath.Join(dir, "leakless")

	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	if kit.FileExists(bin) {
		return bin, nil
	}

	host := l.host
	if host == "" {
		host = "https://github.com/ysmood/leakless/releases/download/"
	}

	urlPrefix := host + lib.Version + "/"
	zipName := ""

	switch runtime.GOOS {
	case "linux":
		zipName += "leakless-linux.tar.gz"
	case "darwin":
		zipName += "leakless-mac.zip"
	case "windows":
		zipName += "leakless-windows.zip"
	}

	data, err := kit.Req(urlPrefix + zipName).Bytes()
	if err != nil {
		return "", err
	}

	zip := filepath.Join(dir, zipName)

	err = kit.OutputFile(zip, data, nil)
	if err != nil {
		return "", err
	}

	_ = os.Remove(bin)

	err = archiver.Unarchive(zip, dir)
	if err != nil {
		return "", err
	}

	return bin, nil
}
