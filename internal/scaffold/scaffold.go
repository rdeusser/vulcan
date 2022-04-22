package scaffold

import (
	"bytes"
	"fmt"
	"go/format"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

var templateFuncs = template.FuncMap{
	"backtick": backtick,
	"toEnv":    toEnv,
}

type Config struct {
	// ModulePath is the name of the Go module.
	ModulePath string

	// ProjectName is the name of this project.
	ProjectName string

	// ProtobufSupport signals whether or not to generate the project with
	// protobuf support.
	ProtobufSupport bool
}

type Scaffolder struct {
	fs afero.Fs
}

func NewScaffolder() *Scaffolder {
	return &Scaffolder{
		fs: afero.NewOsFs(),
	}
}

func (s *Scaffolder) Execute(config *Config, templates ...Template) error {
	log.Info().Msg("Writing scaffold")

	for _, t := range templates {
		if builder, ok := t.(HasProjectName); ok {
			builder.InjectProjectName(config.ProjectName)
		}

		if builder, ok := t.(HasModulePath); ok {
			builder.InjectModulePath(config.ModulePath)
		}

		if builder, ok := t.(HasProtobufSupport); ok {
			builder.InjectProtobufSupport(config.ProtobufSupport)
		}

		if err := t.SetTemplateDefaults(); err != nil {
			return fmt.Errorf("failed setting template defaults: %w", err)
		}

		if err := s.writeFile(t); err != nil {
			return err
		}
	}

	return nil
}

func (s *Scaffolder) writeFile(t Template) error {
	switch t.GetIfExistsAction() {
	case Skip:
		return nil
	case Error:
		return ErrFileAlreadyExists
	case Overwrite:
	default:
		return ErrUnknownAction
	}

	path := t.GetPath()

	if err := s.fs.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create parent directories: %w", err)
	}

	writer, err := s.fs.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}

	contents, err := doTemplate(t)
	if err != nil {
		return fmt.Errorf("failed executing template: %w", err)
	}

	_, err = writer.Write(contents)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	log.Info().Msgf("Successfully wrote %s", path)

	return nil
}

func doTemplate(t Template) ([]byte, error) {
	tmpl, err := template.New(fmt.Sprintf("%T", t)).Funcs(templateFuncs).Parse(t.GetBody())
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, t); err != nil {
		return nil, err
	}

	b := buf.Bytes()

	if filepath.Ext(t.GetPath()) == ".go" {
		b, err = format.Source(b)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

func backtick() string {
	return "`"
}

func toEnv(s string) string {
	s = strings.ToUpper(s)
	s = strings.ReplaceAll(s, "-", "_")

	return s
}
