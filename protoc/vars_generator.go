package protoc

import (
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
)

type varsGenerator struct {
	messages []*protogen.Message
	maxDepth uint
}

func newVarsGenerator(maxDepth uint) *varsGenerator {
	return &varsGenerator{
		maxDepth: maxDepth,
	}
}

// AddMessage adds proto message definitions for which to create an instance of generated fieldmaskpath type
func (x *varsGenerator) AddMessage(messages ...*protogen.Message) {
	x.messages = append(x.messages, messages...)
}

// Generate generates the var definitions for all added messages
func (x *varsGenerator) Generate(g *protogen.GeneratedFile) {
	for _, message := range x.messages {
		messageStructName := getStructName(message)
		structNewFunction := getStructNewFunction(message)
		localVarName := "local" + strcase.ToCamel(messageStructName)
		g.P("var ", localVarName, " = ", structNewFunction, "(\"\",", x.maxDepth, ")")
	}
	g.P()
	for _, message := range x.messages {
		messageInterfaceName := getInterfaceName(message)
		messageStructName := getStructName(message)
		localVarName := "local" + strcase.ToCamel(messageStructName)
		g.P("func (x *", message.GoIdent.GoName, ") ", structSuffix, "() ", messageInterfaceName, " {")
		g.P("return ", localVarName)
		g.P("}")
	}
	g.P()
}
