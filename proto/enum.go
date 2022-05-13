package proto

import (
	"github.com/chanbakjsd/protoc-gen-doc/doc"
	"google.golang.org/protobuf/compiler/protogen"
)

// ConvertEnum converts the provided protogen enum to a doc enum.
func ConvertEnum(e *protogen.Enum) *doc.Enum {
	name := string(e.Desc.Name())
	desc := ConvertCommentSet(e.Comments)
	val := make([]*doc.EnumVal, 0, len(e.Values))
	for _, v := range e.Values {
		val = append(val, ConvertEnumVal(v))
	}
	return &doc.Enum{
		Name:        name,
		Description: desc.Long(name),
		Values:      val,
	}
}

// ConvertEnumVal converts the provided protogen enum value to a doc enum value.
func ConvertEnumVal(v *protogen.EnumValue) *doc.EnumVal {
	name := string(v.Desc.Name())
	desc := ConvertCommentSet(v.Comments)
	return &doc.EnumVal{
		Value:       name,
		Description: desc.Short(name),
	}
}
