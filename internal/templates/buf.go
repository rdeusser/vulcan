package templates

import (
	"github.com/rdeusser/vulcan/internal/scaffold"
)

var _ scaffold.Template = &Buf{}

type Buf struct {
	scaffold.TemplateMixin
	scaffold.ProtobufMixin
}

func (t *Buf) GetIfExistsAction() scaffold.IfExistsAction {
	return t.IfExistsAction
}

func (t *Buf) SetTemplateDefaults() error {
	if t.Path == "" {
		t.Path = "buf.yaml"
	}

	t.TemplateBody = bufTemplate
	t.IfExistsAction = scaffold.Overwrite

	return nil
}

const bufTemplate = `version: v1
breaking:
  use:
    - FILE
lint:
  use:
    - DEFAULT`
