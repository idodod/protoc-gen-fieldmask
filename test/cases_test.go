package test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/andrey308/protoc-gen-fieldmask/test/gen/cases"
	"github.com/andrey308/protoc-gen-fieldmask/test/gen/cases/a"
	"github.com/andrey308/protoc-gen-fieldmask/test/gen/cases/b"
	"github.com/andrey308/protoc-gen-fieldmask/test/gen/cases/thirdpartyimport"
	"github.com/iancoleman/strcase"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type TestSuite struct {
	suite.Suite
	maxDepth int
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupTest() {
	s.maxDepth = 100
}

func (s *TestSuite) TestTypes() {
	testCases := map[string]struct {
		msg proto.Message
	}{
		"all field types get a fieldmask path":                                                             {msg: &cases.Foo{}},
		"an external nested message also gets fieldmasks when it's a parent message":                       {msg: &cases.TestNestedExternalMessage{}},
		"an internal nested message also gets fieldmasks when it's a parent message":                       {msg: &cases.Foo_TestNestedInternalMessage{}},
		"an external nested message from another file also gets fieldmasks when it's a parent message":     {msg: &cases.YetAnotherTestNestedExternalMessage{}},
		"messages with the same name and a different package get fieldmasks (a.Foo, 1/2)":                  {msg: &a.Foo{}},
		"messages with the same name and a different package get fieldmasks (b.Foo, 2/2)":                  {msg: &b.Foo{}},
		"messages from different proto files, in the same package can get fieldmask for 3rd-parties (1/2)": {msg: &thirdpartyimport.FooA{}},
		"messages from different proto files, in the same package can get fieldmask for 3rd-parties (2/2)": {msg: &thirdpartyimport.FooB{}},
		"recursive message works": {msg: &cases.Node{}},
	}

	for name, testCase := range testCases {
		s.Run(name, func() {
			msgTypeCountMap := make(map[string]int)
			msgPtr := testCase.msg
			msgPtrType := reflect.TypeOf(msgPtr)
			fmPaths, res := msgPtrType.MethodByName("FieldMaskPaths")
			s.Run("FieldMaskPaths exists", func() {
				s.Require().True(res)
			})
			fmVal := fmPaths.Func.Call([]reflect.Value{reflect.ValueOf(msgPtr)})[0]
			ref := fmVal.Type()
			s.Run("FieldMaskPaths does not have a String method", func() {
				_, res := ref.MethodByName("String")
				s.Assert().False(res)
			})

			s.Run("All fields have valid paths", func() {
				paths := s.collectAndAssertPaths("", msgPtrType, fmVal, msgTypeCountMap)
				s.Assert().NotZero(len(paths), "number of paths cannot be zero")
				fm, err := fieldmaskpb.New(msgPtr, paths...)
				s.Require().NoError(err)
				s.T().Log(fm)
			})
		})
	}
}

func (s *TestSuite) collectAndAssertPaths(parent string, protoMessagePtrType reflect.Type, fieldMaskValue reflect.Value, typesMap map[string]int) []string {
	el := protoMessagePtrType.Elem()
	name := el.PkgPath() + "." + el.Name()
	if c, exists := typesMap[name]; exists && c > s.maxDepth {
		return nil
	}
	typesMap[name]++

	var paths []string
	for i := 0; i < protoMessagePtrType.NumMethod(); i++ {
		method := protoMessagePtrType.Method(i)
		if !method.IsExported() || !strings.HasPrefix(method.Name, "Get") {
			continue
		}
		funcType := method.Func.Type()
		if funcType.NumOut() != 1 {
			continue
		}

		out := funcType.Out(0)
		if out.Kind() == reflect.Interface {
			// skipping oneof's
			continue
		}
		fieldMaskMethodName := strings.TrimPrefix(method.Name, "Get")
		m := fieldMaskValue.MethodByName(fieldMaskMethodName)
		s.Assert().False(m.IsZero())
		s.Assert().Zero(m.Type().NumIn())
		s.Assert().Equal(1, m.Type().NumOut())
		outType := m.Type().Out(0)
		expected := strcase.ToSnake(fieldMaskMethodName)
		if parent != "" {
			expected = parent + "." + expected
		}
		fv := m.Call(nil)[0]
		var res string
		if outType.Kind() == reflect.String {
			res = fv.String()
			s.Assert().Equal(expected, res)
		} else if outType.Kind() == reflect.Ptr {
			strMethod, exists := outType.MethodByName("String")
			s.Assert().True(exists, "String method does not exists for %s", fieldMaskMethodName)
			res = strMethod.Func.Call([]reflect.Value{fv})[0].String()
			s.Assert().Equal(expected, res)
			paths = append(paths, s.collectAndAssertPaths(res, out, fv, typesMap)...)
		}
		paths = append(paths, res)
	}
	return paths
}
