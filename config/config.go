package config

// Config is the configuration of protoc-gen-doc.
type Config struct {
	// Sections is a list of section available.
	Sections map[string]Section
}

// Section is a part of the documentation as defined. Each section will output
// a separate tag.
type Section struct {
	// DisplayName is the name to display on the documentation tag like a
	// header.
	DisplayName string
	// Packages is the list of packages to include for the section.
	Packages []string
	// PreambleContent is the content to display before the struct and endpoint
	// definitions.
	PreambleContent string
	// Weight is an arbitrary number defining the order of the sections. Lower
	// numbers should be placed nearer to the front.
	Weight int
}
