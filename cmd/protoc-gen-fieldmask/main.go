package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andrey308/protoc-gen-fieldmask/protoc"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	defaultMaxDepth = 7
	defaultLang     = "go"
)

var version = "dev"

func main() {
	app := filepath.Base(os.Args[0])
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("%s %v\n", app, version)
		return
	}

	var flags flag.FlagSet
	maxDepth := flags.Uint("maxdepth", defaultMaxDepth, "")
	lang := flags.String("lang", defaultLang, "")
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		if strings.ToLower(*lang) != defaultLang {
			return errors.New("go is the only supported language at the moment")
		}
		if *maxDepth <= 0 {
			return errors.New("maxdepth must be bigger than 0")
		}
		return protoc.Generate(plugin, *maxDepth)
	})
}
