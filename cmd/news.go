package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"clingo/news"
)

func newNews() *cobra.Command {
	var conf news.ConfigNews

	cmd := &cobra.Command{
		Use:   "news",
		Short: "Top News",
		Long:  "Request top news for the given date, language, and source",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sc := *news.NewServiceNews(conf.Date, conf.Language, conf.Source, conf.Limit, conf.Markup, conf.Token)
			return news.Run(cmd.OutOrStdout(), sc, &conf)
		},
	}

	bindNewsFlags(cmd.Flags(), &conf)

	return cmd
}

func bindNewsFlags(flags *pflag.FlagSet, config *news.ConfigNews) {
	flags.StringVar(&config.Date, "date", "", "news date")
	flags.StringVar(&config.Language, "language", "en", "news language")
	flags.StringVar(&config.Source, "source", "google-news-en", "news source")
	flags.IntVar(&config.Limit, "limit", 100, "news limit")
	flags.BoolVar(&config.Markup, "markup", false, "news markup")
	flags.StringVar(&config.Token, "token", "", "news token")
}
