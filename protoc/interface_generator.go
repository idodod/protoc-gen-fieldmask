package protoc

import (
	"google.golang.org/protobuf/compiler/protogen"
)

type interfaceGenerator struct {
	name      string
	strFields []*protogen.Field
	msgFields []*protogen.Field
}

func newInterfaceGenerator(message *protogen.Message) *interfaceGenerator {
	return &interfaceGenerator{
		name: getInterfaceName(message),
	}
}

// AddStringFields adds fields for which the fieldmask path is a simple string
func (x *interfaceGenerator) AddStringFields(fields ...*protogen.Field) {
	x.strFields = append(x.strFields, fields...)
}

// AddMessageFields adds fields for which the fieldmask path is a nested message with additional nested paths
func (x *interfaceGenerator) AddMessageFields(fields ...*protogen.Field) {
	x.msgFields = append(x.msgFields, fields...)
}

// Generate generates an interface with all fieldmask paths functions for the given type.
func (x *interfaceGenerator) Generate(g *protogen.GeneratedFile) {
	g.P("type ", x.name, " interface {")
	for _, field := range x.strFields {
		g.P(field.GoName, "() string")
	}
	for _, field := range x.msgFields {
		g.P(field.GoName, "() *", getStructName(field.Message))
	}
	g.P("}")
	g.P()
}
