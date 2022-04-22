package scripts

import "github.com/rdeusser/vulcan/internal/scaffold"

var _ scaffold.Template = &Tag{}

type Tag struct {
	scaffold.TemplateMixin
	scaffold.ProjectNameMixin
}

func (t *Tag) GetIfExistsAction() scaffold.IfExistsAction {
	return t.IfExistsAction
}

func (t *Tag) SetTemplateDefaults() error {
	if t.Path == "" {
		t.Path = "scripts/tag.sh"
	}

	t.TemplateBody = tagTemplate
	t.IfExistsAction = scaffold.Overwrite

	return nil
}

const tagTemplate = `#!/usr/bin/env bash

repo_root=$(git rev-parse --show-toplevel)
version=$(grep -oE "[0-9]+[.][0-9]+[.][0-9]+" "${repo_root}/version/version.go")
remote=$(git remote -v | awk '{print $1}' | head -n 1)

git tag -a "v${version}" -m "v${version}"
git push "$remote" "v${version}"`
