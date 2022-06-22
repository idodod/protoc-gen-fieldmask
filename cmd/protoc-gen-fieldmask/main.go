package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andrey308/protoc-gen-fieldmask/protoc"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	defaultMaxDepth = 7
	defaultLang     = "go"
)

var (
	supported = []string{"go", "typescript"}
	version   = "dev"
)

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
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		if !slices.Contains(supported, strings.ToLower(*lang)) {
			return fmt.Errorf("unsupported lang: %s, supported languages: go, typescript", *lang)
		}
		if *maxDepth <= 0 {
			return errors.New("maxdepth must be bigger than 0")
		}
		return protoc.Generate[*lang](plugin, *maxDepth)
	})
}
