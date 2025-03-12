package qatypes

import (
	"github.com/JoshPattman/docqa"
)

// NameEntity is an entity representing a name
type NameEntity struct {
	docqa.EntityAttributes
	FirstName string
	LastName  string
}

// MakeContent implements [docqa.Entity].
func (e *NameEntity) MakeContent() (map[string]any, error) {
	return map[string]any{
		"first_name": e.FirstName,
		"last_name":  e.LastName,
	}, nil
}

// LoadContent implements [docqa.Entity].
func (e *NameEntity) LoadContent(dict map[string]any) error {
	fn, err := get[string](dict, "first_name")
	if err != nil {
		return err
	}
	ln, err := get[string](dict, "last_name")
	if err != nil {
		return err
	}
	e.FirstName, e.LastName = fn, ln
	return nil
}

// NameType defines a type to create [NameEntity].
type NameType struct{}

// NewNameType creates a [NameType].
func NewNameType() *NameType {
	return &NameType{}
}

// Parse implements [docqa.Type].
func (p *NameType) Parse(value map[string]any) (docqa.Entity, error) {
	e := &NameEntity{}
	var err error
	if e.FirstName, err = get[string](value, "first_name"); err != nil {
		return nil, err
	}
	if e.LastName, err = get[string](value, "last_name"); err != nil {
		return nil, err
	}
	e.FirstName = cleanName(e.FirstName)
	e.LastName = cleanName(e.LastName)
	return e, nil
}

// SchemaProperties implements [docqa.Type].
func (p *NameType) SchemaProperties() map[string]any {
	return map[string]any{
		"first_name": map[string]any{
			"type": "string",
		},
		"last_name": map[string]any{
			"type": "string",
		},
	}
}

// Instructions implements [docqa.Type].
func (p *NameType) Instructions() docqa.TypeInstructions {
	return docqa.TypeInstructions{
		OneLiner: "A name of a person, split into first an last name",
		Details: []string{
			"If the person has middle names, include them as further names in the first_name field",
			"If you cannot determine either the first of last name of the person, return an empty string for that property",
		},
	}
}
