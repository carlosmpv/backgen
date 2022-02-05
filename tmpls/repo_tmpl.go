package tmpls

import (
	"io"
	"strings"
	"text/template"
)

const repoTmpl = `
package {{ .Package }}

import "github.com/rs/xid"

func New{{ .Name }}Repository(ps pubsub) *{{ .Name }}Repository {
	return &{{ .Name }}Repository{
		registry: {{ .Name }}Registry{},
		ps: ps,
	}
}

type {{ .Name }}Repository struct {
	registry {{ .Name }}Registry
	ps pubsub
}

func (r {{ .Name }}Repository) Create(m *{{ .Name }}) (string, error) {
	m.ID = xid.New().String()

	err := m.Save(r.ps.Set)
	if err != nil {
		return "", err
	}

	return m.ID, r.registry.Register(r.ps.Set, m.ID)
}

func (r {{ .Name }}Repository) Delete(id string) error {
	return r.registry.Delete(r.ps.Set, id)
}

func (r {{ .Name }}Repository) GetByID(id string) (*{{ .Name }}, error) {
	model := New{{ .Name }}(id)
	err := model.Load(r.ps.Get)

	return model, err
}

func (r {{ .Name }}Repository) GetAll() ([]*{{ .Name }}, error) {
	err := r.registry.Load(r.ps.Get)
	if err != nil {
		return nil, err
	}

	result := make([]*{{ .Name }}, len(r.registry))
	cErr := make(chan error)

	for index, id := range r.registry {
		go func(cErr chan error, index int, id string) {
			model, err := r.GetByID(id)
			if err == nil {
				result[index] = model
			}
			
			cErr <- err
		}(cErr, index, id)
	}
	
	for range r.registry {
		if err := <-cErr; err != nil {
			return nil, err
		}
	}

	return result, nil
}

`

func RenderRepository(ra RenderArgs, w io.Writer) {
	tmpl := template.New(ra.Name).Funcs(template.FuncMap{
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
	})

	tmpl = template.Must(tmpl.Parse(repoTmpl))

	tmpl.Execute(w, ra)
}
