package templates

import (
	"fmt"

	"github.com/rdeusser/vulcan/internal/scaffold"
)

var _ scaffold.Template = &Main{}

type Main struct {
	scaffold.TemplateMixin
	scaffold.ModulePathMixin
	scaffold.ProjectNameMixin
}

func (t *Main) GetIfExistsAction() scaffold.IfExistsAction {
	return t.IfExistsAction
}

func (t *Main) SetTemplateDefaults() error {
	if t.Path == "" {
		t.Path = fmt.Sprintf("cmd/%s/main.go", t.ProjectName)
	}

	t.TemplateBody = mainTemplate
	t.IfExistsAction = scaffold.Skip

	return nil
}

const mainTemplate = `package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rdeusser/stacktrace"

	"{{ .ModulePath }}/version"
)

type Options struct {
	Debug bool
}

func (o *Options) InitDefaults() {
	viper.SetEnvPrefix("{{ toEnv .ProjectName }}")
	viper.AutomaticEnv()

	o.Debug = viper.GetBool("debug")
}

func main() {
	options := &Options{}
	options.InitDefaults()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	cmd := &cobra.Command{
		Use:     "{{ .ProjectName }} [command]",
		Short:   "Project description",
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

	if err := cmd.Execute(); err != nil {
		if options.Debug {
			stacktrace.Throw(err)
		}

		log.Fatal().Msg(stacktrace.Unwrap(err).Error())
	}
}`
