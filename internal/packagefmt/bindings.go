package packagefmt

import (
	"io"
	"os"
	"regexp"
	"strings"
)

var bindingID = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// LoadBindings resolves only public identities and protected references. It must
// never open a credential store or include material in the portable bundle.
func LoadBindings(path string, shortcuts []string, graph *Graph) (map[string]any, error) {
	document := map[string]any{"format": "woobe-package", "schema_version": "1.0", "kind": "ImportBindings",
		"metadata": map[string]any{"name": "Destination"}, "spec": map[string]any{}}
	if path != "" {
		info, err := os.Lstat(path)
		if err != nil || !info.Mode().IsRegular() || info.Size() > 2<<20 {
			return nil, &Diagnostic{Code: "PACKAGE_BINDINGS_INVALID", Message: "Bindings must be a bounded regular file"}
		}
		file, err := openBindingsRead(path)
		if err != nil {
			return nil, &Diagnostic{Code: "PACKAGE_BINDINGS_INVALID", Message: "Bindings file is unavailable"}
		}
		opened, statErr := file.Stat()
		if statErr != nil || !os.SameFile(info, opened) || !opened.Mode().IsRegular() {
			_ = file.Close()
			return nil, &Diagnostic{Code: "PACKAGE_BINDINGS_INVALID", Message: "Bindings file changed while opening"}
		}
		data, readErr := io.ReadAll(io.LimitReader(file, 2<<20+1))
		_ = file.Close()
		if readErr != nil || len(data) > 2<<20 {
			return nil, &Diagnostic{Code: "PACKAGE_LIMIT_EXCEEDED", Message: "Bindings file exceeds 2 MiB"}
		}
		parsed, err := Decode(data, path)
		if err != nil {
			return nil, err
		}
		if err = Validate(parsed); err != nil {
			return nil, err
		}
		if parsed.Value["kind"] != "ImportBindings" {
			return nil, &Diagnostic{Code: "PACKAGE_BINDINGS_INVALID", Message: "Expected ImportBindings"}
		}
		document = parsed.Value
	}
	spec := Object(document["spec"])
	for _, shortcut := range shortcuts {
		name, value, ok := strings.Cut(shortcut, "=")
		group, alias, typed := strings.Cut(name, ".")
		if !ok || !typed || group != "credential" || alias == "" || !bindingID.MatchString(value) {
			return nil, &Diagnostic{Code: "PACKAGE_BINDINGS_INVALID", Message: "--bind requires credential.ALIAS=UUID"}
		}
		credentials := Object(spec["credentials"])
		if credentials == nil {
			credentials = map[string]any{}
			spec["credentials"] = credentials
		}
		if _, exists := credentials[alias]; exists {
			return nil, &Diagnostic{Code: "PACKAGE_BINDINGS_DUPLICATE", Message: "Binding alias is declared more than once"}
		}
		credentials[alias] = map[string]any{"credential_id": value}
	}
	parsed := &Document{Value: document, File: path}
	if err := Validate(parsed); err != nil {
		return nil, err
	}
	requires := Object(Object(graph.Manifest["spec"])["requires"])
	for _, group := range []string{"credentials", "project_environment", "secrets", "knowledge"} {
		aliases := map[string]bool{}
		items, _ := requires[group].([]any)
		for _, item := range items {
			field := "ref"
			if group == "project_environment" {
				field = "key"
			}
			alias, _ := Object(item)[field].(string)
			if aliases[alias] {
				return nil, &Diagnostic{Code: "PACKAGE_REQUIREMENT_DUPLICATE", Message: "Destination requirement is duplicated"}
			}
			aliases[alias] = true
		}
		supplied := Object(spec[group])
		if len(supplied) != len(aliases) {
			return nil, &Diagnostic{Code: "PACKAGE_BINDINGS_INCOMPLETE", Message: "Bindings must match declared destination requirements"}
		}
		for alias := range supplied {
			if !aliases[alias] {
				return nil, &Diagnostic{Code: "PACKAGE_BINDINGS_INCOMPLETE", Message: "Binding alias is not declared by the package"}
			}
		}
	}
	metadata := Object(document["metadata"])
	if _, exists := metadata["description"]; !exists {
		metadata["description"] = ""
	}
	for _, group := range []string{"credentials", "project_environment", "secrets", "knowledge"} {
		values := Object(spec[group])
		if values == nil {
			values = map[string]any{}
			spec[group] = values
		}
		for _, raw := range values {
			for key, value := range Object(raw) {
				if value == nil {
					delete(Object(raw), key)
				}
			}
		}
	}
	return document, nil
}
