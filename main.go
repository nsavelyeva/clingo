package main

import (
	"clingo/cmd"
	_ "clingo/slack"

	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.NewRootCommand().Execute())
}
