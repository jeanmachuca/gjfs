package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jeanmachuca/gjfs/internal/generator"
	"github.com/jeanmachuca/gjfs/pkg/schema"
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

func main() {
	var (
		schemaFile   string
		schemaString string
		outputFile   string
		strictMode   bool
		seed         int64
		showVersion  bool
		pretty       bool
		validateFile string
		useExamples  bool
		help         bool
	)

	flag.StringVar(&schemaFile, "schema", "", "Path to JSON schema file")
	flag.StringVar(&schemaFile, "s", "", "Path to JSON schema file (shorthand)")
	flag.StringVar(&schemaString, "schema-string", "", "JSON schema as string")
	flag.StringVar(&outputFile, "output", "", "Output file (default: stdout)")
	flag.StringVar(&outputFile, "o", "", "Output file (shorthand)")
	flag.BoolVar(&strictMode, "strict", false, "Strict mode (use defaults/examples only, no random values)")
	flag.Int64Var(&seed, "seed", 0, "Random seed for reproducible generation (0 = random)")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information (shorthand)")
	flag.BoolVar(&pretty, "pretty", true, "Pretty print JSON output")
	flag.StringVar(&validateFile, "validate", "", "Validate a JSON file against the schema")
	flag.BoolVar(&useExamples, "use-examples", true, "Use examples from schema if available")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&help, "h", false, "Show help (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `gjfs - Generate JSON examples from JSON Schema

Usage:
  gjfs [options]

Options:
`)
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), `
Examples:
  # Generate from schema file
  gjfs -schema schema.json

  # Generate from schema string
  gjfs -schema-string '{"type": "object", "properties": {"name": {"type": "string"}}}'

  # Generate with strict mode (no random values)
  gjfs -schema schema.json -strict

  # Generate with specific seed for reproducibility
  gjfs -schema schema.json -seed 12345

  # Output to file
  gjfs -schema schema.json -output example.json

  # Validate JSON against schema
  gjfs -schema schema.json -validate data.json

  # Read schema from stdin
  cat schema.json | gjfs -schema -

  # Read JSON to validate from stdin
  gjfs -schema schema.json -validate -
`)
	}

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("gjfs version %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Built: %s\n", buildDate)
		os.Exit(0)
	}

	// Read schema
	var sch *schema.Schema
	var err error

	if schemaString != "" {
		sch, err = schema.ParseSchemaFromString(schemaString)
	} else if schemaFile != "" {
		if schemaFile == "-" {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading schema from stdin: %v\n", err)
				os.Exit(1)
			}
			sch, err = schema.ParseSchema(data)
		} else {
			data, err := os.ReadFile(schemaFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading schema file: %v\n", err)
				os.Exit(1)
			}
			sch, err = schema.ParseSchema(data)
		}
	} else {
		// Try reading from stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading schema from stdin: %v\n", err)
				os.Exit(1)
			}
			sch, err = schema.ParseSchema(data)
		} else {
			fmt.Fprintf(os.Stderr, "Error: schema file or schema string is required\n")
			flag.Usage()
			os.Exit(1)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing schema: %v\n", err)
		os.Exit(1)
	}

	// Handle validation
	if validateFile != "" {
		validateJSON(sch, validateFile)
		return
	}

	// Handle tool manifest format
	if sch.IsToolManifest() {
		generateToolManifest(sch, seed, strictMode, outputFile)
		return
	}

	// Generate example for regular schema
	gen := generator.NewGenerator(generatorOpts(seed, strictMode, sch)...)

	var output []byte
	if pretty {
		output, err = gen.GenerateJSON(sch)
	} else {
		example, err := gen.Generate(sch)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating example: %v\n", err)
			os.Exit(1)
		}
		output, err = json.Marshal(example)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating example: %v\n", err)
		os.Exit(1)
	}

	// Write output
	if outputFile != "" && outputFile != "-" {
		err = os.WriteFile(outputFile, output, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Generated example written to %s\n", outputFile)
	} else {
		fmt.Println(string(output))
	}
}

func generatorOpts(seed int64, strictMode bool, sch *schema.Schema) []generator.GeneratorOption {
	opts := []generator.GeneratorOption{
		generator.WithStrictMode(strictMode),
	}
	if seed != 0 {
		opts = append(opts, generator.WithSeed(seed))
	}
	if len(sch.Definitions) > 0 || len(sch.Defs) > 0 {
		defs := make(map[string]*schema.Schema)
		for k, v := range sch.Definitions {
			defs[k] = v
		}
		for k, v := range sch.Defs {
			defs[k] = v
		}
		opts = append(opts, generator.WithDefinitions(defs))
	}
	return opts
}

func generateToolManifest(sch *schema.Schema, seed int64, strictMode bool, outputFile string) {
	entries := sch.ToolSchemas()
	if len(entries) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no tool schemas found in tool manifest\n")
		os.Exit(1)
	}

	result := make(map[string]interface{})
	for _, entry := range entries {
		gen := generator.NewGenerator(generatorOpts(seed, strictMode, entry.Schema)...)
		val, err := gen.Generate(entry.Schema)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to generate for %s %s: %v\n", entry.ToolName, entry.Kind, err)
			continue
		}

		key := fmt.Sprintf("%s/%s", entry.ToolName, entry.Kind)
		result[key] = map[string]interface{}{
			"description": entry.ToolDescription,
			"value":       val,
		}
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling tool manifest output: %v\n", err)
		os.Exit(1)
	}

	if outputFile != "" && outputFile != "-" {
		err = os.WriteFile(outputFile, output, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Generated tool manifest examples written to %s\n", outputFile)
	} else {
		fmt.Println(string(output))
	}
}

func validateJSON(sch *schema.Schema, validateFile string) {
	var data []byte
	var err error

	if validateFile == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(validateFile)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading JSON file: %v\n", err)
		os.Exit(1)
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	if err := sch.Validate(jsonData); err != nil {
		fmt.Fprintf(os.Stderr, "Validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Validation passed!")
}
