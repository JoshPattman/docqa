package qatypes

import (
	"github.com/JoshPattman/docqa"
)

// TextEntity is an entity representing a text field
type TextEntity struct {
	docqa.EntityAttributes
	Text string
}

// MakeContent implements [docqa.Entity].
func (e *TextEntity) MakeContent() (map[string]any, error) {
	return map[string]any{
		"text": e.Text,
	}, nil
}

// LoadContent implements [docqa.Entity].
func (e *TextEntity) LoadContent(dict map[string]any) error {
	t, err := get[string](dict, "text")
	if err != nil {
		return err
	}
	e.Text = t
	return nil
}

// TextType defines a type to create [TextEntity].
type TextType struct{}

// NewTextType creates a [TextType].
func NewTextType() *TextType {
	return &TextType{}
}

// Parse implements [docqa.Type].
func (t *TextType) Parse(value map[string]any) (docqa.Entity, error) {
	e := &TextEntity{}
	var err error
	if e.Text, err = get[string](value, "text"); err != nil {
		return nil, err
	}
	return e, nil
}

// SchemaProperties implements [docqa.Type].
func (t *TextType) SchemaProperties() map[string]any {
	return map[string]any{
		"text": map[string]any{
			"type": "string",
		},
	}
}

// Instructions implements [docqa.Type].
func (p *TextType) Instructions() docqa.TypeInstructions {
	return docqa.TypeInstructions{
		OneLiner: "A plaintext field",
		Details: []string{
			"Remember that this is json so you will need to json escape the text",
		},
	}
}
