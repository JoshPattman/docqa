package qatypes

import (
	"time"

	"github.com/JoshPattman/docqa"
)

// DateEntity represents a [time.Time].
type DateEntity struct {
	docqa.EntityAttributes
	Date time.Time
}

// MakeContent implements [docqa.Entity].
func (e *DateEntity) MakeContent() (map[string]any, error) {
	return map[string]any{
		"date": e.Date.Format(time.DateOnly),
	}, nil
}

// LoadContent implements [docqa.Entity].
func (e *DateEntity) LoadContent(dict map[string]any) error {
	dateStr, err := get[string](dict, "date")
	if err != nil {
		return err
	}
	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return err
	}
	e.Date = date
	return nil
}

// DateType defines a type to create [DateEntity].
type DateType struct{}

// NewDateType creates a [DateType].
func NewDateType() *DateType {
	return &DateType{}
}

// Parse implements [docqa.Type].
func (p *DateType) Parse(value map[string]any) (docqa.Entity, error) {
	var year, month, day float64
	var err error
	if year, err = get[float64](value, "year"); err != nil {
		return nil, err
	}
	if month, err = get[float64](value, "month"); err != nil {
		return nil, err
	}
	if day, err = get[float64](value, "day"); err != nil {
		return nil, err
	}
	return &DateEntity{
		Date: time.Date(int(year), time.Month(int(month)), int(day), 0, 0, 0, 0, time.UTC),
	}, nil
}

// SchemaProperties implements [docqa.Type].
func (p *DateType) SchemaProperties() map[string]any {
	return map[string]any{
		"year": map[string]any{
			"type": "integer",
		},
		"month": map[string]any{
			"type": "integer",
		},
		"day": map[string]any{
			"type": "integer",
		},
	}
}

// Instructions implements [docqa.Type].
func (p *DateType) Instructions() docqa.TypeInstructions {
	return docqa.TypeInstructions{
		OneLiner: "A date, split into year, month, and day",
		Details: []string{
			"If you cannot determine either day or month, assume the 1st / January",
		},
	}
}
