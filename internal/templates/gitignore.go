package templates

import "github.com/rdeusser/vulcan/internal/scaffold"

var _ scaffold.Template = &GitIgnore{}

type GitIgnore struct {
	scaffold.TemplateMixin
	scaffold.ProjectNameMixin
}

func (t *GitIgnore) GetIfExistsAction() scaffold.IfExistsAction {
	return t.IfExistsAction
}

func (t *GitIgnore) SetTemplateDefaults() error {
	if t.Path == "" {
		t.Path = ".gitignore"
	}

	t.TemplateBody = gitIgnoreTemplate
	t.IfExistsAction = scaffold.Overwrite

	return nil
}

const gitIgnoreTemplate = `### {{ .ProjectName }} ###
bin/

### Emacs ###
# -*- mode: gitignore; -*-
*~
\#*\#
/.emacs.desktop
/.emacs.desktop.lock
*.elc
auto-save-list
tramp
.\#*

# Org-mode
.org-id-locations
*_archive

# flymake-mode
*_flymake.*

# eshell files
/eshell/history
/eshell/lastdir

# elpa packages
/elpa/

# reftex files
*.rel

# AUCTeX auto folder
/auto/

# cask packages
.cask/
dist/

# Flycheck
flycheck_*.el

# server auth directory
/server/

# projectiles files
.projectile

# directory configuration
.dir-locals.el

# network security
/network-security.data


### Go ###
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

### Go Patch ###
/vendor/
/Godeps/

### Vim ###
# Swap
[._]*.s[a-v][a-z]
!*.svg  # comment out if you don't need vector files
[._]*.sw[a-p]
[._]s[a-rt-v][a-z]
[._]ss[a-gi-z]
[._]sw[a-p]

# Session
Session.vim
Sessionx.vim

# Temporary
.netrwhist
# Auto-generated tag files
tags
# Persistent undo
[._]*.un~

### VisualStudioCode ###
.vscode/*
!.vscode/settings.json
!.vscode/tasks.json
!.vscode/launch.json
!.vscode/extensions.json
*.code-workspace

### VisualStudioCode Patch ###
# Ignore all local history of files
.history`
