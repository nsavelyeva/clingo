package cmd

import (
	"clingo/jokes"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newJokes() *cobra.Command {
	var conf jokes.ConfigJokes

	cmd := &cobra.Command{
		Use:   "jokes",
		Short: "Random joke",
		Long:  "Request a random short joke",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return jokes.Run(cmd.OutOrStdout(), &conf)
		},
	}

	bindJokesFlags(cmd.Flags(), &conf)

	return cmd
}

func bindJokesFlags(flags *pflag.FlagSet, config *jokes.ConfigJokes) {
	flags.BoolVar(&config.Emoji, "emoji", false, "jokes emoji")
	flags.StringVar(&config.Token, "token", "", "jokes token")
}
