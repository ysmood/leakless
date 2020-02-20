package leakless

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/mholt/archiver"
	"github.com/ysmood/byframe"
	"github.com/ysmood/kit"
	"github.com/ysmood/leakless/lib"
)

// Leakless ...
type Leakless struct {
	bin  string
	host string
}

// New leakless instance
func New() *Leakless {
	return &Leakless{}
}

// Bin sets the leakless bin location
func (l *Leakless) Bin(path string) *Leakless {
	l.bin = path
	return l
}

// Host sets the host to download leakless bin
func (l *Leakless) Host(host string) *Leakless {
	l.host = host
	return l
}

// Command will try to download the leakless bin and prefix the exec.Cmd with
// the leakless options. If it fails to download the bin a normal exec.Cmd will be returned.
func (l *Leakless) Command(name string, arg ...string) *exec.Cmd {
	bin := l.bin
	if bin == "" {
		var err error
		bin, err = l.getLeaklessBin()
		if err != nil {
			return exec.Command(name, arg...)
		}
	}

	uid := kit.RandString(16)
	addr := l.serve(uid)

	arg = append([]string{uid, addr, name}, arg...)
	return exec.Command(bin, arg...)
}

func (l *Leakless) serve(uid string) string {
	srv, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic("[leakless] serve error: " + err.Error())
	}

	go func() {
		for {
			conn, err := srv.Accept()
			if err != nil {
				return
			}

			_, _ = conn.Write(byframe.Encode([]byte(uid)))
			buf := make([]byte, 1)
			_, _ = conn.Read(buf)
		}
	}()

	return srv.Addr().String()
}

func (l *Leakless) getLeaklessBin() (string, error) {
	dir := filepath.Join(os.TempDir(), "leakless-"+lib.Version)
	bin := filepath.Join(dir, "leakless")

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

	err = archiver.Unarchive(zip, dir)
	if err != nil {
		return "", err
	}

	return bin, nil
}
