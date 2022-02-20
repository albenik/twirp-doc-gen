package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/albenik/twirp-doc-gen/internal/doc"
)

func main() {
	flags := new(flag.FlagSet)
	baseURL := flags.String("base_url", "https://api.example.com/twirp", "")

	pgen := &protogen.Options{
		ParamFunc: flags.Set,
	}
	pgen.Run(func(plugin *protogen.Plugin) error {
		for _, file := range plugin.Files {
			if !file.Generate {
				continue
			}

			for _, service := range file.Services {
				fname := filepath.Join(filepath.Dir(file.GeneratedFilenamePrefix),
					string(service.Desc.Name())+".md")
				f := plugin.NewGeneratedFile(fname, file.GoImportPath)

				if err := doc.NewGenerator(f, *baseURL).GenerateServiceDocument(service); err != nil {
					return fmt.Errorf("%s: schema: %w", file.Desc.Path(), err)
				}
			}
		}

		return nil
	})
}
