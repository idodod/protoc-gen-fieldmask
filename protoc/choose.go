package protoc

import (
	"github.com/andrey308/protoc-gen-fieldmask/protoc/golang"
	"github.com/andrey308/protoc-gen-fieldmask/protoc/typescript"

	"google.golang.org/protobuf/compiler/protogen"
)

var Generate = map[string]func(*protogen.Plugin, uint) error{
	"go":         golang.Generate,
	"typescript": typescript.Generate,
}
