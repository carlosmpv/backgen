package tmpls

import (
	"io"
	"strings"
	"text/template"
)

const apiTmpl = `
package {{ .Package }}

import (
	"github.com/gofiber/fiber/v2"
)

type pubsub interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte) error
}

func New{{ .Name }}API(ps pubsub) *fiber.App {
	api := fiber.New()
	repo := New{{ .Name }}Repository(ps)

	api.Get("/", func(c *fiber.Ctx) error {
		result, err := repo.GetAll()
		if err != nil {
			return err
		}

		return c.JSON(result)
	})
	
	api.Post("/", func(c *fiber.Ctx) error {
		model := &{{ .Name }}{}
		
		err := c.BodyParser(model)
		if err != nil {
			return err
		}

		id, err := repo.Create(model)
		if err != nil {
			return err
		}

		return c.JSON(id)
	})

	api.Get("/:id", func(c *fiber.Ctx) error {
		model := New{{ .Name }}(c.Params("id"))
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
		err := repo.Delete(c.Params("id"))
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
