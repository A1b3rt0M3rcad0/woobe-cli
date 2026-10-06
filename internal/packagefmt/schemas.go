package packagefmt

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed schemas/v1/catalog.json schemas/v1/catalog.sha256
var schemaFiles embed.FS

type offlineLoader struct{}

func (offlineLoader) Load(string) (any, error) {
	return nil, fmt.Errorf("remote schema loading is forbidden")
}

var schemaOnce sync.Once
var compiled map[string]*jsonschema.Schema
var compilationError error

func CatalogBytes() []byte { data, _ := schemaFiles.ReadFile("schemas/v1/catalog.json"); return data }
func CatalogDigest() string {
	data, _ := schemaFiles.ReadFile("schemas/v1/catalog.sha256")
	return string(bytes.TrimSpace(data))
}

func loadSchemas() {
	compiled = map[string]*jsonschema.Schema{}
	var catalog map[string]any
	decoder := json.NewDecoder(bytes.NewReader(CatalogBytes()))
	decoder.UseNumber()
	if err := decoder.Decode(&catalog); err != nil {
		compilationError = err
		return
	}
	for kind, document := range catalog {
		compiler := jsonschema.NewCompiler()
		compiler.UseLoader(offlineLoader{})
		uri := "urn:woobe-package:1.0:" + kind
		if err := compiler.AddResource(uri, document); err != nil {
			compilationError = err
			return
		}
		schema, err := compiler.Compile(uri)
		if err != nil {
			compilationError = err
			return
		}
		compiled[kind] = schema
	}
}

func Validate(document *Document) error {
	failure := func(code, message, path string) error {
		location := document.Locations[path]
		return &Diagnostic{Code: code, Message: message, File: document.File, Path: path, Line: location.Line, Column: location.Column}
	}
	if document.Value["format"] != "woobe-package" || document.Value["schema_version"] != "1.0" {
		return failure("PACKAGE_SCHEMA_UNSUPPORTED", "Expected woobe-package schema_version string 1.0", "")
	}
	kind, ok := document.Value["kind"].(string)
	if !ok {
		return failure("PACKAGE_SCHEMA_UNSUPPORTED", "Unknown component kind", "/kind")
	}
	schemaOnce.Do(loadSchemas)
	if compilationError != nil {
		return failure("PACKAGE_SCHEMA_UNSUPPORTED", "Embedded schema catalog is invalid", "")
	}
	schema, ok := compiled[kind]
	if !ok {
		return failure("PACKAGE_SCHEMA_UNSUPPORTED", "Unknown component kind", "/kind")
	}
	if err := schema.Validate(document.Value); err != nil {
		path := ""
		if validation, ok := err.(*jsonschema.ValidationError); ok {
			for len(validation.Causes) > 0 {
				validation = validation.Causes[0]
			}
			for _, part := range validation.InstanceLocation {
				path += "/" + part
			}
		}
		return failure("PACKAGE_SCHEMA_INVALID", "Document violates its component schema", path)
	}
	return nil
}
