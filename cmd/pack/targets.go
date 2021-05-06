package main

import "github.com/ysmood/leakless/pkg/utils"

// The targets to pack into the main package
var targets = []utils.Target{
	"linux/amd64",
	"linux/arm64",
	"darwin/amd64",
	"darwin/arm64",
	"windows/amd64",
}
