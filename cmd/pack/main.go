package main

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/ysmood/leakless/lib"
)

func main() {
	setVersion()
	lib.Exec("godev", "", "lint")
	lib.Exec("godev", "", "build", "-n", "./cmd/leakless")
	pack("linux")
	pack("darwin")
	pack("windows")
}

func pack(osName string) {
	var bin []byte
	var err error

	switch osName {
	case "linux":
		bin, err = lib.ReadFile(filepath.FromSlash("dist/leakless-linux"))
	case "darwin":
		bin, err = lib.ReadFile(filepath.FromSlash("dist/leakless-mac"))
	case "windows":
		bin, err = lib.ReadFile(filepath.FromSlash("dist/leakless-windows"))
	default:
		panic("unsupported os")
	}

	lib.E(err)

	buf := bytes.Buffer{}
	gw, err := gzip.NewWriterLevel(&buf, 9)
	lib.E(err)
	lib.E(gw.Write(bin))
	lib.E(gw.Close())

	tpl := `package leakless

var leaklessBin = "%s"
`
	code := fmt.Sprintf(tpl, base64.StdEncoding.EncodeToString(buf.Bytes()))

	lib.E(lib.OutputFile(fmt.Sprintf("bin_%s.go", osName), code, nil))
}

func setVersion() {
	files, err := filepath.Glob("cmd/leakless/*.go")
	lib.E(err)

	args := append([]string{"hash-object"}, files...)

	raw, err := exec.Command("git", args...).CombinedOutput()
	lib.E(err)

	hash := md5.Sum(raw)

	lib.E(lib.OutputFile("lib/version.go", fmt.Sprintf(`package lib

// Version ...
const Version = "%x"
`, hash), nil))
}
