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

	"github.com/ysmood/leakless/pkg/utils"
)

func main() {
	setVersion()

	utils.E(os.RemoveAll("dist"))
	utils.E(os.MkdirAll("dist", 0755))

	for _, target := range targets {
		pack(target)
	}
}

func pack(target utils.Target) {
	var bin []byte
	var err error
	name := target.BinName()

	build(target)

	bin, err = utils.ReadFile(filepath.FromSlash("dist/leakless-" + name))
	utils.E(err)

	buf := bytes.Buffer{}
	gw, err := gzip.NewWriterLevel(&buf, 9)
	utils.E(err)
	utils.E(gw.Write(bin))
	utils.E(gw.Close())

	tpl := `package leakless

func init() {
	leaklessBinaries["%s"] = "%s"
}
`
	code := fmt.Sprintf(tpl, name, base64.StdEncoding.EncodeToString(buf.Bytes()))

	utils.E(utils.OutputFile(fmt.Sprintf("bin_%s.go", name), code, nil))
}

func setVersion() {
	a, err := filepath.Glob("cmd/leakless/*.go")
	utils.E(err)

	b, err := filepath.Glob("cmd/pack/*.go")
	utils.E(err)

	files := append(a, b...)

	args := append([]string{"hash-object"}, files...)

	raw, err := exec.Command("git", args...).CombinedOutput()
	utils.E(err)

	hash := md5.Sum(raw)

	utils.E(utils.OutputFile("pkg/shared/version.go", fmt.Sprintf(`package shared

// Version ...
const Version = "%x"
`, hash), nil))
}

func build(target utils.Target) {
	o, err := exec.Command("zig", "build-exe",
		"-O", "ReleaseSmall", "--strip",
		"cmd/leakless/main.zig",
		"--target", string(target),
	).CombinedOutput()
	if err != nil {
		panic(string(o))
	}

	dest := filepath.FromSlash("dist/leakless-" + target.BinName())

	utils.E(os.Rename("main", dest))
}
