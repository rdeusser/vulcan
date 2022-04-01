package scaffold

// IfExistsAction determines what to do if the scaffold file already exists.
type IfExistsAction int

const (
	// Skip skips the file and moves to the next one.
	Skip IfExistsAction = iota

	// Error returns an error and stops processing.
	Error

	// Overwrite truncates and overwrites the existing file.
	Overwrite
)

type TemplateMixin struct {
	// Path is the path of the file.
	Path string

	// TemplateBody is the template body to execute.
	TemplateBody string

	// IFExistsAction determines what to do if the file exists.
	IfExistsAction IfExistsAction
}

func (m *TemplateMixin) GetPath() string {
	return m.Path
}

func (m *TemplateMixin) GetBody() string {
	return m.TemplateBody
}

type ModulePathMixin struct {
	// ModulePath is the name of the Go module.
	ModulePath string
}

// InjectModulePath implements HasModulePath.
func (m *ModulePathMixin) InjectModulePath(modulePath string) {
	if m.ModulePath == "" {
		m.ModulePath = modulePath
	}
}

type ProjectNameMixin struct {
	// ProjectName is the name of the project.
	ProjectName string
}

// InjectProjectName implements HasProjectName.
func (m *ProjectNameMixin) InjectProjectName(projectName string) {
	if m.ProjectName == "" {
		m.ProjectName = projectName
	}
}

type ProtobufMixin struct {
	// ProtobufSupport signals that the project should be generated with
	// protobuf support.
	ProtobufSupport bool
}

// InjectProtobufSupport implements ProtobufSupport.
func (m *ProtobufMixin) InjectProtobufSupport(protobufSupport bool) {
	if !m.ProtobufSupport {
		m.ProtobufSupport = protobufSupport
	}
}
