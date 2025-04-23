package main

import (
	"flag"
	"path"
	"text/template"

	_ "embed"

	"github.com/juliendoutre/protoc-gen-mcp/internal/pb"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

//go:embed main.go.tpl
var mainTemplate string

type Config struct {
	Name    string
	Version string
}

func main() {
	var flags flag.FlagSet

	options := protogen.Options{
		ParamFunc: flags.Set,
	}

	options.Run(func(p *protogen.Plugin) error {
		for _, file := range p.Files {
			if !file.Generate {
				continue
			}

			for _, service := range file.Services {
				option, ok := service.Desc.Options().(*descriptorpb.ServiceOptions)
				if !ok || option == nil {
					continue
				}

				extension, ok := proto.GetExtension(option, pb.E_Server).(*pb.MCPServiceOption)
				if !ok || extension == nil {
					continue
				}

				genFile := p.NewGeneratedFile(path.Join("mcp", service.GoName, "main.go"), "main")

				config := Config{
					Name:    service.GoName,
					Version: extension.GetVersion(),
				}

				for _, method := range service.Methods {
					option, ok := method.Desc.Options().(*descriptorpb.ServiceOptions)
					if !ok || option == nil {
						continue
					}

					extension, ok := proto.GetExtension(option, pb.E_Method).(*pb.MCPMethodOption)
					if !ok || extension == nil {
						continue
					}
				}

				fileTemplate, err := template.New("main").Parse(mainTemplate)
				if err != nil {
					panic(err)
				}

				if err := fileTemplate.ExecuteTemplate(genFile, "main", config); err != nil {
					panic(err)
				}
			}
		}

		return nil
	})
}
