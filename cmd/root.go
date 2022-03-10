package cmd

import (
	"clingo/constants"
	"clingo/helpers"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// The name of our config file, without the file extension because viper supports many config file languages.
	defaultConfigFilename = "clingo-conf" // it is clingo-conf.toml in the repository root folder

	// The environment variable prefix of all environment variables bound to our command line flags.
	// For example, --number is bound to CLINGO_NUMBER.
	envPrefix = "CLINGO"
)

// EventMetadata is a struct to store metadata about a personal or public event
type EventMetadata struct {
	Year   int    `json:"year"`
	Remind int    `json:"remind"`
	Type   string `json:"type"`
	Event  string `json:"event"`
}

// NewRootCommand builds the cobra command that handles our command line tool.
func NewRootCommand() *cobra.Command {
	// Store the result of binding cobra flags and viper config. In a
	// real application these would be data structures, most likely
	// custom structs per command. This is simplified for the demo app and is
	// not recommended that you use one-off variables. The point is that we
	// aren't retrieving the values directly from viper or flags, we read the values
	// from standard Go data structures.
	events := ""
	filter := ""
	output := ""

	// Define our command
	rootCmd := &cobra.Command{
		Use:   "clingo",
		Short: "Check if today is a special day",
		Long:  `Based on a JSON file, check if today is a special day`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			content := helpers.ReadJSON(events)
			var details map[string]EventMetadata
			err := json.Unmarshal(content, &details)
			if err != nil {
				fmt.Printf(`Error loading JSON from "%s": %s`, constants.EventsDefaultJSONFilePath, err)
			}
			today := time.Now()
			// today = time.Date(2022, time.March, 26, 23, 12, 5, 3, time.UTC)

			dt := helpers.GetMonthDay(today, 0)
			if _, ok := details[dt]; ok && (filter == "" || details[dt].Type == filter) {
				output += fmt.Sprintf("Today is %d %s %d: %s [%d year(s)]\n",
					today.Day(), today.Month(), today.Year(),
					details[dt].Event, today.Year()-details[dt].Year)
			}

			// Now scan for the upcoming events with reminders
			for i := 1; i < 10; i++ {
				dt = helpers.GetMonthDay(today, i)
				if _, ok := details[dt]; ok {
					if i <= details[dt].Remind && (filter == "" || details[dt].Type == filter) {
						output += fmt.Sprintf("In %d day(s) will be %d-%s: %s [%d year(s)]\n",
							i, today.Year(), dt, details[dt].Event, today.Year()-details[dt].Year)
					}
				}
			}
			if output == "" {
				output = "No events today.\nNo reminders today.\n"
			}
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			out := cmd.OutOrStdout()

			// Print the final resolved value from binding cobra flags and viper config
			_, _ = fmt.Fprint(out, "", output)
		},
	}

	// Define cobra flags, the precedence is as follows:
	// The path to JSON file with events should come from flag first,
	// then env var CLINGO_EVENTS,
	// then the config file,
	// then the default last.
	rootCmd.Flags().StringVarP(&events, "events", "e", constants.EventsDefaultJSONFilePath, "Is today a special day?")
	rootCmd.Flags().StringVarP(&filter, "filter", "f", "", "Filter events by type")

	rootCmd.AddCommand(
		newWeather(),
		newCurrency(),
		newJokes(),
	)

	return rootCmd
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	// Set the base name of the config file, without the file extension.
	v.SetConfigName(defaultConfigFilename)

	// Set as many paths as you like where viper should look for the
	// config file. We are only looking in the current working directory.
	v.AddConfigPath(".")

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	v.SetEnvPrefix(envPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, v)

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to CLINGO_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			err := v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
			if err != nil {
				return
			}
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				return
			}
		}
	})
}
