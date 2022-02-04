package tmpls

import (
	"html/template"
	"io"
	"strings"
)

const apiTmpl = `
package {{ .Package }}

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"
)

type pubsub interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte) error
}

func New{{ .Name }}API(ps pubsub) *fiber.App {
	api := fiber.New()
	registry := {{ .Name }}Registry{}

	api.Get("/", func(c *fiber.Ctx) error {
		result := []*{{ .Name }}{}
		
		err := registry.Load(ps.Get)
		if err != nil {
			return err
		}

		for _, id := range registry {
			model := New{{ .Name }}(id)
			err = model.Load(ps.Get)
			if err != nil {
				return err
			}

			result = append(result, model)
		}

		return c.JSON(result)
	})
	
	api.Post("/", func(c *fiber.Ctx) error {
		model := New{{ .Name }}(xid.New().String())
		
		err := c.BodyParser(model)
		if err != nil {
			return err
		}

		err = model.Save(ps.Set)
		if err != nil {
			return err
		}

		err = registry.Register(ps.Set, model.ID)
		if err != nil {
			return err
		}

		return c.JSON(model.ID)
	})

	api.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		model := New{{ .Name }}(id)
		err := model.Load(ps.Get)
		if err != nil {
			return err
		}

		return c.JSON(model)
	})

	api.Patch("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		model := New{{ .Name }}(id)
		
		err := c.BodyParser(model)
		if err != nil {
			return err
		}

		err = model.Save(ps.Set)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	})

	api.Delete("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		err := registry.Delete(ps.Set, id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	})

	{{ range $n, $t := .Fields }}
	api.Put("/:id/{{ lower $n | print }}", func(c *fiber.Ctx) error {
		id := c.Params("id")
		model := New{{ $.Name }}(id)

		var value {{ $t }}
		err := c.BodyParser(&value)
		if err != nil {
			return err
		}

		err = model.Set{{ $n }}(ps.Set, value)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	})

	api.Get("/:id/{{ lower $n | print }}", func(c *fiber.Ctx) error {
		id := c.Params("id")
		model := New{{ $.Name }}(id)
		ret, err := model.Get{{ $n }}(ps.Get)
		if err != nil {
			return err
		}

		return c.JSON(ret)
	})
	{{ end }}

	return api
}
`

func RenderAPI(ra RenderArgs, w io.Writer) {
	tmpl := template.New(ra.Name).Funcs(template.FuncMap{
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
	})

	tmpl = template.Must(tmpl.Parse(apiTmpl))
	tmpl.Execute(w, ra)
}
