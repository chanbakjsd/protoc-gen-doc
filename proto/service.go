package proto

import (
	"fmt"

	"github.com/chanbakjsd/protoc-gen-doc/doc"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

// ConvertService converts the provided protogen service to a doc service.
func ConvertService(s *protogen.Service) *doc.Service {
	name := string(s.GoName)
	desc := ConvertCommentSet(s.Comments)
	endpoints := make([]*doc.Endpoint, 0, len(s.Methods))
	for _, m := range s.Methods {
		endpoint := ConvertMethod(m)
		if endpoint != nil {
			endpoints = append(endpoints, endpoint)
		}
	}
	return &doc.Service{
		Name:        name,
		Description: desc.Long(name),
		Endpoints:   endpoints,
	}
}

// ConvertMethod converts the provided protogen method to a doc endpoint.
// If the method is not an endpoint, nil is returned instead.
func ConvertMethod(m *protogen.Method) *doc.Endpoint {
	name := string(m.GoName)
	desc := ConvertCommentSet(m.Comments)
	opt := m.Desc.Options()
	req, _ := ConvertMessage(m.Input)
	resp, _ := ConvertMessage(m.Output)
	// Parse HTTP verb and path.
	if !proto.HasExtension(opt, annotations.E_Http) {
		return nil
	}
	rule := proto.GetExtension(opt, annotations.E_Http).(*annotations.HttpRule)
	var method, path string
	switch r := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		method = "GET"
		path = r.Get
	case *annotations.HttpRule_Put:
		method = "PUT"
		path = r.Put
	case *annotations.HttpRule_Post:
		method = "POST"
		path = r.Post
	case *annotations.HttpRule_Delete:
		method = "DELETE"
		path = r.Delete
	case *annotations.HttpRule_Patch:
		method = "PATCH"
		path = r.Patch
	case *annotations.HttpRule_Custom:
		method = r.Custom.Kind
		path = r.Custom.Path
	default:
		panic(fmt.Sprintf("unknown HTTP rule type: %T", rule.Pattern))
	}
	// Find body field.
	var bodyName string
	if rule.Body != "" {
		for _, v := range m.Input.Fields {
			if v.GoName == rule.Body {
				bodyName = v.Desc.JSONName()
				break
			}
		}
		if bodyName == "" {
			panic(fmt.Sprintf(
				"cannot find body field %q for method %q",
				rule.Body, m.Input.GoIdent.GoName,
			))
		}
	}
	return &doc.Endpoint{
		Name:              name,
		Description:       desc.Long(name),
		Method:            method,
		Path:              path,
		BodyField:         bodyName,
		Request:           req,
		Response:          resp,
		StreamingRequest:  m.Desc.IsStreamingClient(),
		StreamingResponse: m.Desc.IsStreamingServer(),
	}
}
