package golang

import (
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
)

type fullPathGenerator struct {
	messages []*protogen.Message
	maxDepth uint
}

func newFullPathGenerator(maxDepth uint) *fullPathGenerator {
	return &fullPathGenerator{
		maxDepth: maxDepth,
	}
}

// AddMessage adds proto message definitions for which to create an instance of generated fieldmaskpath type
func (x *fullPathGenerator) AddMessage(messages ...*protogen.Message) {
	x.messages = append(x.messages, messages...)
}

// Generate generates the var definitions for all added messages
func (x *fullPathGenerator) Generate(g *protogen.GeneratedFile) {
	for _, message := range x.messages {
		messageStructName := getStructName(message)
		structNewFunction := getStructNewFunction(message)
		fullPathVarName := "fullPath" + strcase.ToCamel(messageStructName)
		g.P("var ", fullPathVarName, " = ", structNewFunction, "(\"", getMessageName(message), "\",", x.maxDepth, ")")
	}
	g.P()
	for _, message := range x.messages {
		messageInterfaceName := getInterfaceName(message)
		messageStructName := getStructName(message)
		fullPathVarName := "fullPath" + strcase.ToCamel(messageStructName)
		g.P("func (x *", message.GoIdent.GoName, ") ", fullPathSuffix, "() ", messageInterfaceName, " {")
		g.P("return ", fullPathVarName)
		g.P("}")
	}
	g.P()
}
