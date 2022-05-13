package proto

import (
	"github.com/chanbakjsd/protoc-gen-doc/doc"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// wellKnownTypes are types defined that are common but not a basic type.
var wellKnownTypes = map[string]doc.Type{
	"google.protobuf.Any":       &doc.Basic{"Any", ""},
	"google.protobuf.Duration":  &doc.Basic{"Duration", ""},
	"google.protobuf.Empty":     &doc.Basic{"Empty", ""},
	"google.protobuf.Value":     &doc.Basic{"JSON", ""},
	"google.protobuf.List":      &doc.Basic{"JSON List", ""},
	"google.protobuf.Struct":    &doc.Basic{"JSON Struct", ""},
	"google.protobuf.Timestamp": &doc.Basic{"Timestamp", ""},
}

// fieldType returns the type of a field.
func fieldType(f *protogen.Field) doc.Type {
	if f.Desc.IsMap() {
		key := fieldType(f.Message.Fields[0])
		value := fieldType(f.Message.Fields[1])
		return &doc.Map{
			Key:   key,
			Value: value,
		}
	}
	var typ doc.Type
	switch f.Desc.Kind() {
	case protoreflect.BoolKind:
		typ = &doc.Basic{"Boolean", ""}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind,
		protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind:
		typ = &doc.Basic{"Integer", ""}
	case protoreflect.Uint32Kind:
		typ = &doc.Basic{"Unsigned Integer", ""}
	case protoreflect.Int64Kind, protoreflect.Sint64Kind,
		protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:
		typ = &doc.Basic{"Integer(64)", ""}
	case protoreflect.Uint64Kind:
		typ = &doc.Basic{"Unsigned Integer(64)", ""}
	case protoreflect.FloatKind:
		typ = &doc.Basic{"Float(32)", ""}
	case protoreflect.DoubleKind:
		typ = &doc.Basic{"Float(64)", ""}
	case protoreflect.StringKind:
		typ = &doc.Basic{"String", ""}
	case protoreflect.BytesKind:
		typ = &doc.Basic{"Bytes", ""}
	case protoreflect.EnumKind:
		typ = &doc.Ref{string(f.Desc.Enum().FullName())}
	case protoreflect.MessageKind:
		fullName := string(f.Desc.Message().FullName())
		typ = &doc.Ref{fullName}
		if t, ok := wellKnownTypes[fullName]; ok {
			typ = t
		}
	default:
		panic("Unknown protobuf type: " + f.Desc.Kind().String())
	}
	if f.Desc.IsList() {
		typ = &doc.Array{typ}
	}
	return typ
}
