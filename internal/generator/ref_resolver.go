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
	if r.definitions == nil {
		return nil, fmt.Errorf("no definitions loaded for reference: %s", ref)
	}

	ref = strings.TrimPrefix(ref, "#/")
	parts := strings.Split(ref, "/")

	// Handle #/$defs/Name, #/definitions/Name (common cases)
	if len(parts) >= 2 && isContainer(parts[0]) {
		name := strings.Join(parts[1:], "/")
		name = strings.ReplaceAll(name, "~1", "/")
		name = strings.ReplaceAll(name, "~0", "~")
		if sch, ok := r.definitions[name]; ok {
			return sch, nil
		}
		return nil, fmt.Errorf("definition %q not found in %s container", name, parts[0])
	}

	// Handle nested paths like #/properties/foo/properties/bar
	return r.resolvePath(parts)
}

func isContainer(part string) bool {
	return part == "$defs" || part == "definitions"
}

func (r *RefResolver) resolvePath(parts []string) (*schema.Schema, error) {
	var current interface{} = r.definitions
	for _, part := range parts {
		part = strings.ReplaceAll(part, "~1", "/")
		part = strings.ReplaceAll(part, "~0", "~")

		switch v := current.(type) {
		case map[string]*schema.Schema:
			if sch, ok := v[part]; ok {
				current = sch
			} else {
				return nil, fmt.Errorf("reference not found: %s", strings.Join(parts, "/"))
			}
		case *schema.Schema:
			next, err := navigateSchema(v, part, strings.Join(parts, "/"))
			if err != nil {
				return nil, err
			}
			current = next
		case []*schema.Schema:
			idx := 0
			if _, err := fmt.Sscanf(part, "%d", &idx); err == nil {
				if idx >= 0 && idx < len(v) {
					current = v[idx]
				} else {
					return nil, fmt.Errorf("index %d out of bounds in %s", idx, strings.Join(parts, "/"))
				}
			} else {
				return nil, fmt.Errorf("expected array index for segment %q", part)
			}
		default:
			return nil, fmt.Errorf("invalid reference path: %s", strings.Join(parts, "/"))
		}
	}

	if sch, ok := current.(*schema.Schema); ok {
		return sch, nil
	}
	return nil, fmt.Errorf("reference does not resolve to a schema: %s", strings.Join(parts, "/"))
}

func navigateSchema(s *schema.Schema, part, ref string) (interface{}, error) {
	switch part {
	case "properties":
		return s.Properties, nil
	case "items":
		return s.Items, nil
	case "prefixItems":
		// Will need index next
		return s.PrefixItems, nil
	case "contains":
		return s.Contains, nil
	case "not":
		return s.Not, nil
	case "if":
		return s.If, nil
	case "then":
		return s.Then, nil
	case "else":
		return s.Else, nil
	case "propertyNames":
		return s.PropertyNames, nil
	case "contentSchema":
		return s.ContentSchema, nil
	case "allOf":
		return s.AllOf, nil
	case "anyOf":
		return s.AnyOf, nil
	case "oneOf":
		return s.OneOf, nil
	case "additionalProperties":
		if ap, ok := s.AdditionalProperties.(*schema.Schema); ok {
			return ap, nil
		}
		return nil, fmt.Errorf("additionalProperties is not a schema: %s", ref)
	case "unevaluatedItems":
		return s.UnevaluatedItems, nil
	case "unevaluatedProperties":
		return s.UnevaluatedProperties, nil
	default:
		return nil, fmt.Errorf("cannot navigate into schema segment %q", part)
	}
}
