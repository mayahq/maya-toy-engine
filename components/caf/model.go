package caf

import (
	"errors"
	"log"
	"strings"

	"github.com/oleksandr/conditions"
)

// Constants
const (
	PropTypeFloat  = "float"
	PropTypeBool   = "bool"
	PropTypeString = "string"
)

// Property describes a context property
type Property struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Group       string   `json:"group"`
	Timestamp   *int64   `json:"t,omitempty"`
	Value       *float64 `json:"v,omitempty"`
	BoolValue   *bool    `json:"bv,omitempty"`
	StringValue *string  `json:"sv,omitempty"`
}

// Context describes a context (primary or secondary)
type Context struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Group       string        `json:"group"`
	Description string        `json:"description"`
	Condition   string        `json:"condition"`
	Determined  bool          `json:"determined"` // Whether the context state is determined
	Matching    bool          `json:"matching"`
	Timestamp   int64         `json:"timestamp"`
	Entries     []interface{} `json:"entries"`
}

// GetValue returns the value of a Property (the one that is set)
func (p *Property) GetValue() (interface{}, error) {
	if p.Value != nil && p.BoolValue != nil && p.StringValue != nil {
		return nil, errors.New("Only one of {Value, BoolValue, StringValue} should be set")
	}
	if p.Value != nil {
		return *p.Value, nil
	} else if p.BoolValue != nil {
		return *p.BoolValue, nil
	} else if p.StringValue != nil {
		return *p.StringValue, nil
	}
	return nil, errors.New("None of {Value, BoolValue, StringValue} is set")
}

// Evaluate evaluates the context against its condition
func (c *Context) Evaluate() (bool, error) {
	// Only evaluate contexts that are determined
	if !c.Determined {
		return false, errors.New("Undetermined context cannot be evaluated")
	}

	// Check whehter the condition is valid
	if c.Condition == "" {
		return false, errors.New("Empty condition")
	}

	p := conditions.NewParser(strings.NewReader(c.Condition))
	cond, err := p.Parse()
	if err != nil {
		return false, errors.New("Invalid condition")
	}

	// Obtain the values vector
	var values []interface{}
	for _, e := range c.Entries {
		var v interface{}
		switch e.(type) {
		case Context:
			if e.(Context).Determined {
				v = e.(Context).Matching
			}
		case Property:
			prop := e.(Property)
			v, err = prop.GetValue()
			if err != nil {
				log.Println("Error obtaining value of a property: ", err.Error())
			}
		default:
			log.Println("Entry of unsuported type")
		}
		if v != nil {
			values = append(values, v)
		}
	}

	// Check if all values are indeed specified
	if len(values) != len(c.Entries) {
		return false, errors.New("Undetermined context with determined=true?")
	}

	// Evaluate the values aganst the condition
	c.Matching, err = conditions.Evaluate(cond, values...)
	if err != nil {
		return false, err
	}
	return c.Matching, nil
}
