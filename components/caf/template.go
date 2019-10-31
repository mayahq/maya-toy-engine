package caf

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"
	"time"
)

// ContextTemplate is used to configure the context component
type ContextTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Condition   string `json:"condition"`
}

// PropertyTemplate is used to configure *-property components
type PropertyTemplate struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Group    string `json:"group"`
	Template string `json:"template"`
}

// Fill fills the template with provided data and returns the resulting property
func (pt *PropertyTemplate) Fill(data interface{}) (*Property, error) {
	var buf bytes.Buffer
	ts := time.Now().Unix()
	prop := &Property{
		ID:        pt.ID,
		Name:      pt.Name,
		Group:     pt.Group,
		Timestamp: &ts,
	}

	tmpl, err := template.New("value").Parse(pt.Template)
	if err != nil {
		return nil, err
	}

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	switch pt.Type {
	case PropTypeString:
		v := buf.String()
		prop.StringValue = &v

	case PropTypeFloat:
		v, err := strconv.ParseFloat(buf.String(), 64)
		if err != nil {
			return nil, err
		}
		prop.Value = &v

	case PropTypeBool:
		v, err := strconv.ParseBool(buf.String())
		if err != nil {
			return nil, err
		}
		prop.BoolValue = &v

	default:
		return nil, fmt.Errorf("Unsupported property type: %v", pt)
	}
	return prop, nil
}
