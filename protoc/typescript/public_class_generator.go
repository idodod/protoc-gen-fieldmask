package typescript

import (
	"google.golang.org/protobuf/compiler/protogen"
)

type publicClassGenerator struct {
	messages []*protogen.Message
	maxDepth uint
}

func newPublicClassGenerator(maxDepth uint) *publicClassGenerator {
	return &publicClassGenerator{
		maxDepth: maxDepth,
	}
}

// AddMessage adds proto message definitions for which to create an instance of generated fieldmaskpath type
func (x *publicClassGenerator) AddMessage(messages ...*protogen.Message) {
	x.messages = append(x.messages, messages...)
}

// Generate generates the public class definitions for all added messages
func (x *publicClassGenerator) Generate(g *protogen.GeneratedFile) {
	for _, message := range x.messages {
		publicClassName := getPublicClassName(message)
		privateClassName := getPrivateClassName(message)
		g.P("export class ", publicClassName, "{")
		g.P("\tfieldMaskPaths(): ", privateClassName, "{")
		g.P("\t\treturn new ", privateClassName, "(\"\", ", x.maxDepth, ")")
		g.P("\t}")
		g.P("\tfullFieldMaskPaths(): ", privateClassName, "{")
		g.P("\t\treturn new ", privateClassName, "(\"", getMessageName(message), "\",", x.maxDepth, ")")
		g.P("\t}")
		g.P("}")
	}
	g.P()
}
