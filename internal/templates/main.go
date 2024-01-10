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

	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"

	"github.com/rdeusser/stacktrace"
	"github.com/rdeusser/log"

	"{{ .ModulePath }}/version"
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

	logger, atom := log.New()
	defer logger.Sync()

	cmd := &cobra.Command{
		Use:     "{{ .ProjectName }} [command]",
		Short:   "Project description",
		Version: version.GetHumanVersion(),
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			atom.SetLevel(zapcore.InfoLevel)

			if options.Debug {
				atom.SetLevel(zapcore.DebugLevel)
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

	cmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	cmd.PersistentFlags().BoolVar(&options.Debug, "debug", options.Debug, "Run in debug mode.")

	if err := cmd.Execute(); err != nil {
		if options.Debug {
			stacktrace.Throw(err)
		}

		logger.Fatal(stacktrace.Unwrap(err).Error())
	}
}`
