package templates

import "github.com/rdeusser/vulcan/internal/scaffold"

var _ scaffold.Template = &Dockerfile{}

type Dockerfile struct {
	scaffold.TemplateMixin
	scaffold.ProjectNameMixin
}

func (t *Dockerfile) SetTemplateDefaults() error {
	if t.Path == "" {
		t.Path = "Dockerfile"
	}

	t.TemplateBody = dockerfileTemplate
	t.IfExistsAction = scaffold.Overwrite

	return nil
}

const dockerfileTemplate = `FROM golang:1.18 as builder

WORKDIR /src/{{ .ProjectName }}

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN make build

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /src/{{ .ProjectName }}/bin/{{ .ProjectName }} /
ENTRYPOINT ["/{{ .ProjectName }}"]`
