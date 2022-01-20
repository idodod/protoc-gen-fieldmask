package protoc

import (
	"google.golang.org/protobuf/compiler/protogen"
)

type generator interface {
	Generate(file *protogen.GeneratedFile)
}
