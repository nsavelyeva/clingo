package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"clingo/weather"
)

func newWeather() *cobra.Command {
	var conf weather.ConfigWeather

	cmd := &cobra.Command{
		Use:   "weather",
		Short: "Current weather information in the given city",
		Long:  "Request current weather information for the given city",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sw := *weather.NewServiceWeather(conf.City, conf.Token)
			return weather.Run(cmd.OutOrStdout(), sw)
		},
	}

	bindWeatherFlags(cmd.Flags(), &conf)

	return cmd
}

func bindWeatherFlags(flags *pflag.FlagSet, config *weather.ConfigWeather) {
	flags.StringVar(&config.City, "city", "Amsterdam", "weather city")
	flags.StringVar(&config.Token, "token", "", "weather token")
}
