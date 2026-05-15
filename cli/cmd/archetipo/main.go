package main

import (
	"os"

	appcli "github.com/techreloaded-ar/ARchetipo/cli/internal/cli"
)

func main() {
	os.Exit(appcli.Execute(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
