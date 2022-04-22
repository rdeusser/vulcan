package scripts

import "github.com/rdeusser/vulcan/internal/scaffold"

var _ scaffold.Template = &CleanupNoise{}

type CleanupNoise struct {
	scaffold.TemplateMixin
	scaffold.ProjectNameMixin
}

func (t *CleanupNoise) GetIfExistsAction() scaffold.IfExistsAction {
	return t.IfExistsAction
}

func (t *CleanupNoise) SetTemplateDefaults() error {
	if t.Path == "" {
		t.Path = "scripts/cleanup-noise.sh"
	}

	t.TemplateBody = cleanupNoiseTemplate
	t.IfExistsAction = scaffold.Overwrite

	return nil
}

const cleanupNoiseTemplate = `#!/usr/bin/env bash

SED_BIN=${SED_BIN:-sed}

${SED_BIN} -i 's/[ \t]*$//' "$@"`
