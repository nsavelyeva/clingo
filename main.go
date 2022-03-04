package main

import (
	"clingo/cmd"
	_ "clingo/weather"

	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.NewRootCommand().Execute())
}
