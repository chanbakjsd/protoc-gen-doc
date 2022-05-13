package proto

import (
	"github.com/chanbakjsd/protoc-gen-doc/doc"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// pkgPath is the path to the comment of a package hard-coded in descriptorpb.
var pkgPath = protoreflect.SourcePath{2}

// ConvertFile converts the provided protogen file to a package.
func ConvertFile(f *protogen.File) *doc.Package {
	name := string(f.GoPackageName)
	path := string(f.Proto.GetPackage())
	desc := ParseDesc(f.Desc.SourceLocations().ByPath(pkgPath).LeadingComments)
	// Convert types.
	typ := make(map[string]doc.Type, len(f.Enums)+len(f.Messages))
	for _, e := range f.Enums {
		conv := ConvertEnum(e)
		typ[conv.Name] = conv
	}
	for _, msg := range f.Messages {
		if msg.Desc.IsMapEntry() {
			continue
		}
		msg, extra := ConvertMessage(msg)
		typ[msg.Name] = msg
		for _, t := range extra {
			t := t.(doc.NamedType)
			typ[t.TypeName()] = t
		}
	}
	services := make([]*doc.Service, 0, len(f.Services))
	for _, s := range f.Services {
		services = append(services, ConvertService(s))
	}
	return &doc.Package{
		Name:        name,
		ID:          path,
		Description: desc.Text,
		Services:    services,
		Types:       typ,
	}
}
