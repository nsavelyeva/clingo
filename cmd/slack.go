package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"clingo/slack"
)

func newSlack() *cobra.Command {
	var conf slack.Config

	cmd := &cobra.Command{
		Use:   "slack",
		Short: "Prints the Slack config",
		Long:  "Requires and prints the Slack config",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return slack.Run(cmd.OutOrStdout(), conf)
		},
	}

	bindSlackFlags(cmd.Flags(), &conf)

	return cmd
}

func bindSlackFlags(flags *pflag.FlagSet, config *slack.Config) {
	flags.StringVar(&config.String, "slack-string", "slack default", "slack string field")
	flags.Float64Var(&config.Float, "slack-float", 0.1, "slack float field")
	flags.BoolVar(&config.Bool, "slack-bool", true, "slack bool field")
}
