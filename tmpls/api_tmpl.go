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

func New{{ .Name }}API(repo *{{ .Name }}Repository) *fiber.App {
	api := fiber.New()

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
		result, err := repo.GetByID(c.Params("id"))
		if err != nil {
			return err
		}

		return c.JSON(result)
	})

	api.Patch("/:id", func(c *fiber.Ctx) error {
		model := New{{ .Name }}(c.Params("id"))

		err := c.BodyParser(model)
		if err != nil {
			return err
		}

		err = repo.Edit(model)
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
		var value {{ $t }}
		err := c.BodyParser(&value)
		if err != nil {
			return err
		}

		err = repo.Set{{ $.Name }}{{ $n }}(c.Params("id"), value)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	})

	api.Get("/:id/{{ lower $n | print }}", func(c *fiber.Ctx) error {
		result, err := repo.Get{{ $.Name }}{{ $n }}(c.Params("id"))
		if err != nil {
			return err
		}

		return c.JSON(result)
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
