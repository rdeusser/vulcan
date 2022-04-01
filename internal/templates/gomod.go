package templates

import "github.com/rdeusser/vulcan/internal/scaffold"

var _ scaffold.Template = &GoMod{}

type GoMod struct {
	scaffold.TemplateMixin
	scaffold.ModulePathMixin
	scaffold.ProjectNameMixin
}

func (t *GoMod) SetTemplateDefaults() error {
	if t.Path == "" {
		t.Path = "go.mod"
	}

	t.TemplateBody = goModTemplate
	t.IfExistsAction = scaffold.Overwrite

	return nil
}

const goModTemplate = `module {{ .ModulePath }}

go 1.18`
