package protoc

import (
	"github.com/idodod/protoc-gen-fieldmask/protoc/golang"
	"github.com/idodod/protoc-gen-fieldmask/protoc/typescript"

	"google.golang.org/protobuf/compiler/protogen"
)

var Generate = map[string]func(*protogen.Plugin, uint) error{
	"go":         golang.Generate,
	"typescript": typescript.Generate,
}
