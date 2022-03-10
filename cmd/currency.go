package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"clingo/currency"
)

func newCurrency() *cobra.Command {
	var conf currency.ConfigCurrency

	cmd := &cobra.Command{
		Use:   "currency",
		Short: "Currency rate",
		Long:  "Request currency rate information for the given currency using specified base currency",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sc := *currency.NewServiceCurrency(conf.From, conf.To, conf.Token)
			return currency.Run(cmd.OutOrStdout(), sc, &conf)
		},
	}

	bindCurrencyFlags(cmd.Flags(), &conf)

	return cmd
}

func bindCurrencyFlags(flags *pflag.FlagSet, config *currency.ConfigCurrency) {
	flags.StringVar(&config.From, "from", "EUR", "currency from")
	flags.StringVar(&config.To, "to", "USD", "currency to")
	flags.StringVar(&config.Token, "token", "", "currency token")
}
