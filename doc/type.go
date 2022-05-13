package doc

import "encoding/json"

// Type is the documentation for a data type.
type Type interface {
	typ() string
}

// NamedType is a data type that is named and not purely derived from other
// data types.
// Arrays and maps are not NamedTypes.
type NamedType interface {
	Type
	TypeName() string
}

// Message is the documentation for a message.
type Message struct {
	// Name is the name of the data type.
	Name string `json:"name"`
	// Description is the description of the data type.
	Description string `json:"description"`
	// Fields is a list of fields in the data type.
	Fields []*Field `json:"fields"`
}

// Field is the documentation for a field in a message.
type Field struct {
	// Name is the name of the field.
	Name string `json:"name"`
	// GunkName is the gunk name of the field.
	GunkName string `json:"-"`
	// Description is the description of the field.
	Description string `json:"description"`
	// Type is the type of the field.
	Type Type `json:"type"`
}

// Enum is the documentation for an enum.
type Enum struct {
	// Name is the name of the data type.
	Name string `json:"name"`
	// Description is the description of the data type.
	Description string `json:"description"`
	// Values are the list of values for enum.
	Values []*EnumVal `json:"values"`
}

// EnumVal is the documentation for enum values.
type EnumVal struct {
	// Value is the value of the enum.
	Value string `json:"value"`
	// Description is the description of the enum value.
	Description string `json:"description"`
}

// Ref is a reference to a Message or Enum type.
type Ref struct {
	// Name is the fully qualified name of the referenced type.
	Name string `json:"name"`
}

// Map is the documentation for a map.
type Map struct {
	// Key is the type of the map key.
	Key Type `json:"key"`
	// Value is the type of the map value.
	Value Type `json:"value"`
}

// Array is the documentation for an array.
type Array struct {
	// Value is the type of the array.
	Value Type `json:"value"`
}

// Basic is the documentation for a basic type such as string or int.
type Basic struct {
	// Name is the name of the data type.
	Name string `json:"name"`
	// Example is an example of the value of the type.
	Example string `json:"example"`
}

func (*Message) typ() string { return "message" }
func (*Enum) typ() string    { return "enum" }
func (*Ref) typ() string     { return "ref" }
func (*Map) typ() string     { return "map" }
func (*Array) typ() string   { return "array" }
func (*Basic) typ() string   { return "basic" }

func (m *Message) TypeName() string { return m.Name }
func (e *Enum) TypeName() string    { return e.Name }
func (r *Ref) TypeName() string     { return r.Name }
func (b *Basic) TypeName() string   { return b.Name }

// MarshalJSON marshals the message to JSON with a type field.
func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Message
	}{
		Type:    "message",
		Message: *m,
	})
}

// MarshalJSON marshals the enum to JSON with a type field.
func (e *Enum) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Enum
	}{
		Type: "enum",
		Enum: *e,
	})
}

// MarshalJSON marshals the enum to JSON with a type field.
func (r *Ref) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Ref
	}{
		Type: "ref",
		Ref:  *r,
	})
}

// MarshalJSON marshals the map to JSON with a type field.
func (m *Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Map
	}{
		Type: "map",
		Map:  *m,
	})
}

// MarshalJSON marshals the array to JSON with a type field.
func (a *Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Array
	}{
		Type:  "array",
		Array: *a,
	})
}

// MarshalJSON marshals the basic type to JSON with a type field.
func (b *Basic) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Basic
	}{
		Type:  "basic",
		Basic: *b,
	})
}
