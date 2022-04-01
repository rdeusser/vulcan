package main

import (
	"fmt"
	"os"

	"github.com/rdeusser/stacktrace"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rdeusser/vulcan/cmd/vulcan/app"
	"github.com/rdeusser/vulcan/version"
)

type Options struct {
	Debug bool
}

func (o *Options) InitDefaults() {
	o.Debug = false
}

func main() {
	options := &Options{}
	options.InitDefaults()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	cmd := &cobra.Command{
		Use:     "vulcan [command]",
		Short:   "A project generator",
		Version: version.GetHumanVersion(),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)

			if options.Debug {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println(cmd.UsageString())
			}
		},
		SilenceUsage:  true,
		SilenceErrors: true, // we log and return our own errors and if this is false the errors are printed twice
	}

	cmd.PersistentFlags().BoolVar(&options.Debug, "debug", options.Debug, "Run in debug mode.")

	cmd.AddCommand(app.NewCmdNew())

	viper.AutomaticEnv()

	if err := cmd.Execute(); err != nil {
		if options.Debug {
			stacktrace.Throw(err)
		}

		log.Fatal().Msg(stacktrace.Error(err).Error())
	}
}
