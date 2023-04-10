package main

import (
	_ "embed"

	"github.com/SECRYPT-2022/SECRYPT/command/root"
	"github.com/SECRYPT-2022/SECRYPT/licenses"
)

var (
	//go:embed LICENSE
	license string
)

func main() {
	licenses.SetLicense(license)

	root.NewRootCommand().Execute()
}
