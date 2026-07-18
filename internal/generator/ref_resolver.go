package generator

import (
	"fmt"
	"strings"

	"github.com/jeanmachuca/gjfs/pkg/schema"
)

// RefResolver resolves JSON references
type RefResolver struct {
	definitions map[string]*schema.Schema
}

// NewRefResolver creates a new reference resolver
func NewRefResolver(defs map[string]*schema.Schema) *RefResolver {
	return &RefResolver{
		definitions: defs,
	}
}

// AddDefinition adds a definition to the resolver
func (r *RefResolver) AddDefinition(name string, s *schema.Schema) {
	if r.definitions == nil {
		r.definitions = make(map[string]*schema.Schema)
	}
	r.definitions[name] = s
}

// Resolve resolves a JSON reference
func (r *RefResolver) Resolve(ref string) (*schema.Schema, error) {
	// Handle #/definitions/Name or #/$defs/Name
	ref = strings.TrimPrefix(ref, "#/")
	parts := strings.Split(ref, "/")

	var current interface{} = r.definitions
	for _, part := range parts {
		part = strings.ReplaceAll(part, "~1", "/")
		part = strings.ReplaceAll(part, "~0", "~")

		switch v := current.(type) {
		case map[string]*schema.Schema:
			if sch, ok := v[part]; ok {
				current = sch
			} else {
				return nil, fmt.Errorf("reference not found: %s", ref)
			}
		case *schema.Schema:
			// Navigate into schema properties
			if part == "properties" {
				current = v.Properties
			} else if part == "items" {
				current = v.Items
			} else if part == "prefixItems" {
				current = v.PrefixItems
			} else if part == "contains" {
				current = v.Contains
			} else if part == "additionalProperties" {
				if addProps, ok := v.AdditionalProperties.(*schema.Schema); ok {
					current = addProps
				} else {
					return nil, fmt.Errorf("additionalProperties is not a schema: %s", ref)
				}
			} else if part == "propertyNames" {
				current = v.PropertyNames
			} else if part == "contentSchema" {
				current = v.ContentSchema
			} else if part == "not" {
				current = v.Not
			} else if part == "if" {
				current = v.If
			} else if part == "then" {
				current = v.Then
			} else if part == "else" {
				current = v.Else
			} else if part == "allOf" || part == "anyOf" || part == "oneOf" {
				// These would need index access, simplified for now
				return nil, fmt.Errorf("cannot navigate into %s without index: %s", part, ref)
			} else if strings.HasPrefix(part, "definitions/") || strings.HasPrefix(part, "$defs/") {
				name := strings.TrimPrefix(strings.TrimPrefix(part, "definitions/"), "$defs/")
				if sch, ok := v.Definitions[name]; ok {
					current = sch
				} else if sch, ok := v.Defs[name]; ok {
					current = sch
				} else {
					return nil, fmt.Errorf("definition not found: %s", name)
				}
			} else {
				return nil, fmt.Errorf("cannot navigate into schema: %s", part)
			}
		case []*schema.Schema:
			// Array of schemas (e.g., prefixItems, allOf, etc.)
			idx := 0
			if _, err := fmt.Sscanf(part, "%d", &idx); err == nil {
				if idx < len(v) {
					current = v[idx]
				} else {
					return nil, fmt.Errorf("index out of bounds: %s", ref)
				}
			} else {
				return nil, fmt.Errorf("expected array index for %s: %s", part, ref)
			}
		default:
			return nil, fmt.Errorf("invalid reference path: %s", ref)
		}
	}

	if schema, ok := current.(*schema.Schema); ok {
		return schema, nil
	}
	return nil, fmt.Errorf("reference does not resolve to a schema: %s", ref)
}
