package app

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/rdeusser/cli"
	"github.com/rdeusser/stacktrace"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/rdeusser/vulcan/internal/scaffold"
	"github.com/rdeusser/vulcan/internal/templates"
	"github.com/rdeusser/vulcan/internal/templates/scripts"
	"github.com/rdeusser/vulcan/pkg/git"
	"github.com/rdeusser/vulcan/pkg/osutil"
)

type NewOptions struct {
	Branch          string
	ModulePath      string
	ProtobufSupport bool
}

func (o *NewOptions) InitDefaults() {
	o.Branch = "main"
	o.ModulePath = ""
	o.ProtobufSupport = false
}

func NewCmdNew() *cobra.Command {
	options := &NewOptions{}
	options.InitDefaults()

	cmd := &cobra.Command{
		Use:   "new [module]",
		Short: "Create a new Go project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				log.Error().Msg("the Go module name should be the only argument to 'new'")

				os.Exit(1)
			}

			options.ModulePath = args[0]

			if err := RunNew(options); err != nil {
				return stacktrace.Propagate(err, "running new command")
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&options.ProtobufSupport, "protobuf", options.ProtobufSupport, "Generates a project with protobuf support")
	cmd.Flags().StringVar(&options.Branch, "branch", options.Branch, "Default branch to use for the repository")

	return cmd
}

func RunNew(options *NewOptions) error {
	config := &scaffold.Config{
		ModulePath:      options.ModulePath,
		ProjectName:     filepath.Base(options.ModulePath),
		ProtobufSupport: options.ProtobufSupport,
	}

	if err := os.Mkdir(config.ProjectName, 0o755); err != nil {
		return err
	}

	if err := os.Chdir(config.ProjectName); err != nil {
		return err
	}

	currentDir, err := filepath.Abs(".")
	if err != nil {
		return stacktrace.Propagate(err, "getting absolute path of current directory")
	}

	if err := scaffold.NewScaffolder().Execute(
		config,
		&scripts.BumpVersion{},
		&scripts.CleanupNoise{},
		&scripts.Tag{},
		&templates.Buf{},
		&templates.BufGen{},
		&templates.Dockerfile{},
		&templates.GitIgnore{},
		&templates.GoMod{},
		&templates.Main{},
		&templates.Makefile{},
		&templates.Readme{},
		&templates.VersionTest{},
		&templates.Version{},
	); err != nil {
		return stacktrace.Propagate(err, "scaffolding project")
	}

	scriptsDir := filepath.Join(currentDir, "scripts")

	err = fs.WalkDir(os.DirFS(currentDir), filepath.Base(scriptsDir), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return stacktrace.Propagate(err, "walking '%s' directory", path)
		}

		if d.IsDir() {
			cli.Debug("Walking '%s'", scriptsDir)
		} else {
			cli.Debug("Looking at '%s'", path)
		}

		if !d.IsDir() && filepath.Ext(path) == ".sh" {
			if err := os.Chmod(path, 0o755); err != nil {
				return stacktrace.Propagate(err, "marking script as executable")
			}
		}

		return nil
	})
	if err != nil {
		return stacktrace.Propagate(err, "walking '%s' directory", scriptsDir)
	}

	if err := osutil.RunCommand("go mod tidy"); err != nil {
		return stacktrace.Propagate(err, "running 'go mod tidy'")
	}

	if err := git.Init(currentDir, options.Branch, options.ModulePath); err != nil {
		return stacktrace.Propagate(err, "initializing git repo")
	}

	log.Info().Msg("Initialized git repo")
	log.Info().Msgf("Successfully generated project at %s", currentDir)

	return nil
}
