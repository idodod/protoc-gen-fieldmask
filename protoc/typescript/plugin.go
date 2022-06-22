package typescript

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	generatedExtension = "_pb.fieldmask.ts"
	classSuffix        = "FieldMaskPaths"
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
		g.P()

		pcGenerator := newPublicClassGenerator(maxDepth)
		generators := []generator{pcGenerator}

		for _, message := range f.Messages {
			generators = append(generators, generateFieldMaskPaths(g, message, "", seen, pcGenerator, maxDepth)...)
		}

		for _, gen := range generators {
			gen.Generate(g)
		}
	}
}

// generateFieldMaskPaths generates a FieldMaskPath struct for each proto message which will contain the fieldmask paths
func generateFieldMaskPaths(g *protogen.GeneratedFile, message *protogen.Message, currFieldPath string, seen map[string]struct{}, publicClassGenerator *publicClassGenerator, maxDepth uint) []generator {
	messageName := string(message.Desc.FullName())
	if _, exists := seen[messageName]; exists {
		return nil
	}
	seen[messageName] = struct{}{}

	privateClGenerator := newPrivateClassGenerator(message, maxDepth)

	var generators []generator
	generators = append(generators, privateClGenerator)
	publicClassGenerator.AddMessage(message)

	for _, field := range message.Fields {
		if (field.Desc.Kind() != protoreflect.MessageKind && field.Desc.Kind() != protoreflect.GroupKind) || field.Desc.IsList() || field.Desc.IsMap() {
			privateClGenerator.AddStringFields(field)
		} else {
			privateClGenerator.AddMessageFields(field)
			nextFieldPath := string(field.Desc.Name())
			if currFieldPath != "" {
				nextFieldPath = currFieldPath + "." + nextFieldPath
			}
			generators = append(generators, generateFieldMaskPaths(g, field.Message, nextFieldPath, seen, publicClassGenerator, maxDepth)...)
		}
	}
	return generators
}
