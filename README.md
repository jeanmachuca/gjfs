# gjfs - Generate JSON From Schema

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

**gjfs** is a command-line tool and Go library that generates JSON examples from JSON Schema. It's useful for testing, documentation, and API development.

## Features

- 🚀 **Generate realistic JSON examples** from any JSON Schema (Draft 4/6/7/2019-09/2020-12)
- 🎯 **Multiple generation modes**: Random, strict (defaults/examples only), seeded for reproducibility
- ✅ **Schema validation** - Validate JSON data against schemas
- 🔧 **Flexible output** - stdout, file, or as a library
- 📦 **Zero dependencies** - Pure Go implementation
- 🎨 **Format support** - email, uuid, date-time, uri, ipv4, ipv6, hostname, byte, password

## Installation

### From Source (Recommended)

```bash
git clone https://github.com/jeanmachuca/gjfs.git
cd gjfs
make install
```

### Using Go Install

```bash
go install github.com/jeanmachuca/gjfs/cmd/gjfs@latest
```

### Using the Install Script

```bash
curl -sSL https://raw.githubusercontent.com/jeanmachuca/gjfs/main/scripts/install.sh | bash
```

### Pre-built Binaries

Download from [Releases](https://github.com/jeanmachuca/gjfs/releases) for your platform.

## Quick Start

```bash
# Generate from schema file
gjfs -schema schema.json

# Generate from schema string
gjfs -schema-string '{"type": "object", "properties": {"name": {"type": "string"}}}'

# Generate with strict mode (no random values)
gjfs -schema schema.json -strict

# Generate with fixed seed for reproducibility
gjfs -schema schema.json -seed 42

# Output to file
gjfs -schema schema.json -output example.json

# Validate JSON against schema
gjfs -schema schema.json -validate data.json
```

## Usage

```
gjfs - Generate JSON examples from JSON Schema

Usage:
  gjfs [options]

Options:
  -schema string
        Path to JSON schema file (use '-' for stdin)
  -s string
        Path to JSON schema file (shorthand)
  -schema-string string
        JSON schema as string
  -output string
        Output file (default: stdout)
  -o string
        Output file (shorthand)
  -strict
        Strict mode (use defaults/examples only, no random values)
  -seed int
        Random seed for reproducible generation (0 = random)
  -version
        Show version information
  -v
        Show version information (shorthand)
  -pretty
        Pretty print JSON output (default: true)
  -validate string
        Validate a JSON file against the schema (use '-' for stdin)
  -use-examples
        Use examples from schema if available (default: true)
  -help
        Show help
  -h
        Show help (shorthand)
```

## Examples

### Basic Object Schema

```json
{
  "type": "object",
  "properties": {
    "name": {"type": "string"},
    "age": {"type": "integer", "minimum": 0},
    "email": {"type": "string", "format": "email"},
    "active": {"type": "boolean", "default": true}
  },
  "required": ["name", "email"]
}
```

```bash
gjfs -schema user.json
```

**Output:**
```json
{
  "name": "John Doe",
  "age": 28,
  "email": "user@example.com",
  "active": true
}
```

### With Enum and Const

```json
{
  "type": "object",
  "properties": {
    "status": {
      "type": "string",
      "enum": ["pending", "active", "inactive"],
      "default": "pending"
    },
    "role": {
      "type": "string",
      "const": "user"
    }
  }
}
```

```bash
gjfs -schema user.json -seed 123
```

### Array Generation

```json
{
  "type": "array",
  "items": {
    "type": "object",
    "properties": {
      "id": {"type": "integer"},
      "name": {"type": "string"}
    },
    "required": ["id"]
  },
  "minItems": 2,
  "maxItems": 5
}
```

```bash
gjfs -schema items.json
```

### Complex Nested Schema

```bash
gjfs -schema examples/schema.json
```

### Strict Mode (Deterministic)

```bash
gjfs -schema schema.json -strict
```

In strict mode, only `default` values, `examples`, and `const` are used. No random values are generated.

### Reproducible Generation

```bash
gjfs -schema schema.json -seed 42
```

Same seed always produces the same output.

### Validate JSON

```bash
# Validate a file
gjfs -schema schema.json -validate data.json

# Validate from stdin
cat data.json | gjfs -schema schema.json -validate -
```

### Read Schema from Stdin

```bash
cat schema.json | gjfs -schema -
```

## Library Usage

```go
package main

import (
    "fmt"
    "github.com/jeanmachuca/gjfs/internal/generator"
    "github.com/jeanmachuca/gjfs/pkg/schema"
)

func main() {
    // Parse schema
    schemaStr := `{"type": "object", "properties": {"name": {"type": "string"}}}`
    sch, err := schema.ParseSchemaFromString(schemaStr)
    if err != nil {
        panic(err)
    }

    // Create generator with options
    gen := generator.NewGenerator(
        generator.WithSeed(42),
        generator.WithStrictMode(false),
    )

    // Generate JSON
    jsonData, err := gen.GenerateJSON(sch)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(jsonData))
}
```

## Supported JSON Schema Features

| Feature | Support |
|---------|---------|
| Types (string, number, integer, boolean, array, object, null) | ✅ |
| Required properties | ✅ |
| Enum / Const | ✅ |
| Default values | ✅ |
| Examples | ✅ |
| String formats (email, uuid, date-time, uri, etc.) | ✅ |
| String validation (minLength, maxLength, pattern) | ✅ |
| Number validation (min/max, exclusive, multipleOf) | ✅ |
| Array validation (min/maxItems, uniqueItems, contains) | ✅ |
| Object validation (min/maxProperties, patternProperties) | ✅ |
| AdditionalProperties | ✅ |
| Composition (allOf, anyOf, oneOf, not) | ✅ |
| Conditional (if/then/else) | ⚠️ Partial |
| $ref (local definitions) | ✅ |
| DependentSchemas / DependentRequired | ✅ |
| PropertyNames | ⚠️ Partial |

## Development

### Prerequisites

- Go 1.21+

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Running Examples

```bash
make run-example
```

### Formatting & Linting

```bash
make fmt
make vet
make lint  # requires golangci-lint
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please read our contributing guidelines before submitting PRs.

## Related Projects

- [gojsonschema](https://github.com/xeipuuv/gojsonschema) - JSON Schema validation
- [json-schema-faker](https://github.com/json-schema-faker/json-schema-faker) - JavaScript equivalent
