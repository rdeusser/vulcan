package templates

import "github.com/rdeusser/vulcan/internal/scaffold"

var _ scaffold.Template = &Version{}

type Version struct {
	scaffold.TemplateMixin
	scaffold.ProjectNameMixin
}

func (t *Version) GetIfExistsAction() scaffold.IfExistsAction {
	return t.IfExistsAction
}

func (t *Version) SetTemplateDefaults() error {
	if t.Path == "" {
		t.Path = "version/version.go"
	}

	t.TemplateBody = versionTemplate
	t.IfExistsAction = scaffold.Skip

	return nil
}

const versionTemplate = `package version

import (
	"fmt"
	"strings"
)

var (
	// GitCommit represents the git commit that {{ .ProjectName }} was compiled
	// with. These will be filled in by the compiler.
	GitCommit   string
	GitDescribe string

	// The main version number that is being run at the moment.
	//
	// Version must conform to the format expected by
	// github.com/hashicorp/go-version for tests to work.
	Version = "0.1.0"

	// VersionPrerelease is a pre-release marker for the version. If this is
	// "" (empty string) then it means that that this is a final
	// release. Otherwise, this is a pre-release such as "dev" (in
	// development), "beta", "rc1", etc.
	VersionPrerelease = "dev"
)

// GetHumanVersion composes the parts of the version in a way that's suitable for displaying to humans.
func GetHumanVersion() string {
	version := Version

	if GitDescribe != "" {
		version = GitDescribe
	}

	release := VersionPrerelease
	if GitDescribe == "" && release == "" {
		release = "dev"
	}

	if release != "" {
		if !strings.HasSuffix(version, "-"+release) {
			// If we tagged a pre-release version then the release is in the version already.
			version += fmt.Sprintf("-%s", release)
		}

		if GitCommit != "" {
			version += fmt.Sprintf(" (%s)", GitCommit)
		}
	}

	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	// Strip off any single quotes added by Git.
	return strings.ReplaceAll(version, "'", "")
}`
