package proto

import (
	"strings"

	"github.com/chanbakjsd/protoc-gen-doc/doc"
	"google.golang.org/protobuf/compiler/protogen"
)

// ConvertMessage converts the provided protogen message to a doc message. It
// also returns every nested type.
func ConvertMessage(m *protogen.Message) (*doc.Message, []doc.Type) {
	pkgName := string(m.Desc.ParentFile().Package())
	fullName := string(m.Desc.FullName())
	name := strings.TrimPrefix(fullName, pkgName+".")
	desc := ConvertCommentSet(m.Comments)
	nestedTypes := make([]doc.Type, 0, len(m.Enums)+len(m.Messages))
	for _, e := range m.Enums {
		nestedTypes = append(nestedTypes, ConvertEnum(e))
	}
	for _, nestedMsg := range m.Messages {
		converted, recursedTypes := ConvertMessage(nestedMsg)
		nestedTypes = append(nestedTypes, converted)
		nestedTypes = append(nestedTypes, recursedTypes...)
	}
	fields := make([]*doc.Field, 0, len(m.Fields))
	for _, f := range m.Fields {
		fields = append(fields, ConvertField(f))
	}
	msg := &doc.Message{
		Name:        name,
		Description: desc.Long(m.GoIdent.GoName),
		Fields:      fields,
	}
	return msg, nestedTypes
}

// ConvertField converts the provided protogen field to a doc field.
func ConvertField(f *protogen.Field) *doc.Field {
	jsonName := f.Desc.JSONName()
	desc := ConvertCommentSet(f.Comments)
	return &doc.Field{
		Name:        jsonName,
		GunkName:    f.GoName,
		Description: desc.Short(f.GoName),
		Type:        fieldType(f),
	}
}
