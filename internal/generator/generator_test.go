package generator

import (
	"encoding/json"
	"testing"

	"github.com/jeanmachuca/gjfs/pkg/schema"
)

func TestGenerateString(t *testing.T) {
	sch := &schema.Schema{
		Type: "string",
	}
	gen := NewGenerator(WithSeed(12345))
	val, err := gen.Generate(sch)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if _, ok := val.(string); !ok {
		t.Errorf("Expected string, got %T", val)
	}
}

func TestGenerateInteger(t *testing.T) {
	sch := &schema.Schema{
		Type: "integer",
	}
	gen := NewGenerator(WithSeed(12345))
	val, err := gen.Generate(sch)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if _, ok := val.(int64); !ok {
		t.Errorf("Expected int64, got %T", val)
	}
}

func TestGenerateNumber(t *testing.T) {
	sch := &schema.Schema{
		Type: "number",
	}
	gen := NewGenerator(WithSeed(12345))
	val, err := gen.Generate(sch)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if _, ok := val.(float64); !ok {
		t.Errorf("Expected float64, got %T", val)
	}
}

func TestGenerateBoolean(t *testing.T) {
	sch := &schema.Schema{
		Type: "boolean",
	}
	gen := NewGenerator(WithSeed(12345))
	val, err := gen.Generate(sch)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if _, ok := val.(bool); !ok {
		t.Errorf("Expected bool, got %T", val)
	}
}

func TestGenerateArray(t *testing.T) {
	sch := &schema.Schema{
		Type: "array",
		Items: &schema.Schema{
			Type: "string",
		},
		MinItems: &[]int{2}[0],
		MaxItems: &[]int{5}[0],
	}
	gen := NewGenerator(WithSeed(12345))
	val, err := gen.Generate(sch)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	arr, ok := val.([]interface{})
	if !ok {
		t.Errorf("Expected []interface{}, got %T", val)
	}
	if len(arr) < 2 || len(arr) > 5 {
		t.Errorf("Array length %d not in range [2,5]", len(arr))
	}
	for _, item := range arr {
		if _, ok := item.(string); !ok {
			t.Errorf("Array item is not string: %T", item)
		}
	}
}

func TestGenerateObject(t *testing.T) {
	sch := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"name": {Type: "string"},
			"age":  {Type: "integer"},
		},
		Required: []string{"name"},
	}
	gen := NewGenerator(WithSeed(12345))
	val, err := gen.Generate(sch)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	obj, ok := val.(map[string]interface{})
	if !ok {
		t.Errorf("Expected map[string]interface{}, got %T", val)
	}
	if _, ok := obj["name"]; !ok {
		t.Error("Required property 'name' missing")
	}
	if _, ok := obj["name"].(string); !ok {
		t.Errorf("Property 'name' is not string: %T", obj["name"])
	}
}

func TestGenerateWithEnum(t *testing.T) {
	sch := &schema.Schema{
		Type: "string",
		Enum: []interface{}{"red", "green", "blue"},
	}
	gen := NewGenerator(WithSeed(12345))
	val, err := gen.Generate(sch)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	str, ok := val.(string)
	if !ok {
		t.Errorf("Expected string, got %T", val)
	}
	valid := false
	for _, e := range sch.Enum {
		if str == e {
			valid = true
			break
		}
	}
	if !valid {
		t.Errorf("Generated value '%s' not in enum", str)
	}
}

func TestGenerateWithConst(t *testing.T) {
	sch := &schema.Schema{
		Type:  "string",
		Const: "fixed-value",
	}
	gen := NewGenerator(WithSeed(12345))
	val, err := gen.Generate(sch)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if val != "fixed-value" {
		t.Errorf("Expected 'fixed-value', got %v", val)
	}
}

func TestGenerateWithFormat(t *testing.T) {
	tests := []struct {
		format string
		check  func(string) bool
	}{
		{"email", func(s string) bool { return len(s) > 0 }},
		{"uri", func(s string) bool { return len(s) > 0 }},
		{"uuid", func(s string) bool { return len(s) > 0 }},
		{"date-time", func(s string) bool { return len(s) > 0 }},
	}
	for _, tc := range tests {
		sch := &schema.Schema{
			Type:   "string",
			Format: tc.format,
		}
		gen := NewGenerator(WithSeed(12345))
		val, err := gen.Generate(sch)
		if err != nil {
			t.Fatalf("Generate failed for format %s: %v", tc.format, err)
		}
		str, ok := val.(string)
		if !ok {
			t.Errorf("Expected string for format %s, got %T", tc.format, val)
		}
		if !tc.check(str) {
			t.Errorf("Generated value '%s' failed check for format %s", str, tc.format)
		}
	}
}

func TestGenerateWithDefaults(t *testing.T) {
	sch := &schema.Schema{
		Type:    "object",
		Default: map[string]interface{}{"default": true},
	}
	gen := NewGenerator(WithDefaults(true))
	val, err := gen.Generate(sch)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	obj, ok := val.(map[string]interface{})
	if !ok {
		t.Errorf("Expected map, got %T", val)
	}
	if obj["default"] != true {
		t.Errorf("Default not used: %v", obj)
	}
}

func TestGenerateWithExamples(t *testing.T) {
	sch := &schema.Schema{
		Type:     "string",
		Examples: []interface{}{"example-value"},
	}
	gen := NewGenerator(WithExamples(true))
	val, err := gen.Generate(sch)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if val != "example-value" {
		t.Errorf("Expected 'example-value', got %v", val)
	}
}

func TestGenerateStrictMode(t *testing.T) {
	sch := &schema.Schema{
		Type: "string",
	}
	gen := NewGenerator(WithStrictMode(true))
	val, err := gen.Generate(sch)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	// In strict mode, should return empty string or default
	if val != "" {
		t.Errorf("Strict mode should return empty string, got %v", val)
	}
}

func TestGenerateObjectWithMinMaxProperties(t *testing.T) {
	sch := &schema.Schema{
		Type:           "object",
		Properties:     map[string]*schema.Schema{"a": {Type: "string"}, "b": {Type: "string"}, "c": {Type: "string"}},
		MinProperties:  &[]int{2}[0],
		MaxProperties:  &[]int{2}[0],
		Required:       []string{},
	}
	gen := NewGenerator(WithSeed(12345))
	val, err := gen.Generate(sch)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	obj, ok := val.(map[string]interface{})
	if !ok {
		t.Errorf("Expected map, got %T", val)
	}
	if len(obj) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(obj))
	}
}

func TestGenerateAllOf(t *testing.T) {
	sch := &schema.Schema{
		AllOf: []*schema.Schema{
			{Type: "object", Properties: map[string]*schema.Schema{"a": {Type: "string"}}},
			{Type: "object", Properties: map[string]*schema.Schema{"b": {Type: "integer"}}},
		},
	}
	gen := NewGenerator(WithSeed(12345))
	val, err := gen.Generate(sch)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	obj, ok := val.(map[string]interface{})
	if !ok {
		t.Errorf("Expected map, got %T", val)
	}
	if _, ok := obj["a"]; !ok {
		t.Error("Property 'a' missing from allOf result")
	}
	if _, ok := obj["b"]; !ok {
		t.Error("Property 'b' missing from allOf result")
	}
}

func TestGenerateAnyOf(t *testing.T) {
	sch := &schema.Schema{
		AnyOf: []*schema.Schema{
			{Type: "string", Const: "option1"},
			{Type: "integer", Const: 42},
		},
	}
	gen := NewGenerator(WithSeed(12345))
	val, err := gen.Generate(sch)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	// Should match one of the const values
	valid := val == "option1" || val == float64(42) || val == int64(42)
	if !valid {
		t.Errorf("Generated value %v not in anyOf options", val)
	}
}

func TestGenerateJSONOutput(t *testing.T) {
	sch := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"name": {Type: "string"},
			"age":  {Type: "integer"},
		},
		Required: []string{"name"},
	}
	gen := NewGenerator(WithSeed(12345))
	data, err := gen.GenerateJSON(sch)
	if err != nil {
		t.Fatalf("GenerateJSON failed: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Generated JSON is invalid: %v", err)
	}
	if _, ok := result["name"]; !ok {
		t.Error("Required property 'name' missing in JSON output")
	}
}

func TestReproducibleGeneration(t *testing.T) {
	sch := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"value": {Type: "number"},
		},
	}
	gen1 := NewGenerator(WithSeed(42))
	gen2 := NewGenerator(WithSeed(42))

	val1, _ := gen1.Generate(sch)
	val2, _ := gen2.Generate(sch)

	if !equalValues(val1, val2) {
		t.Errorf("Same seed should produce same output: %v != %v", val1, val2)
	}
}

func equalValues(a, b interface{}) bool {
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	return string(aj) == string(bj)
}

func TestComplexSchema(t *testing.T) {
	sch := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"user": {
				Type: "object",
				Properties: map[string]*schema.Schema{
					"id":    {Type: "integer", Minimum: float64Ptr(1)},
					"name":  {Type: "string", MinLength: intPtr(1)},
					"email": {Type: "string", Format: "email"},
					"roles": {
						Type:  "array",
						Items: &schema.Schema{Type: "string"},
					},
				},
				Required: []string{"id", "name"},
			},
			"active": {Type: "boolean"},
		},
		Required: []string{"user"},
	}
	gen := NewGenerator(WithSeed(12345))
	val, err := gen.Generate(sch)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	obj, ok := val.(map[string]interface{})
	if !ok {
		t.Errorf("Expected map, got %T", val)
	}
	user, ok := obj["user"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected user object, got %T", obj["user"])
	}
	if _, ok := user["id"]; !ok {
		t.Error("Missing required 'id' in user")
	}
	if _, ok := user["name"]; !ok {
		t.Error("Missing required 'name' in user")
	}
	if roles, ok := user["roles"].([]interface{}); ok {
		for _, role := range roles {
			if _, ok := role.(string); !ok {
				t.Errorf("Role is not string: %T", role)
			}
		}
	}
}

func float64Ptr(f float64) *float64 { return &f }
func intPtr(i int) *int { return &i }
