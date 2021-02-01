package main

import (
	"github.com/tzvetkoff-go/pasteur/commands"
)

func main() {
	_ = commands.NewRootCommand().Execute()
}
