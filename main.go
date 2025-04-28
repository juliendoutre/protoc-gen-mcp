package main

import (
	_ "embed"
	"flag"
	"strings"
	"text/template"

	"github.com/juliendoutre/protoc-gen-mcp/internal/pb"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

var version = "unknown"

//go:embed main.go.tpl
var mainTemplate string

type Config struct {
	Source        string
	PluginVersion string
	GoPackageName string
	Name          string
	Version       string
	Tools         []Tool
}

type Tool struct {
	Name        string
	Description string
}

func main() {
	var flags flag.FlagSet

	options := protogen.Options{
		ParamFunc: flags.Set,
	}

	options.Run(func(plugin *protogen.Plugin) error {
		for _, file := range plugin.Files {
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

				genFile := plugin.NewGeneratedFile(file.GeneratedFilenamePrefix+"_mcp.pb.go", file.GoImportPath)

				config := Config{
					PluginVersion: version,
					Source:        file.Proto.GetSourceCodeInfo().ProtoReflect().Descriptor().ParentFile().Path(),
					GoPackageName: string(file.GoPackageName),
					Name:          service.GoName,
					Version:       extension.GetVersion(),
					Tools:         []Tool{},
				}

				for _, method := range service.Methods {
					option, ok := method.Desc.Options().(*descriptorpb.MethodOptions)
					if !ok || option == nil {
						continue
					}

					extension, ok := proto.GetExtension(option, pb.E_Method).(*pb.MCPMethodOption)
					if !ok || extension == nil {
						continue
					}

					config.Tools = append(config.Tools, Tool{
						Name:        method.GoName,
						Description: strings.TrimSpace(strings.TrimPrefix(method.Comments.Leading.String(), "//")),
					})
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
