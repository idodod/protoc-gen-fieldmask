package protoc

import (
	"go/token"

	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
)

type structGenerator struct {
	name      string
	strFields []*protogen.Field
	msgFields []*protogen.Field
	maxDepth  uint
}

func newStructGenerator(message *protogen.Message, maxDepth uint) *structGenerator {
	return &structGenerator{
		name:     strcase.ToLowerCamel(string(message.Desc.Parent().FullName())) + message.GoIdent.GoName + structSuffix,
		maxDepth: maxDepth,
	}
}

// AddStringFields adds fields for which the fieldmask path is a simple string
func (x *structGenerator) AddStringFields(fields ...*protogen.Field) {
	x.strFields = append(x.strFields, fields...)
}

// AddMessageFields adds fields for which the fieldmask path is a nested message with additional nested paths
func (x *structGenerator) AddMessageFields(fields ...*protogen.Field) {
	x.msgFields = append(x.msgFields, fields...)
}

// safeFieldName generates a valid go identifier for the field name
//
// Additionally the and underscore is appended if the field name is 'fieldPath' or 'prefix'
func safeFieldName(field *protogen.Field) string {
	out := strcase.ToLowerCamel(field.GoName)
	if token.IsKeyword(out) || out == "fieldPath" || out == "prefix" {
		return out + "_"
	}
	return out
}

// Generate generates a struct with all fieldmask paths functions for the given type.
func (x *structGenerator) Generate(g *protogen.GeneratedFile) {
	// generate struct with all fields
	g.P("type ", x.name, " struct {")
	g.P("fieldPath string")
	g.P("prefix string")
	for _, field := range x.strFields {
		g.P(safeFieldName(field), " string")
	}
	for _, field := range x.msgFields {
		g.P(safeFieldName(field), " *", getStructName(field.Message))
	}
	g.P("}")
	g.P()

	// generate ctor
	g.P("func ", "new"+strcase.ToCamel(x.name), "(fieldPath string, maxDepth int) *", x.name, " { ")
	g.P("if maxDepth <= 0 {")
	g.P("return nil")
	g.P("}")
	g.P("prefix := \"\"")
	g.P("if fieldPath != \"\" {")
	g.P("prefix = fieldPath + \".\"")
	g.P("}")
	g.P("return &", x.name, "{")
	g.P("fieldPath: fieldPath,")
	g.P("prefix: prefix,")
	for _, field := range x.strFields {
		g.P(safeFieldName(field), ": prefix + \"", field.Desc.Name(), "\",")
	}
	for _, field := range x.msgFields {
		fieldStructNewFunction := getStructNewFunction(field.Message)
		g.P(safeFieldName(field), ": ", fieldStructNewFunction, "(prefix + \"", field.Desc.Name(), "\", maxDepth - 1),")
	}
	g.P("}")
	g.P("}")
	g.P()

	// generate receiver methods
	g.P("func (x *", x.name, ") String() string { return x.fieldPath }")
	for _, field := range x.strFields {
		g.P("func (x *", x.name, ") ", field.GoName, "() string { return x.", safeFieldName(field), "}")
	}
	for _, field := range x.msgFields {
		varName := safeFieldName(field)
		fieldStructNewFunction := getStructNewFunction(field.Message)
		g.P("func (x *", x.name, ") ", field.GoName, "() *", getStructName(field.Message), " {")
		g.P("if x.", varName, "!= nil {")
		g.P("return x.", varName)
		g.P("}")
		g.P("return ", fieldStructNewFunction, "(x.prefix + \"", field.Desc.Name(), "\",", x.maxDepth, ")")
		g.P("}")
	}
	g.P()
}
