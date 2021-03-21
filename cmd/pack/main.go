package main

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ysmood/leakless/lib"
)

func main() {
	setVersion()

	lib.E(os.RemoveAll("dist"))

	pack("linux")
	pack("darwin")
	pack("windows")
}

func pack(osName string) {
	var bin []byte
	var err error

	build(osName)

	bin, err = lib.ReadFile(filepath.FromSlash("dist/leakless-" + osName))
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

func build(osName string) {
	flags := []string{
		"build",
		"-trimpath",
		"-o", filepath.FromSlash("dist/leakless-" + osName),
	}

	ldFlags := "-ldflags=-w -s"
	if osName == "windows" {
		// On Windows, -H windowsgui writes a "GUI binary" instead of a "console binary."
		ldFlags += " -H=windowsgui"
	}
	flags = append(flags, ldFlags)

	flags = append(flags, filepath.FromSlash("./cmd/leakless"))

	cmd := exec.Command("go", flags...)
	cmd.Env = append(os.Environ(), []string{
		"GOOS=" + osName,
		"GOARCH=amd64",
	}...)
	lib.E(cmd.Run())
}
