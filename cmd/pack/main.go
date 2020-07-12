package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"path/filepath"

	"github.com/ysmood/kit"
)

func main() {
	kit.Exec("godev", "build", "-n", "./cmd/leakless").MustDo()
	pack("linux")
	pack("darwin")
	pack("windows")
}

func pack(osName string) {
	var bin []byte
	var err error

	switch osName {
	case "linux":
		bin, err = kit.ReadFile(filepath.FromSlash("dist/leakless-linux"))
	case "darwin":
		bin, err = kit.ReadFile(filepath.FromSlash("dist/leakless-mac"))
	case "windows":
		bin, err = kit.ReadFile(filepath.FromSlash("dist/leakless-windows"))
	default:
		panic("unsupported os")
	}

	kit.E(err)

	buf := bytes.Buffer{}
	gw, err := gzip.NewWriterLevel(&buf, 9)
	kit.E(err)
	kit.E(gw.Write(bin))
	kit.E(gw.Close())

	tpl := `package leakless

var leaklessBin = "%s"
`
	code := fmt.Sprintf(tpl, base64.StdEncoding.EncodeToString(buf.Bytes()))

	kit.E(kit.OutputFile(fmt.Sprintf("bin_%s.go", osName), code, nil))
}
