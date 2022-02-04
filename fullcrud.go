package main

import (
	"log"
	"os"
	"path"

	"github.com/carlosmpv/fullcrud/tmpls"
)

func parseArgs(args ...string) map[string]string {
	dt := map[string]string{}

	for index, arg := range args {
		if index%2 == 0 {
			continue
		}

		dt[arg] = args[index-1]
	}

	return dt
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dt := parseArgs(os.Args[2:]...)

	renderArgs := tmpls.MakeRenderArgs(os.Args[1], path.Base(dir), dt)

	modelsFile, err := os.Create("models.go")
	if err != nil {
		log.Fatal(err)
	}

	tmpls.RenderModel(renderArgs, modelsFile)

	apiFile, err := os.Create("api.go")
	if err != nil {
		log.Fatal(err)
	}

	tmpls.RenderAPI(renderArgs, apiFile)
}
