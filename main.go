package main

import (
	_ "embed"
	"flag"
	"fmt"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"
)

var version = "(unknown)"

//go:embed main.go.tpl
var mainTemplate string

type Config struct {
	SourcePath    string
	ProtocVersion string
	PluginVersion string
	GoPackageName string
	Services      []Service
}

type Service struct {
	Name    string
	Methods []Method
}

type Method struct {
	Name  string
	Input Input
}

type Input struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name      string
	IsPointer bool
	Type      string
}

func main() {
	var flags flag.FlagSet

	options := protogen.Options{
		ParamFunc: flags.Set,
	}

	options.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		return run(plugin)
	})
}

func run(plugin *protogen.Plugin) error {
	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}

		genFile := plugin.NewGeneratedFile(file.GeneratedFilenamePrefix+"_mcp.pb.go", file.GoImportPath)

		config := Config{
			SourcePath:    file.Desc.Path(),
			ProtocVersion: protocVersion(plugin),
			PluginVersion: version,
			GoPackageName: string(file.GoPackageName),
			Services:      []Service{},
		}

		for _, service := range file.Services {
			serviceTemplate := Service{
				Name:    service.GoName,
				Methods: []Method{},
			}

			for _, method := range service.Methods {
				methodTemplate := Method{
					Name: method.GoName,
					Input: Input{
						Name:   method.Input.GoIdent.GoName,
						Fields: []Field{},
					},
				}

				for _, field := range method.Input.Fields {
					fieldTemplate := Field{
						Name: field.GoName,
					}

					fieldTemplate.Type, fieldTemplate.IsPointer = fieldGoType(genFile, field)

					methodTemplate.Input.Fields = append(methodTemplate.Input.Fields, fieldTemplate)
				}

				serviceTemplate.Methods = append(serviceTemplate.Methods, methodTemplate)
			}

			config.Services = append(config.Services, serviceTemplate)
		}

		fileTemplate, err := template.New("main").Parse(mainTemplate)
		if err != nil {
			return fmt.Errorf("parsing template: %w", err)
		}

		if err := fileTemplate.ExecuteTemplate(genFile, "main", config); err != nil {
			return fmt.Errorf("executing template: %w", err)
		}
	}

	return nil
}

func protocVersion(gen *protogen.Plugin) string {
	version := gen.Request.GetCompilerVersion()
	if version == nil {
		return "(unknown)"
	}

	var suffix string
	if s := version.GetSuffix(); s != "" {
		suffix = "-" + s
	}

	return fmt.Sprintf("v%d.%d.%d%s", version.GetMajor(), version.GetMinor(), version.GetPatch(), suffix)
}

//nolint:lll
// Taken from https://github.com/protocolbuffers/protobuf-go/blob/1a3946737f2cd954c76cd035fa6968468f568b6a/cmd/protoc-gen-go/internal_gengo/main.go#L647

//nolint:cyclop,nonamedreturns,varnamelen
func fieldGoType(g *protogen.GeneratedFile, field *protogen.Field) (goType string, pointer bool) {
	pointer = field.Desc.HasPresence()

	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		goType = "bool"
	case protoreflect.EnumKind:
		goType = g.QualifiedGoIdent(field.Enum.GoIdent)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		goType = "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		goType = "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		goType = "int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		goType = "uint64"
	case protoreflect.FloatKind:
		goType = "float32"
	case protoreflect.DoubleKind:
		goType = "float64"
	case protoreflect.StringKind:
		goType = "string"
	case protoreflect.BytesKind:
		goType = "[]byte"
		pointer = false // rely on nullability of slices for presence
	case protoreflect.MessageKind, protoreflect.GroupKind:
		goType = "*" + g.QualifiedGoIdent(field.Message.GoIdent)
		pointer = false // pointer captured as part of the type
	}

	switch {
	case field.Desc.IsList():
		return "[]" + goType, false
	case field.Desc.IsMap():
		keyType, _ := fieldGoType(g, field.Message.Fields[0])
		valType, _ := fieldGoType(g, field.Message.Fields[1])

		return fmt.Sprintf("map[%v]%v", keyType, valType), false
	}

	return goType, pointer
}
