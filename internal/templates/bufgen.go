package templates

import (
	"github.com/rdeusser/vulcan/internal/scaffold"
)

var _ scaffold.Template = &BufGen{}

type BufGen struct {
	scaffold.TemplateMixin
	scaffold.ProtobufMixin
}

func (t *BufGen) GetIfExistsAction() scaffold.IfExistsAction {
	return t.IfExistsAction
}

func (t *BufGen) SetTemplateDefaults() error {
	if t.Path == "" {
		t.Path = "buf.gen.yaml"
	}

	t.TemplateBody = bufTemplate
	t.IfExistsAction = scaffold.Overwrite

	return nil
}

const bufGenTemplate = `version: v1
plugins:
  - name: go
    out: .
    opt:
      - paths=source_relative
  - name: go-grpc
    out: .
    opt:
      - paths=source_relative`
