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

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/rdeusser/stacktrace"
	"github.com/rdeusser/x/zappretty"

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

	atom := zap.NewAtomicLevel()
	cfg := zap.NewProductionEncoderConfig()

	zappretty.Register(cfg)

	cliEncoder := zappretty.NewCLIEncoder(cfg)
	jsonEncoder := zapcore.NewJSONEncoder(cfg)

	leveler := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= atom.Level()
	})

	var core zapcore.Core

	if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		core = zapcore.NewCore(cliEncoder, os.Stdout, leveler)
	} else {
		core = zapcore.NewCore(jsonEncoder, os.Stdout, leveler)
	}

	logger := zap.New(core, zap.AddStacktrace(atom)).Named("{{ .ProjectName }}")
	defer logger.Sync()

	cmd := &cobra.Command{
		Use:     "{{ .ProjectName }} [command]",
		Short:   "Project description",
		Version: version.GetHumanVersion(),
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

	cmd.PersistentFlags().BoolVar(&options.Debug, "debug", options.Debug, "Run in debug mode.")

	if err := cmd.Execute(); err != nil {
		if options.Debug {
			stacktrace.Throw(err)
		}

		logger.Fatal(stacktrace.Unwrap(err).Error())
	}
}`
