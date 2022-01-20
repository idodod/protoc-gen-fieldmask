package protoc

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	generatedExtension = ".pb.fieldmask.go"
	structSuffix       = "FieldMaskPaths"
)

// Generate will iterate over all given proto files and will generate fieldmask paths functions for each message
func Generate(plugin *protogen.Plugin, maxDepth uint) error {
	seen := make(map[string]map[string]struct{})
	for _, f := range plugin.Files {
		if !f.Generate {
			continue
		}
		m, exists := seen[f.GoImportPath.String()]
		if !exists {
			m = make(map[string]struct{})
			seen[f.GoImportPath.String()] = m
		}
		generateFile(f, plugin, m, maxDepth)
	}
	return nil
}

func generateFile(f *protogen.File, plugin *protogen.Plugin, seen map[string]struct{}, maxDepth uint) {

	if len(f.Messages) > 0 {
		g := plugin.NewGeneratedFile(getFilePath(f), f.GoImportPath)
		g.P(getFileHeaderComment(f.Desc.Path()))
		g.P("package " + f.GoPackageName)
		g.P("")

		varsGenerator := newVarsGenerator(maxDepth)
		generators := []generator{varsGenerator}

		packageName := string(f.GoImportPath)
		for _, message := range f.Messages {
			generators = append(generators, generateFieldMaskPaths(g, packageName, message, "", seen, varsGenerator, maxDepth)...)
		}

		for _, generator := range generators {
			generator.Generate(g)
		}
	}
}

// generateFieldMaskPaths generates a FieldMaskPath struct for each proto message which will contain the fieldmask paths
func generateFieldMaskPaths(g *protogen.GeneratedFile, generatedFileImportPath string, message *protogen.Message, currFieldPath string, seen map[string]struct{}, varsGenerator *varsGenerator, maxDepth uint) []generator {
	if len(message.Fields) == 0 {
		return nil
	}
	messageName := string(message.Desc.FullName())
	if _, exists := seen[messageName]; exists {
		return nil
	}
	seen[messageName] = struct{}{}

	msgStructGenerator := newStructGenerator(message, maxDepth)
	msgInterfaceGenerator := newInterfaceGenerator(message)

	var generators []generator
	// only generate the fieldmask function if the message belongs to the current file's package
	if string(message.GoIdent.GoImportPath) == generatedFileImportPath {
		varsGenerator.AddMessage(message)
		generators = append(generators, msgInterfaceGenerator)
	}
	generators = append(generators, msgStructGenerator)

	for _, field := range message.Fields {
		if (field.Desc.Kind() != protoreflect.MessageKind && field.Desc.Kind() != protoreflect.GroupKind) || field.Desc.IsList() || field.Desc.IsMap() {
			msgInterfaceGenerator.AddStringFields(field)
			msgStructGenerator.AddStringFields(field)
		} else {
			msgInterfaceGenerator.AddMessageFields(field)
			msgStructGenerator.AddMessageFields(field)
			nextFieldPath := string(field.Desc.Name())
			if currFieldPath != "" {
				nextFieldPath = currFieldPath + "." + nextFieldPath
			}
			generators = append(generators, generateFieldMaskPaths(g, generatedFileImportPath, field.Message, nextFieldPath, seen, varsGenerator, maxDepth)...)
		}
	}
	g.P()
	return generators
}
