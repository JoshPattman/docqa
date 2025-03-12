package docqa

// Range defines a range of characters from the source document.
type Range struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// IsIndef checks whether the specific [Range] in indefinite.
func (r Range) IsIndef() bool {
	return r == IndefRange()
}

// Len checks how far the [Range] spans.
func (r Range) Len() int {
	return r.End - r.Start
}

// IndefRange creates a new [Range] that is indefinite.
func IndefRange() Range {
	return Range{-1, -1}
}

// EntityAttributes define all common attributes that an [Entity] must have.
type EntityAttributes struct {
	EvidenceRanges []Range `json:"evidence_positions"`
	LocalisedRange Range   `json:"localised_range"`
}

// Attr gets the [EntityAttributes] for this [Entity].
func (a *EntityAttributes) Attr() *EntityAttributes {
	return a
}

// Entity describes a single typed piece of information about a document.
type Entity interface {
	// MakeContent converts the content (does not include attributes) and converts them into a map.
	MakeContent() (map[string]any, error)
	// LoadContent loads a map created with [Entity.ContentToDict] into this.
	LoadContent(map[string]any) error
	// Attr returns the [EntityAttributes]..
	Attr() *EntityAttributes
}
