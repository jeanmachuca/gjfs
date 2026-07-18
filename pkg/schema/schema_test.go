package schema

import (
	"testing"
)

const toolManifestJSON = `{
  "tools": [
    {
      "name": "test.tool",
      "description": "A test tool",
      "inputSchema": {
        "type": "object",
        "properties": {
          "name": {"type": "string"}
        },
        "required": ["name"]
      },
      "outputSchema": {
        "type": "object",
        "properties": {
          "result": {"type": "string"}
        }
      }
    }
  ]
}`

func TestParseToolManifest(t *testing.T) {
	s, err := ParseSchemaFromString(toolManifestJSON)
	if err != nil {
		t.Fatalf("ParseSchemaFromString failed: %v", err)
	}
	if !s.IsToolManifest() {
		t.Fatal("Expected IsToolManifest to be true")
	}
	entries := s.ToolSchemas()
	if len(entries) != 2 {
		t.Fatalf("Expected 2 tool schema entries, got %d", len(entries))
	}
	if entries[0].ToolName != "test.tool" {
		t.Errorf("Expected tool name 'test.tool', got %s", entries[0].ToolName)
	}
	if entries[0].Kind != ToolInputSchema {
		t.Errorf("Expected first entry to be inputSchema, got %s", entries[0].Kind)
	}
	if entries[1].Kind != ToolOutputSchema {
		t.Errorf("Expected second entry to be outputSchema, got %s", entries[1].Kind)
	}
	if entries[0].Schema.GetType() != "object" {
		t.Errorf("Expected input schema type 'object', got %s", entries[0].Schema.GetType())
	}
}

func TestParseRegularSchema(t *testing.T) {
	s, err := ParseSchemaFromString(`{"type": "object", "properties": {"x": {"type": "string"}}}`)
	if err != nil {
		t.Fatalf("ParseSchemaFromString failed: %v", err)
	}
	if s.IsToolManifest() {
		t.Fatal("Expected IsToolManifest to be false for regular schema")
	}
	entries := s.ToolSchemas()
	if len(entries) != 0 {
		t.Fatalf("Expected 0 tool schema entries, got %d", len(entries))
	}
}

func TestEmptyToolManifest(t *testing.T) {
	s, err := ParseSchemaFromString(`{"tools": []}`)
	if err != nil {
		t.Fatalf("ParseSchemaFromString failed: %v", err)
	}
	if s.IsToolManifest() {
		t.Fatal("Expected IsToolManifest to be false for empty tools array")
	}
}

func TestToolWithoutOutput(t *testing.T) {
	s, err := ParseSchemaFromString(`{
		"tools": [{
			"name": "test",
			"description": "test",
			"inputSchema": {"type": "object"}
		}]
	}`)
	if err != nil {
		t.Fatalf("ParseSchemaFromString failed: %v", err)
	}
	entries := s.ToolSchemas()
	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry (no outputSchema), got %d", len(entries))
	}
	if entries[0].Kind != ToolInputSchema {
		t.Errorf("Expected inputSchema, got %s", entries[0].Kind)
	}
}
