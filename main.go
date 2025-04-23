package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/juliendoutre/protoc-gen-mcp/internal/pb"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

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

				genFile := p.NewGeneratedFile(file.GeneratedFilenamePrefix+".txt", file.GoImportPath)
				genFile.P("hello world")

				// TODO: create folder with main.go with skeletton for MCP server

				fmt.Fprintln(os.Stderr, extension.String())

				for _, method := range service.Methods {
					option, ok := method.Desc.Options().(*descriptorpb.ServiceOptions)
					if !ok || option == nil {
						continue
					}

					extension, ok := proto.GetExtension(option, pb.E_Method).(*pb.MCPMethodOption)
					if !ok || extension == nil {
						continue
					}

					// TODO: register method and implement handler

					fmt.Fprintln(os.Stderr, options)
				}
			}
		}
		return nil
	})
}
