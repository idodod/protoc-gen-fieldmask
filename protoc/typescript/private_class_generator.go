package typescript

import (
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
)

type privateClassGenerator struct {
	name      string
	path      string
	strFields []*protogen.Field
	msgFields []*protogen.Field
	maxDepth  uint
}

func newPrivateClassGenerator(message *protogen.Message, maxDepth uint) *privateClassGenerator {
	return &privateClassGenerator{
		name:     strcase.ToLowerCamel(string(message.Desc.Parent().FullName())) + message.GoIdent.GoName + classSuffix,
		path:     string(message.Desc.Parent().FullName()),
		maxDepth: maxDepth,
	}
}

// AddStringFields adds fields for which the fieldmask path is a simple string
func (x *privateClassGenerator) AddStringFields(fields ...*protogen.Field) {
	x.strFields = append(x.strFields, fields...)
}

// AddMessageFields adds fields for which the fieldmask path is a nested message with additional nested paths
func (x *privateClassGenerator) AddMessageFields(fields ...*protogen.Field) {
	x.msgFields = append(x.msgFields, fields...)
}

// Generate generates a struct with all fieldmask paths functions for the given type.
func (x *privateClassGenerator) Generate(g *protogen.GeneratedFile) {
	// generate fields
	g.P("class ", x.name, " {")
	g.P("\t#fieldPath: string;")
	g.P("\t#prefix: string;")
	for _, field := range x.strFields {
		g.P("\t#", strcase.ToLowerCamel(field.GoName), ": string;")
	}
	for _, field := range x.msgFields {
		g.P("\t#", strcase.ToLowerCamel(field.GoName), ": ", getPrivateClassName(field.Message))
	}
	g.P()

	// generate constructor
	g.P("\tconstructor(fieldPath: string, maxDepth: number) {")
	g.P("\t\tif (maxDepth <= 0) {")
	g.P("\t\t\treturn;")
	g.P("\t\t}")
	g.P("\t\tlet prefix = \"\";")
	g.P("\t\tif (fieldPath != \"\") {")
	g.P("\t\t\tprefix = fieldPath + \".\";")
	g.P("\t\t}")
	g.P()

	g.P("\t\tthis.#fieldPath = fieldPath;")
	g.P("\t\tthis.#prefix = prefix;")
	for _, field := range x.strFields {
		g.P("\t\tthis.#", strcase.ToLowerCamel(field.GoName), " = prefix + \"", field.Desc.Name(), "\";")
	}
	for _, field := range x.msgFields {
		fieldStructNewFunction := getPrivateClassName(field.Message)
		g.P("\t\tthis.#", strcase.ToLowerCamel(field.GoName), " = new ", fieldStructNewFunction, "(prefix + \"", field.Desc.Name(), "\", maxDepth - 1);")
	}
	g.P("\t}")
	g.P()

	// generate receiver methods
	g.P("\tString(): string { return this.#fieldPath }")
	for _, field := range x.strFields {
		g.P("\t", field.GoName, "(): string { return this.#", strcase.ToLowerCamel(field.GoName), " }")
	}
	for _, field := range x.msgFields {
		varName := strcase.ToLowerCamel(field.GoName)
		g.P("\t", field.GoName, "(): ", getPrivateClassName(field.Message), " { return this.#", varName, " }")
	}
	g.P("}")
	g.P()
}
