package docqa

// Type describes an object capable of parsing a response to a [Question] into an [Entity].
type Type interface {
	// Parse takes a json object defined by the schema from [Type.SchemaProperties] and parses it into an [Entity].
	Parse(value map[string]any) (Entity, error)
	// Instructions provides the [TypeInstructions] to send to the LLM.
	Instructions() TypeInstructions
	// SchemaProperties provides the `properties` field of a jsonschema object.
	SchemaProperties() map[string]any
}

// TypeInstructions define how to respond with a [Type] for the LLM.
type TypeInstructions struct {
	OneLiner string
	Details  []string
}
