package templates

import (
	"github.com/rdeusser/vulcan/internal/scaffold"
)

var _ scaffold.Template = &Buf{}

type Buf struct {
	scaffold.TemplateMixin
	scaffold.ProtobufMixin
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
build:
  roots:
    - .
breaking:
  use:
    - FILE`
