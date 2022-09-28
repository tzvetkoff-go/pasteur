package main

import (
	"github.com/tzvetkoff-go/pasteur/pkg/cli"
)

func main() {
	_ = cli.NewRootCommand().Execute()
}
