package generator

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/jeanmachuca/gjfs/pkg/schema"
)

// Generator generates JSON examples from JSON Schema
type Generator struct {
	seed         int64
	strictMode   bool
	useExamples  bool
	useDefaults  bool
	requiredOnly bool
	maxArraySize int
	maxDepth     int
	currentDepth int
	rand         *rand.Rand
	resolver     *RefResolver
}

// GeneratorOption is a functional option for Generator
type GeneratorOption func(*Generator)

// WithSeed sets the random seed for reproducible generation
func WithSeed(seed int64) GeneratorOption {
	return func(g *Generator) {
		g.seed = seed
		g.rand = rand.New(rand.NewSource(seed))
	}
}

// WithStrictMode enables strict mode (no random values, uses defaults/examples only)
func WithStrictMode(strict bool) GeneratorOption {
	return func(g *Generator) {
		g.strictMode = strict
	}
}

// WithExamples enables using examples from schema
func WithExamples(use bool) GeneratorOption {
	return func(g *Generator) {
		g.useExamples = use
	}
}

// WithDefaults enables using default values from schema
func WithDefaults(use bool) GeneratorOption {
	return func(g *Generator) {
		g.useDefaults = use
	}
}

// WithRequiredOnly generates only required properties
func WithRequiredOnly(requiredOnly bool) GeneratorOption {
	return func(g *Generator) {
		g.requiredOnly = requiredOnly
	}
}

// WithMaxArraySize sets the maximum array size
func WithMaxArraySize(size int) GeneratorOption {
	return func(g *Generator) {
		g.maxArraySize = size
	}
}

// WithMaxDepth sets the maximum recursion depth
func WithMaxDepth(depth int) GeneratorOption {
	return func(g *Generator) {
		g.maxDepth = depth
	}
}

// WithDefinitions provides schema definitions for $ref resolution
func WithDefinitions(defs map[string]*schema.Schema) GeneratorOption {
	return func(g *Generator) {
		g.resolver = NewRefResolver(defs)
	}
}

// NewGenerator creates a new JSON example generator
func NewGenerator(opts ...GeneratorOption) *Generator {
	g := &Generator{
		seed:         time.Now().UnixNano(),
		strictMode:   false,
		useExamples:  true,
		useDefaults:  true,
		requiredOnly: false,
		maxArraySize: 5,
		maxDepth:     50,
		rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
		resolver:     NewRefResolver(nil),
	}
	for _, opt := range opts {
		opt(g)
	}
	if g.rand == nil {
		g.rand = rand.New(rand.NewSource(g.seed))
	}
	return g
}

// Generate generates a JSON example from a schema
func (g *Generator) Generate(s *schema.Schema) (interface{}, error) {
	g.currentDepth = 0
	return g.generateValue(s, make(map[string]bool))
}

// GenerateJSON generates a JSON example and returns it as JSON bytes
func (g *Generator) GenerateJSON(s *schema.Schema) ([]byte, error) {
	example, err := g.Generate(s)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(example, "", "  ")
}

// GenerateJSONString generates a JSON example as a formatted string
func (g *Generator) GenerateJSONString(s *schema.Schema) (string, error) {
	data, err := g.GenerateJSON(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (g *Generator) generateValue(s *schema.Schema, visitedRefs map[string]bool) (interface{}, error) {
	if s == nil {
		return nil, nil
	}

	if g.currentDepth > g.maxDepth {
		return nil, fmt.Errorf("maximum recursion depth exceeded")
	}

	// Handle $ref
	if s.HasRef() {
		return g.resolveRef(s.GetRef(), visitedRefs)
	}

	// Handle const
	if s.GetConst() != nil {
		return s.GetConst(), nil
	}

	// Handle enum
	if len(s.GetEnum()) > 0 {
		return g.generateEnum(s.GetEnum()), nil
	}

	// Handle examples
	if g.useExamples && s.HasExamples() && !g.strictMode {
		return s.GetExample(), nil
	}

	// Handle default
	if g.useDefaults && s.HasDefault() {
		return s.GetDefault(), nil
	}

	// Handle combined schemas
	if len(s.AllOf) > 0 {
		return g.generateAllOf(s.AllOf, visitedRefs)
	}
	if len(s.AnyOf) > 0 {
		return g.generateAnyOf(s.AnyOf, visitedRefs)
	}
	if len(s.OneOf) > 0 {
		return g.generateOneOf(s.OneOf, visitedRefs)
	}
	if s.Not != nil {
		return g.generateNot(s.Not, visitedRefs)
	}
	if s.If != nil {
		return g.generateIfThenElse(s.If, s.Then, s.Else, visitedRefs)
	}

	// Generate based on type
	primaryType := s.GetType()
	switch primaryType {
	case "string":
		return g.generateString(s), nil
	case "number", "integer":
		return g.generateNumber(s), nil
	case "boolean":
		return g.generateBoolean(), nil
	case "array":
		return g.generateArray(s, visitedRefs)
	case "object":
		return g.generateObject(s, visitedRefs)
	case "null":
		return nil, nil
	default:
		// Try to infer from properties or items
		if len(s.Properties) > 0 || len(s.PatternProperties) > 0 || s.AdditionalProperties != nil {
			return g.generateObject(s, visitedRefs)
		}
		if s.Items != nil {
			return g.generateArray(s, visitedRefs)
		}
		return g.generateString(s), nil
	}
}

func (g *Generator) resolveRef(ref string, visitedRefs map[string]bool) (interface{}, error) {
	if visitedRefs[ref] {
		return nil, fmt.Errorf("circular reference detected: %s", ref)
	}
	visitedRefs[ref] = true
	defer delete(visitedRefs, ref)

	resolved, err := g.resolver.Resolve(ref)
	if err != nil {
		return nil, err
	}
	return g.generateValue(resolved, visitedRefs)
}

func (g *Generator) generateAllOf(schemas []*schema.Schema, visitedRefs map[string]bool) (interface{}, error) {
	result := make(map[string]interface{})
	for _, sch := range schemas {
		val, err := g.generateValue(sch, visitedRefs)
		if err != nil {
			return nil, err
		}
		if m, ok := val.(map[string]interface{}); ok {
			for k, v := range m {
				result[k] = v
			}
		}
	}
	return result, nil
}

func (g *Generator) generateAnyOf(schemas []*schema.Schema, visitedRefs map[string]bool) (interface{}, error) {
	if len(schemas) == 0 {
		return nil, nil
	}
	for _, sch := range schemas {
		val, err := g.generateValue(sch, visitedRefs)
		if err == nil && val != nil {
			return val, nil
		}
	}
	return g.generateValue(schemas[0], visitedRefs)
}

func (g *Generator) generateOneOf(schemas []*schema.Schema, visitedRefs map[string]bool) (interface{}, error) {
	if len(schemas) == 0 {
		return nil, nil
	}
	idx := g.rand.Intn(len(schemas))
	return g.generateValue(schemas[idx], visitedRefs)
}

func (g *Generator) generateNot(s *schema.Schema, visitedRefs map[string]bool) (interface{}, error) {
	// For 'not', generate a default value for the type
	if s.GetType() != "" {
		return g.generateDefaultForType(s.GetType()), nil
	}
	return nil, nil
}

func (g *Generator) generateIfThenElse(ifSchema, thenSchema, elseSchema *schema.Schema, visitedRefs map[string]bool) (interface{}, error) {
	// Simplified: always try then branch first
	if thenSchema != nil {
		return g.generateValue(thenSchema, visitedRefs)
	}
	if elseSchema != nil {
		return g.generateValue(elseSchema, visitedRefs)
	}
	return nil, nil
}

func (g *Generator) generateString(s *schema.Schema) string {
	if s.GetConst() != nil {
		if str, ok := s.GetConst().(string); ok {
			return str
		}
	}

	if len(s.GetEnum()) > 0 {
		return g.generateEnum(s.GetEnum()).(string)
	}

	format := s.GetFormat()
	switch format {
	case "date-time":
		return time.Now().Format(time.RFC3339)
	case "date":
		return time.Now().Format("2006-01-02")
	case "time":
		return time.Now().Format("15:04:05")
	case "email":
		return "user@example.com"
	case "hostname":
		return "example.com"
	case "ipv4":
		return "192.168.1.1"
	case "ipv6":
		return "2001:db8::1"
	case "uri", "url":
		return "https://example.com"
	case "uuid":
		return "550e8400-e29b-41d4-a716-446655440000"
	case "byte":
		return "YWJjZGVmZw=="
	case "password":
		return "secret123"
	default:
		if s.GetPattern() != "" {
			return g.generateFromPattern(s.GetPattern())
		}
		if s.Title != "" {
			return s.Title
		}
		if s.Description != "" {
			words := strings.Fields(s.Description)
			if len(words) > 0 {
				return strings.Join(words[:min(3, len(words))], " ")
			}
		}
		return "example string"
	}
}

func (g *Generator) generateFromPattern(pattern string) string {
	re := regexp.MustCompile(`\[([^\]]+)\]`)
	if re.MatchString(pattern) {
		return "a"
	}
	return "match"
}

func (g *Generator) generateNumber(s *schema.Schema) interface{} {
	if s.GetConst() != nil {
		return s.GetConst()
	}

	if len(s.GetEnum()) > 0 {
		return g.generateEnum(s.GetEnum())
	}

	isInteger := s.GetType() == "integer"
	min := s.GetMinimum()
	max := s.GetMaximum()

	if g.strictMode {
		if min != 0 {
			if isInteger {
				return int64(min)
			}
			return min
		}
		return 0
	}

	var val float64
	if min != 0 && max != 0 {
		val = min + g.rand.Float64()*(max-min)
	} else if min != 0 {
		val = min + g.rand.Float64()*100
	} else if max != 0 {
		val = max - g.rand.Float64()*100
	} else {
		val = g.rand.Float64() * 100
	}

	if isInteger {
		return int64(math.Round(val))
	}

	// Handle multipleOf
	if s.MultipleOf != nil && *s.MultipleOf > 0 {
		multiple := *s.MultipleOf
		steps := int((max - min) / multiple)
		if steps < 1 {
			steps = 1
		}
		step := g.rand.Intn(steps + 1)
		return min + float64(step)*multiple
	}

	return val
}

func (g *Generator) generateBoolean() bool {
	if g.strictMode {
		return false
	}
	return g.rand.Float64() > 0.5
}

func (g *Generator) generateArray(s *schema.Schema, visitedRefs map[string]bool) ([]interface{}, error) {
	g.currentDepth++
	defer func() { g.currentDepth-- }()

	minItems := s.GetMinItems()
	maxItems := s.GetMaxItems()
	if maxItems < 0 || maxItems > g.maxArraySize {
		maxItems = g.maxArraySize
	}
	if minItems > maxItems {
		minItems = maxItems
	}

	count := minItems
	if maxItems > minItems && !g.strictMode {
		count = minItems + g.rand.Intn(maxItems-minItems+1)
	}
	if count < 0 {
		count = 0
	}

	var itemsSchema *schema.Schema
	if s.Items != nil {
		itemsSchema = s.Items
	}

	result := make([]interface{}, 0, count)

	// Handle prefixItems (tuple validation)
	if len(s.PrefixItems) > 0 {
		for i := 0; i < count && i < len(s.PrefixItems); i++ {
			val, err := g.generateValue(s.PrefixItems[i], visitedRefs)
			if err != nil {
				return nil, err
			}
			result = append(result, val)
		}
		// Fill remaining with items schema if available
		for i := len(s.PrefixItems); i < count; i++ {
			if itemsSchema != nil {
				val, err := g.generateValue(itemsSchema, visitedRefs)
				if err != nil {
					return nil, err
				}
				result = append(result, val)
			}
		}
		return result, nil
	}

	// Regular array with items schema
	for i := 0; i < count; i++ {
		if itemsSchema != nil {
			val, err := g.generateValue(itemsSchema, visitedRefs)
			if err != nil {
				return nil, err
			}
			result = append(result, val)
		} else {
			result = append(result, g.randomValue())
		}
	}

	// Handle uniqueItems
	if s.UniqueItems != nil && *s.UniqueItems {
		result = g.uniqueValues(result)
	}

	return result, nil
}

func (g *Generator) uniqueValues(arr []interface{}) []interface{} {
	seen := make(map[string]bool)
	result := make([]interface{}, 0, len(arr))
	for _, v := range arr {
		key := fmt.Sprintf("%v", v)
		if !seen[key] {
			seen[key] = true
			result = append(result, v)
		}
	}
	return result
}

func (g *Generator) generateObject(s *schema.Schema, visitedRefs map[string]bool) (map[string]interface{}, error) {
	g.currentDepth++
	defer func() { g.currentDepth-- }()

	result := make(map[string]interface{})

	// Handle required properties first
	for _, req := range s.Required {
		if propSchema, ok := s.Properties[req]; ok {
			val, err := g.generateValue(propSchema, visitedRefs)
			if err != nil {
				return nil, err
			}
			result[req] = val
		}
	}

	// Handle optional properties
	if !g.requiredOnly && !g.strictMode {
		for propName, propSchema := range s.Properties {
			if _, exists := result[propName]; !exists {
				if g.rand.Float32() > 0.3 {
					val, err := g.generateValue(propSchema, visitedRefs)
					if err != nil {
						return nil, err
					}
					result[propName] = val
				}
			}
		}
	}

	// Handle patternProperties
	for pattern, propSchema := range s.PatternProperties {
		if !g.strictMode && g.rand.Float32() > 0.5 {
			key := g.generateKeyFromPattern(pattern)
			val, err := g.generateValue(propSchema, visitedRefs)
			if err != nil {
				return nil, err
			}
			result[key] = val
		}
	}

	// Handle additionalProperties
	if s.AdditionalProperties != nil && !g.strictMode {
		if addSchema, ok := s.AdditionalProperties.(*schema.Schema); ok {
			for i := 0; i < g.rand.Intn(3); i++ {
				key := fmt.Sprintf("extraProperty%d", i)
				val, err := g.generateValue(addSchema, visitedRefs)
				if err == nil {
					result[key] = val
				}
			}
		} else if addProps, ok := s.AdditionalProperties.(bool); ok && addProps {
			for i := 0; i < g.rand.Intn(3); i++ {
				result[fmt.Sprintf("additional%d", i)] = g.randomValue()
			}
		}
	}

	// Handle minProperties
	if s.MinProperties != nil && len(result) < *s.MinProperties && !g.strictMode {
		for propName, propSchema := range s.Properties {
			if _, exists := result[propName]; !exists && len(result) < *s.MinProperties {
				val, err := g.generateValue(propSchema, visitedRefs)
				if err == nil {
					result[propName] = val
				}
			}
		}
	}

	// Handle maxProperties
	if s.MaxProperties != nil && len(result) > *s.MaxProperties {
		keys := make([]string, 0, len(result))
		for k := range result {
			keys = append(keys, k)
		}
		for i := *s.MaxProperties; i < len(keys); i++ {
			delete(result, keys[i])
		}
	}

	return result, nil
}

func (g *Generator) generateKeyFromPattern(pattern string) string {
	patterns := map[string]string{
		"^[a-z]+$":       "key",
		"^[A-Z]+$":       "KEY",
		"^[a-zA-Z]+$":    "KeyName",
		"^\\w+$":         "key_name",
		"^[a-z_]+$":      "key_name",
		"^\\d+$":         "123",
		"^[a-z]+\\d+$":   "key1",
	}
	if val, ok := patterns[pattern]; ok {
		return val
	}
	return "key"
}

func (g *Generator) generateEnum(values []interface{}) interface{} {
	if len(values) == 0 {
		return nil
	}
	if g.strictMode {
		return values[0]
	}
	return values[g.rand.Intn(len(values))]
}

func (g *Generator) generateDefaultForType(t string) interface{} {
	switch t {
	case "string":
		return "example"
	case "number", "integer":
		return 0
	case "boolean":
		return false
	case "array":
		return []interface{}{}
	case "object":
		return map[string]interface{}{}
	case "null":
		return nil
	default:
		return nil
	}
}

func (g *Generator) randomValue() interface{} {
	switch g.rand.Intn(5) {
	case 0:
		return g.randomString(10)
	case 1:
		return g.rand.Intn(1000)
	case 2:
		return g.rand.Float64() * 1000
	case 3:
		return g.rand.Intn(2) == 1
	case 4:
		return nil
	}
	return "random"
}

func (g *Generator) randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[g.rand.Intn(len(charset))]
	}
	return string(b)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
	ref = strings.TrimPrefix(ref, "#/")
	parts := strings.Split(ref, "/")

	var current interface{} = r.definitions
	for _, part := range parts {
		part = strings.ReplaceAll(part, "~1", "/")
		part = strings.ReplaceAll(part, "~0", "~")

		switch v := current.(type) {
		case map[string]*schema.Schema:
			if schema, ok := v[part]; ok {
				current = schema
			} else {
				return nil, fmt.Errorf("reference not found: %s", ref)
			}
		case *schema.Schema:
			if part == "properties" {
				current = v.Properties
			} else if part == "items" {
				current = v.Items
			} else if strings.HasPrefix(part, "definitions/") || strings.HasPrefix(part, "$defs/") {
				name := strings.TrimPrefix(part, "definitions/")
				name = strings.TrimPrefix(name, "$defs/")
				if schema, ok := v.Definitions[name]; ok {
					current = schema
				} else if schema, ok := v.Defs[name]; ok {
					current = schema
				} else {
					return nil, fmt.Errorf("definition not found: %s", name)
				}
			} else {
				return nil, fmt.Errorf("cannot navigate into schema: %s", part)
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
