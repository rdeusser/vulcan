package scaffold

type Template interface {
	// GetPath returns the path to the file's location.
	GetPath() string

	// GetBody returns the template body.
	GetBody() string

	// GetIfExistsAction determines what to do if the file already exists.
	GetIfExistsAction() IfExistsAction

	// SetTemplateDefaults sets the default values for templates.
	SetTemplateDefaults() error
}

// HasModulePath allows a module path to be used in a template.
type HasModulePath interface {
	// InjectModulePath sets the template module path.
	InjectModulePath(string)
}

// HasProjectName allows a project name to be used in a template.
type HasProjectName interface {
	// InjectProjectName sets the template project name.
	InjectProjectName(string)
}

// HasProtobufSupport signals that the project should be generated with protobuf
// support.
type HasProtobufSupport interface {
	// InjectProtobufSupport injects the template with protobuf support.
	InjectProtobufSupport(bool)
}
