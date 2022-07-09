# protoc-gen-fieldmask

[![CI](https://github.com/idodod/protoc-gen-fieldmask/actions/workflows/ci.yml/badge.svg)](https://github.com/idodod/protoc-gen-fieldmask/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/idodod/protoc-gen-fieldmask)](https://goreportcard.com/report/github.com/idodod/protoc-gen-fieldmask)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/idodod/protoc-gen-fieldmask)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/idodod/protoc-gen-fieldmask)
![GitHub](https://img.shields.io/github/license/idodod/protoc-gen-fieldmask)

A protoc plugin that generates fieldmask paths as static type properties for proto messages, which elimantes the usage of error-prone strings.

For example, given the following proto messages:

```proto

syntax = "proto3";

package example;

option go_package = "example/;example";

import "google/type/date.proto";

message Foo {
  string baz = 1;
  int32 xyz = 2;
  Bar my_bar = 3;
  google.type.Date some_date = 4;
}

message Bar {
  string some_field = 1;
  bool another_field = 2;
}
```

For Golang fieldMasks paths can be used as follows:

```golang
  foo := &example.Foo{}

  fmt.Println(foo.FieldMaskPaths().Baz())
  // Prints "baz"
  
  fmt.Println(foo.FieldMaskPaths().Xyz())
  // Prints "xyz"

  fmt.Println(foo.FieldMaskPaths().MyBar().String())
  // Prints "my_bar"

  // Since baz is a nested message, we can print a nested path:
  fmt.Println(foo.FieldMaskPaths().MyBar().SomeField())
  // Prints "my_bar.some_field"

  // Third party messages work the same way:
  fmt.Println(foo.FieldMaskPaths().SomeDate().String())
  // Prints "some_date"

  fmt.Println(foo.FieldMaskPaths().SomeDate().Year())
  // Prints "some_date.year"
  
  // Full path with message name also available:
  fmt.Println(foo.FullFieldMaskPaths().MyBar().SomeField())
  // Prints "foo.my_bar.some_field"
```

For TypeScript:

```typescript
let foo = new FooFieldMaskPaths();

console.log(foo.fullFieldMaskPaths().MyBar().SomeField())
// Prints "foo.my_bar.some_field"
```

## Usage

### Installation

The plugin can be downloaded from the [release page](https://github.com/idodod/protoc-gen-fieldmask/releases/latest), and should be ideally installed somewhere available in your `$PATH`.

### Executing the plugin

```sh
protoc --fieldmask_out=out_dir protos/example.proto
protoc --fieldmask_out=out_dir --fieldmask_opt=lang=typescript protos/example.proto

# if the plugin is not in your $PATH:
protoc --fieldmask_out=out_dir protos/example.proto --plugin=protoc-gen-fieldmask=/path/to/protoc-gen-fieldmask
```

### Parameters

The following parameters can be set by passing `--fieldmask_opt` to the command:

* `lang`: Language to generate. Supported variants: `go`, `typescript`. Default - `go`.

* `maxdepth`: This option is relevant for a recursive message case.\
    Specify the max depth for which the paths will be pregenerated. If the path depth gets over the max value, it will be generated at runtime.
    default value is `7`.

## Features

*   Currently, supported languages are only `go` and `TypeScript`.
*   All paths are pregenerated (except for recursive messages past `maxdepth`).
*   Support all type of fields including repeated fields, maps, oneofs, third parties, nested messages and recursive messages.
