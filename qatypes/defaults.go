package qatypes

import "github.com/JoshPattman/docqa"

// GetDefaultTypes returns a list of all the [docqa.Type] in this package,
// keyed by type key.
func GetDefaultTypes() map[string]docqa.Type {
	return map[string]docqa.Type{
		"name": NewNameType(),
		"date": NewDateType(),
		"text": NewTextType(),
	}
}

func factory[T docqa.Entity]() func() docqa.Entity {
	return func() docqa.Entity {
		var t T
		return t
	}
}

// GetDefaultFactories returns a list of factories that create empty [docqa.Entity],
// keyed by entity key.
func GetDefaultFactories() map[string]func() docqa.Entity {
	return map[string]func() docqa.Entity{
		"name": factory[*NameEntity](),
		"date": factory[*DateEntity](),
		"text": factory[*TextEntity](),
	}
}
