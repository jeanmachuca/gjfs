package schema

import (
	"encoding/json"
	"fmt"
)

// Schema represents a JSON Schema
type Schema struct {
	Type                 interface{}            `json:"type,omitempty"`
	Title                string                 `json:"title,omitempty"`
	Description          string                 `json:"description,omitempty"`
	Default              interface{}            `json:"default,omitempty"`
	Enum                 []interface{}          `json:"enum,omitempty"`
	Const                interface{}            `json:"const,omitempty"`
	Properties           map[string]*Schema     `json:"properties,omitempty"`
	PatternProperties    map[string]*Schema     `json:"patternProperties,omitempty"`
	AdditionalProperties interface{}            `json:"additionalProperties,omitempty"`
	Required             []string               `json:"required,omitempty"`
	Items                *Schema                `json:"items,omitempty"`
	PrefixItems          []*Schema              `json:"prefixItems,omitempty"`
	Contains             *Schema                `json:"contains,omitempty"`
	MinContains          *int                   `json:"minContains,omitempty"`
	MaxContains          *int                   `json:"maxContains,omitempty"`
	MinItems             *int                   `json:"minItems,omitempty"`
	MaxItems             *int                   `json:"maxItems,omitempty"`
	UniqueItems          *bool                  `json:"uniqueItems,omitempty"`
	MinLength            *int                   `json:"minLength,omitempty"`
	MaxLength            *int                   `json:"maxLength,omitempty"`
	Pattern              string                 `json:"pattern,omitempty"`
	Format               string                 `json:"format,omitempty"`
	Minimum              *float64               `json:"minimum,omitempty"`
	ExclusiveMinimum     interface{}            `json:"exclusiveMinimum,omitempty"`
	Maximum              *float64               `json:"maximum,omitempty"`
	ExclusiveMaximum     interface{}            `json:"exclusiveMaximum,omitempty"`
	MultipleOf           *float64               `json:"multipleOf,omitempty"`
	AllOf                []*Schema              `json:"allOf,omitempty"`
	AnyOf                []*Schema              `json:"anyOf,omitempty"`
	OneOf                []*Schema              `json:"oneOf,omitempty"`
	Not                  *Schema                `json:"not,omitempty"`
	If                   *Schema                `json:"if,omitempty"`
	Then                 *Schema                `json:"then,omitempty"`
	Else                 *Schema                `json:"else,omitempty"`
	Definitions          map[string]*Schema     `json:"definitions,omitempty"`
	Defs                 map[string]*Schema     `json:"$defs,omitempty"`
	Ref                  string                 `json:"$ref,omitempty"`
	Schema               string                 `json:"$schema,omitempty"`
	ID                   string                 `json:"$id,omitempty"`
	Comment              string                 `json:"$comment,omitempty"`
	Examples             []interface{}          `json:"examples,omitempty"`
	Deprecated           *bool                  `json:"deprecated,omitempty"`
	ReadOnly             *bool                  `json:"readOnly,omitempty"`
	WriteOnly            *bool                  `json:"writeOnly,omitempty"`
	ContentEncoding      string                 `json:"contentEncoding,omitempty"`
	ContentMediaType     string                 `json:"contentMediaType,omitempty"`
	ContentSchema        *Schema                `json:"contentSchema,omitempty"`
	UnevaluatedItems     *Schema                `json:"unevaluatedItems,omitempty"`
	UnevaluatedProperties *Schema               `json:"unevaluatedProperties,omitempty"`
	PropertyNames        *Schema                `json:"propertyNames,omitempty"`
	MinProperties        *int                   `json:"minProperties,omitempty"`
	MaxProperties        *int                   `json:"maxProperties,omitempty"`
	Dependencies         map[string]interface{} `json:"dependencies,omitempty"`
	DependentRequired    map[string][]string    `json:"dependentRequired,omitempty"`
	DependentSchemas     map[string]*Schema     `json:"dependentSchemas,omitempty"`
}

// ParseSchema parses a JSON schema from bytes
func ParseSchema(data []byte) (*Schema, error) {
	var schema Schema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}
	return &schema, nil
}

// ParseSchemaFromString parses a JSON schema from a string
func ParseSchemaFromString(s string) (*Schema, error) {
	return ParseSchema([]byte(s))
}

// HasRef checks if the schema has a $ref
func (s *Schema) HasRef() bool {
	return s.Ref != ""
}

// GetRef returns the reference path
func (s *Schema) GetRef() string {
	return s.Ref
}

// IsRequired checks if a property is required
func (s *Schema) IsRequired(prop string) bool {
	for _, req := range s.Required {
		if req == prop {
			return true
		}
	}
	return false
}

// GetDefault returns the default value if set
func (s *Schema) GetDefault() interface{} {
	return s.Default
}

// HasDefault checks if the schema has a default value
func (s *Schema) HasDefault() bool {
	return s.Default != nil
}

// GetExample returns the first example if available
func (s *Schema) GetExample() interface{} {
	if len(s.Examples) > 0 {
		return s.Examples[0]
	}
	return nil
}

// HasExamples checks if the schema has examples
func (s *Schema) HasExamples() bool {
	return len(s.Examples) > 0
}

// GetFormat returns the format string
func (s *Schema) GetFormat() string {
	return s.Format
}

// GetPattern returns the pattern string
func (s *Schema) GetPattern() string {
	return s.Pattern
}

// GetMinLength returns the minimum length
func (s *Schema) GetMinLength() int {
	if s.MinLength != nil {
		return *s.MinLength
	}
	return 0
}

// GetMaxLength returns the maximum length
func (s *Schema) GetMaxLength() int {
	if s.MaxLength != nil {
		return *s.MaxLength
	}
	return -1
}

// GetMinimum returns the minimum value
func (s *Schema) GetMinimum() float64 {
	if s.Minimum != nil {
		return *s.Minimum
	}
	return 0
}

// GetMaximum returns the maximum value
func (s *Schema) GetMaximum() float64 {
	if s.Maximum != nil {
		return *s.Maximum
	}
	return 0
}

// GetMinItems returns the minimum items for arrays
func (s *Schema) GetMinItems() int {
	if s.MinItems != nil {
		return *s.MinItems
	}
	return 0
}

// GetMaxItems returns the maximum items for arrays
func (s *Schema) GetMaxItems() int {
	if s.MaxItems != nil {
		return *s.MaxItems
	}
	return -1
}

// GetMinProperties returns the minimum properties for objects
func (s *Schema) GetMinProperties() int {
	if s.MinProperties != nil {
		return *s.MinProperties
	}
	return 0
}

// GetMaxProperties returns the maximum properties for objects
func (s *Schema) GetMaxProperties() int {
	if s.MaxProperties != nil {
		return *s.MaxProperties
	}
	return -1
}

// GetEnum returns the enum values
func (s *Schema) GetEnum() []interface{} {
	return s.Enum
}

// GetConst returns the const value
func (s *Schema) GetConst() interface{} {
	return s.Const
}

// GetType returns the primary type as a string
func (s *Schema) GetType() string {
	if s.Type == nil {
		return ""
	}
	switch t := s.Type.(type) {
	case string:
		return t
	case []interface{}:
		if len(t) > 0 {
			if str, ok := t[0].(string); ok {
				return str
			}
		}
	}
	return ""
}

// GetTypes returns all possible types (handles array of types)
func (s *Schema) GetTypes() []string {
	if s.Type == nil {
		return nil
	}
	switch t := s.Type.(type) {
	case string:
		return []string{t}
	case []interface{}:
		types := make([]string, 0, len(t))
		for _, v := range t {
			if str, ok := v.(string); ok {
				types = append(types, str)
			}
		}
		return types
	}
	return nil
}

// Validate validates data against the schema
func (s *Schema) Validate(data interface{}) error {
	return s.validate(data, "")
}

func (s *Schema) validate(data interface{}, path string) error {
	if s == nil {
		return nil
	}

	if s.HasRef() {
		return fmt.Errorf("%s: $ref validation not implemented", path)
	}

	if types := s.GetTypes(); len(types) > 0 {
		if err := s.validateType(data, path, types); err != nil {
			return err
		}
	}

	if s.Const != nil {
		if !equal(data, s.Const) {
			return fmt.Errorf("%s: value does not match const", path)
		}
	}

	if len(s.Enum) > 0 {
		found := false
		for _, e := range s.Enum {
			if equal(data, e) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%s: value not in enum", path)
		}
	}

	switch s.GetType() {
	case "string":
		if err := s.validateString(data, path); err != nil {
			return err
		}
	case "number", "integer":
		if err := s.validateNumber(data, path); err != nil {
			return err
		}
	case "array":
		if err := s.validateArray(data, path); err != nil {
			return err
		}
	case "object":
		if err := s.validateObject(data, path); err != nil {
			return err
		}
	}

	for i, sub := range s.AllOf {
		if err := sub.validate(data, fmt.Sprintf("%s/allOf/%d", path, i)); err != nil {
			return err
		}
	}
	if len(s.AnyOf) > 0 {
		valid := false
		for i, sub := range s.AnyOf {
			if err := sub.validate(data, fmt.Sprintf("%s/anyOf/%d", path, i)); err == nil {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("%s: does not match anyOf", path)
		}
	}
	if len(s.OneOf) > 0 {
		validCount := 0
		for i, sub := range s.OneOf {
			if err := sub.validate(data, fmt.Sprintf("%s/oneOf/%d", path, i)); err == nil {
				validCount++
			}
		}
		if validCount != 1 {
			return fmt.Errorf("%s: matches %d of oneOf, expected 1", path, validCount)
		}
	}
	if s.Not != nil {
		if err := s.Not.validate(data, path+"/not"); err == nil {
			return fmt.Errorf("%s: matches not schema", path)
		}
	}

	if s.If != nil {
		if err := s.If.validate(data, path+"/if"); err == nil {
			if s.Then != nil {
				if err := s.Then.validate(data, path+"/then"); err != nil {
					return err
				}
			}
		} else if s.Else != nil {
			if err := s.Else.validate(data, path+"/else"); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Schema) validateType(data interface{}, path string, types []string) error {
	for _, t := range types {
		if s.checkType(data, t) {
			return nil
		}
	}
	return fmt.Errorf("%s: type mismatch, expected one of %v, got %T", path, types, data)
}

func (s *Schema) checkType(data interface{}, t string) bool {
	switch t {
	case "string":
		_, ok := data.(string)
		return ok
	case "number":
		switch data.(type) {
		case float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return true
		}
		return false
	case "integer":
		switch data.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return true
		case float64:
			if f := data.(float64); f == float64(int64(f)) {
				return true
			}
		}
		return false
	case "boolean":
		_, ok := data.(bool)
		return ok
	case "array":
		_, ok := data.([]interface{})
		return ok
	case "object":
		_, ok := data.(map[string]interface{})
		return ok
	case "null":
		return data == nil
	}
	return false
}

func (s *Schema) validateString(data interface{}, path string) error {
	str, ok := data.(string)
	if !ok {
		return nil
	}
	if s.MinLength != nil && len(str) < *s.MinLength {
		return fmt.Errorf("%s: string too short (min %d)", path, *s.MinLength)
	}
	if s.MaxLength != nil && len(str) > *s.MaxLength {
		return fmt.Errorf("%s: string too long (max %d)", path, *s.MaxLength)
	}
	return nil
}

func (s *Schema) validateNumber(data interface{}, path string) error {
	var num float64
	switch v := data.(type) {
	case float64:
		num = v
	case int:
		num = float64(v)
	case int64:
		num = float64(v)
	default:
		return nil
	}

	if s.Minimum != nil && num < *s.Minimum {
		return fmt.Errorf("%s: value %f less than minimum %f", path, num, *s.Minimum)
	}
	if s.Maximum != nil && num > *s.Maximum {
		return fmt.Errorf("%s: value %f greater than maximum %f", path, num, *s.Maximum)
	}
	if s.ExclusiveMinimum != nil {
		switch v := s.ExclusiveMinimum.(type) {
		case float64:
			if num <= v {
				return fmt.Errorf("%s: value %f not greater than exclusive minimum %f", path, num, v)
			}
		case bool:
			if v && s.Minimum != nil && num <= *s.Minimum {
				return fmt.Errorf("%s: value %f not greater than exclusive minimum %f", path, num, *s.Minimum)
			}
		}
	}
	if s.ExclusiveMaximum != nil {
		switch v := s.ExclusiveMaximum.(type) {
		case float64:
			if num >= v {
				return fmt.Errorf("%s: value %f not less than exclusive maximum %f", path, num, v)
			}
		case bool:
			if v && s.Maximum != nil && num >= *s.Maximum {
				return fmt.Errorf("%s: value %f not less than exclusive maximum %f", path, num, *s.Maximum)
			}
		}
	}
	if s.MultipleOf != nil && *s.MultipleOf > 0 {
		if num/(*s.MultipleOf) != float64(int64(num/(*s.MultipleOf))) {
			return fmt.Errorf("%s: value %f not a multiple of %f", path, num, *s.MultipleOf)
		}
	}
	return nil
}

func (s *Schema) validateArray(data interface{}, path string) error {
	arr, ok := data.([]interface{})
	if !ok {
		return nil
	}
	if s.MinItems != nil && len(arr) < *s.MinItems {
		return fmt.Errorf("%s: array too short (min %d)", path, *s.MinItems)
	}
	if s.MaxItems != nil && len(arr) > *s.MaxItems {
		return fmt.Errorf("%s: array too long (max %d)", path, *s.MaxItems)
	}
	if s.UniqueItems != nil && *s.UniqueItems {
		seen := make(map[string]bool)
		for i, item := range arr {
			key := fmt.Sprintf("%v", item)
			if seen[key] {
				return fmt.Errorf("%s: duplicate item at index %d", path, i)
			}
			seen[key] = true
		}
	}
	if s.Items != nil {
		for i, item := range arr {
			if err := s.Items.validate(item, fmt.Sprintf("%s/%d", path, i)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Schema) validateObject(data interface{}, path string) error {
	obj, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}
	if s.MinProperties != nil && len(obj) < *s.MinProperties {
		return fmt.Errorf("%s: too few properties (min %d)", path, *s.MinProperties)
	}
	if s.MaxProperties != nil && len(obj) > *s.MaxProperties {
		return fmt.Errorf("%s: too many properties (max %d)", path, *s.MaxProperties)
	}
	for _, req := range s.Required {
		if _, ok := obj[req]; !ok {
			return fmt.Errorf("%s: missing required property '%s'", path, req)
		}
	}
	for name, value := range obj {
		if propSchema, ok := s.Properties[name]; ok {
			if err := propSchema.validate(value, path+"/"+name); err != nil {
				return err
			}
		} else if s.AdditionalProperties != nil {
			switch addProps := s.AdditionalProperties.(type) {
			case bool:
				if !addProps {
					return fmt.Errorf("%s: additional property '%s' not allowed", path, name)
				}
			case *Schema:
				if err := addProps.validate(value, path+"/"+name); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func equal(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
