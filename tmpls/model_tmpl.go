package tmpls

import (
	"html/template"
	"io"
	"strings"
)

const defaultValue = `{{ define "defaultValue" }}{{ if (eq . "bool") }}false{{ end }}{{ if (eq . "int" "int8" "int16" "int32" "int64" "uint" "uint8" "uint16" "uint32" "uint64" "float32" "float64") }}0{{ end }}{{ if (eq . "string")}}""{{ end }}{{ end }}`

const encodeBoolTmpl = `{{ define "encodeBool" }}
	var value []byte
	if v {
		value = []byte{1}
	} else {
		value = []byte{0}
	}{{ end }}`

const encodeNumberTmpl = `{{ define "encodeNumber" }}
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.LittleEndian, v)
	if err != nil {
		return err
	}

	value := buff.Bytes(){{ end }}`

const encodeStringTmpl = `{{ define "encodeString" }}
	value := []byte(v){{ end }}`

const encoderTmpl = `{{ define "encode" }}
{{ if (eq . "bool") }}{{ template "encodeBool" }}{{ end }}
{{ if (eq . "int" "int8" "int16" "int32" "int64" "uint" "uint8" "uint16" "uint32" "uint64" "float32" "float64") }}{{ template "encodeNumber" }}{{ end }}
{{ if (eq . "string")}}{{ template "encodeString" }}{{ end }}
{{ end }}`

const decodeBoolTmpl = `{{ define "decodeBool" }}
	value := v[0] == byte(1){{ end }}`

const decodeNumberTmpl = `{{ define "decodeNumber" }}
	var value {{ . }}
	buff := bytes.NewBuffer(v)
	err = binary.Read(buff, binary.LittleEndian, &value)
	if err != nil {
		return {{ template "defaultValue" . }}, err
	}{{ end }}`

const decodeStringTmpl = `{{ define "decodeString" }}
	value := string(v){{ end }}`

const decoderTmpl = `{{ define "decoder" }}
{{ if (eq . "bool") }}{{ template "decodeBool" }}{{ end }}
{{ if (eq . "int" "int8" "int16" "int32" "int64" "uint" "uint8" "uint16" "uint32" "uint64" "float32" "float64") }}{{ template "decodeNumber" . }}{{ end }}
{{ if (eq . "string")}}{{ template "decodeString" }}{{ end }}
{{ end }}`

const modelTmpl = `
package {{ .Package }}

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
)

var {{ .Name }}RegistryKey = []byte("{{ .Name }}RegistryKey")

type {{ .Name }}Registry []string

func (il *{{ .Name }}Registry) Register(pub func([]byte, []byte) error, id string) error {
	*il = append(*il, id)

	buff := new(bytes.Buffer)
	err := gob.NewEncoder(buff).Encode(il)
	if err != nil {
		return err
	}

	return pub({{ .Name }}RegistryKey, buff.Bytes())
}

func (il *{{ .Name }}Registry) Load(sub func([]byte) ([]byte, error)) error {
	dt, err := sub({{ .Name }}RegistryKey)
	if err != nil {
		return err
	}

	buff := bytes.NewBuffer(dt)
	return gob.NewDecoder(buff).Decode(il)
}

func (il *{{ .Name }}Registry) Delete(pub func([]byte, []byte) error, id string) error {
	copy := {{ .Name }}Registry{}

	for _, iID := range *il {
		if iID != id {
			copy = append(copy, iID)
		}
	}

	*il = copy

	buff := new(bytes.Buffer)
	err := gob.NewEncoder(buff).Encode(il)
	if err != nil {
		return err
	}

	return pub({{ .Name }}RegistryKey, buff.Bytes())
}

func New{{ .Name }}(id string) *{{ .Name }} {
	return &{{ .Name }}{
		ID: id,
		{{ range $n, $t := .Fields }}
		{{ $n }}: {{ template "defaultValue" $t }},{{ end }}
	}
}

type {{ .Name }} struct {
	ID string ` + "`" + `json:"id"` + "`" + `
	{{ range $n, $t := .Fields }}
	{{ $n }} {{ $t }} ` + "`" + `json:"{{ lower $n | print }}"` + "`" + `{{ end }}
}

func (m *{{ .Name }}) Save(pub func([]byte, []byte) error) error {
	var err error
	{{ range $n, $t := .Fields }}
	err = m.Set{{ $n }}(pub, m.{{ $n }})
	if err != nil {
		return err
	}
	{{ end }}

	return nil
}

func (m *{{ .Name }}) Load(sub func([]byte) ([]byte, error)) error {
	var err error
	{{ range $n, $t := .Fields }}
	_, err = m.Get{{ $n }}(sub)
	if err != nil {
		return err
	}{{ end }}

	return nil
}

{{ range $n, $t := .Fields }}
func (m *{{ $.Name }}) Set{{ $n }}(pub func([]byte, []byte) error, v {{ $t }}) error {
	key := []byte(fmt.Sprintf("%s:%s/%s", "{{ $.Name }}", m.ID, "{{ $n }}"))
	{{ template "encode" $t }}

	m.{{ $n }} = v
	return pub(key, value)
}

func (m *{{ $.Name }}) Get{{ $n }}(sub func([]byte) ([]byte, error)) ({{ $t }}, error) {
	key := []byte(fmt.Sprintf("%s:%s/%s", "{{ $.Name }}", m.ID, "{{ $n }}"))
	v, err := sub(key)
	if err != nil {
		return {{ template "defaultValue" $t }}, err
	}

	if v == nil {
		return {{ template "defaultValue" $t }}, err
	}

	{{ template "decoder" $t }}
	m.{{ $n }} = value
	return value, nil
}
{{ end }}
`

func MakeRenderArgs(name, pkg string, fields map[string]string) RenderArgs {
	return RenderArgs{pkg, name, fields}
}

type RenderArgs struct {
	Package string
	Name    string
	Fields  map[string]string
}

func RenderModel(ra RenderArgs, w io.Writer) {
	tmpl := template.New(ra.Name).Funcs(template.FuncMap{
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
	})

	tmpl = template.Must(tmpl.Parse(defaultValue))

	tmpl = template.Must(tmpl.Parse(encodeNumberTmpl))
	tmpl = template.Must(tmpl.Parse(encodeBoolTmpl))
	tmpl = template.Must(tmpl.Parse(encodeStringTmpl))
	tmpl = template.Must(tmpl.Parse(strings.ReplaceAll(encoderTmpl, "\n", "")))

	tmpl = template.Must(tmpl.Parse(decodeBoolTmpl))
	tmpl = template.Must(tmpl.Parse(decodeNumberTmpl))
	tmpl = template.Must(tmpl.Parse(decodeStringTmpl))
	tmpl = template.Must(tmpl.Parse(strings.ReplaceAll(decoderTmpl, "\n", "")))
	tmpl = template.Must(tmpl.Parse(modelTmpl))

	tmpl.Execute(w, ra)
}
