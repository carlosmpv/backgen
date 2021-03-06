package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/carlosmpv/backgen/tmpls"
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

	modelsFile, err := os.Create(fmt.Sprintf("%s_models.go", strings.ToLower(renderArgs.Name)))
	if err != nil {
		log.Fatal(err)
	}

	tmpls.RenderModel(renderArgs, modelsFile)

	repoFile, err := os.Create(fmt.Sprintf("%s_repository.go", strings.ToLower(renderArgs.Name)))
	if err != nil {
		log.Fatal(err)
	}

	tmpls.RenderRepository(renderArgs, repoFile)

	apiFile, err := os.Create(fmt.Sprintf("%s_api.go", strings.ToLower(renderArgs.Name)))
	if err != nil {
		log.Fatal(err)
	}

	tmpls.RenderAPI(renderArgs, apiFile)
}
