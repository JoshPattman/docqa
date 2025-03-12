package docqa

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// EntityJsoner is capable of converting an [Entity] back and forward into a
// json-marshallable object.
type EntityJsoner struct {
	entityFactories  map[string]func() Entity
	entityTypeLookup map[reflect.Type]string
}

// NewEntityJsoner builds a new, empty [EntityJsoner]
// (no types of entity are registered).
func NewEntityJsoner() *EntityJsoner {
	return &EntityJsoner{
		entityFactories:  make(map[string]func() Entity),
		entityTypeLookup: make(map[reflect.Type]string),
	}
}

// AddFactories registers the given functions that can each create a new [Entity]
// with their provided keys.
func (e *EntityJsoner) AddFactories(factories map[string]func() Entity) {
	for k, v := range factories {
		e.entityFactories[k] = v
		e.entityTypeLookup[reflect.TypeOf(v())] = k
	}
}

// Encode converts an [Entity] into a json-serialisable object.
func (enc *EntityJsoner) Encode(e Entity) (any, error) {
	key, ok := enc.entityTypeLookup[reflect.TypeOf(e)]
	if !ok {
		return nil, fmt.Errorf("unrecognised entity type %T", e)
	}
	content, err := e.MakeContent()
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"type":       key,
		"attributes": e.Attr(),
		"content":    content,
	}, nil

}

// Decode converts a json-serialisable object created with [EntityJsoner.Encode]
// into an [Entity].
func (enc *EntityJsoner) Decode(d any) (Entity, error) {
	dict, ok := d.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("was not a dict")
	}

	// Create the recieving entity
	typeAny, ok := dict["type"]
	if !ok {
		return nil, fmt.Errorf("did not have type")
	}
	typeStr, ok := typeAny.(string)
	if !ok {
		return nil, fmt.Errorf("type was not a string")
	}
	factory, ok := enc.entityFactories[typeStr]
	if !ok {
		return nil, fmt.Errorf("unrecognised enitity type")
	}
	into := factory()

	// Load the attributes (for now we will just encode then decode again)
	attributesAny, ok := dict["attributes"]
	if !ok {
		return nil, fmt.Errorf("no attributes")
	}
	bs, err := json.Marshal(attributesAny)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bs, &into)
	if err != nil {
		return nil, err
	}

	// Load the content
	contentAny, ok := dict["content"]
	if !ok {
		return nil, fmt.Errorf("content did not exist")
	}
	content, ok := contentAny.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("content had wrong type")
	}
	err = into.LoadContent(content)
	if err != nil {
		return nil, err
	}
	return into, nil
}
