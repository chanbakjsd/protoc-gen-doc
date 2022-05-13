package doc

// Tag contains the packages for a specific tag, as well as the preamble
// information.
type Tag struct {
	// Name is the name of the tag.
	Name string `json:"name"`
	// Preamble is the preamble for this tag.
	Preamble string `json:"preamble"`
	// Weight is the weight of the tag, used for sorting.
	Weight int `json:"weight"`
	// Packages is the list of packages for this tag.
	Packages []*Package `json:"packages"`
}

// Package is the documentation for a package.
type Package struct {
	// Name is the name of the package. It may not necessarily be unique.
	Name string `json:"name"`
	// ID is the full qualified path of the package and can serve as a unique
	// identifier.
	ID string `json:"id"`
	// Description is the description of the comment.
	Description string `json:"description"`
	// Services is a list of services in the package.
	Services []*Service `json:"services"`
	// Types is a list of data types in the package.
	Types map[string]Type `json:"types"`
}

// Service is the documentation for a service.
type Service struct {
	// Name is the name of the service.
	Name string `json:"name"`
	// Description is the description of the service.
	Description string `json:"description"`
	// Methods is a list of methods in the service.
	Endpoints []*Endpoint `json:"endpoints"`
}

// Endpoint is the documentation for an endpoint.
type Endpoint struct {
	// Name is the name of the endpoint.
	Name string `json:"name"`
	// Description is the description of the endpoint.
	Description string `json:"description"`
	// Method is the HTTP method to trigger the endpoint.
	Method string `json:"method"`
	// Path is the HTTP path to trigger the endpoint.
	Path string `json:"path"`
	// BodyField is the name of the field that contains the request body.
	BodyField string `json:"body_field"`
	// Request is the data type of the request.
	Request Type `json:"request"`
	// Response is the data type of the response.
	Response Type `json:"response"`
	// StreamingRequest is true if the request is streamed.
	StreamingRequest bool `json:"streaming_request"`
	// StreamingResponse is true if the response is streamed.
	StreamingResponse bool `json:"streaming_response"`
}
