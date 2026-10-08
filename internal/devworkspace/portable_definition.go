package devworkspace

import (
	"fmt"
	"sort"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/asac"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagebundle"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagefmt"
)

const PortableDefinitionScope = "woobe-portable-definition@1.0"

// PortableDefinition seals execution intent from the closed package. Registry
// aliases, physical paths, package names/versions and native credentials do not
// determine it. Unreferenced captured dependencies remain part of the closure.
func PortableDefinition(bundle *packagebundle.Bundle) (string, map[string]string, error) {
	files := map[string]string{}
	for _, item := range bundle.Inventory {
		files[item.Path] = "sha256:" + item.SHA256
	}
	requirements := packagefmt.Object(packagefmt.Object(bundle.Graph.Manifest["spec"])["requires"])
	intents := map[string]any{}
	for _, raw := range packagefmt.List(requirements["credentials"]) {
		requirement := clone(packagefmt.Object(raw))
		alias := packagefmt.Text(requirement["ref"])
		delete(requirement, "ref")
		intents[alias] = requirement
	}
	components := map[string]string{}
	for _, key := range bundle.Graph.Order {
		document := clone(bundle.Graph.Components[key])
		delete(packagefmt.Object(document["metadata"]), "key")
		delete(document, "provenance")
		RewriteReferences(document, components)
		spec := packagefmt.Object(document["spec"])
		if document["kind"] == "Model" {
			alias := packagefmt.Text(packagefmt.Object(spec["credential"])["ref"])
			if intents[alias] == nil {
				return "", nil, fmt.Errorf("Model has no portable credential intent")
			}
			spec["credential"] = intents[alias]
		}
		hash := func(raw any) (string, error) {
			path, err := packagebundle.PortablePath(packagefmt.Text(raw), bundle.Graph.Paths[key])
			if err != nil {
				return "", err
			}
			if files[path] == "" {
				return "", fmt.Errorf("support inventory is incomplete")
			}
			return files[path], nil
		}
		if document["kind"] == "Skill" {
			pkg := packagefmt.Object(spec["package"])
			manifest, err := hash(pkg["manifest"])
			if err != nil {
				return "", nil, err
			}
			pkg["manifest"] = manifest
			resources := []any{}
			for _, raw := range packagefmt.List(pkg["resources"]) {
				digest, err := hash(raw)
				if err != nil {
					return "", nil, err
				}
				resources = append(resources, digest)
			}
			pkg["resources"] = resources
		}
		if document["kind"] == "Knowledge" {
			for _, raw := range packagefmt.List(spec["documents"]) {
				item := packagefmt.Object(raw)
				digest, err := hash(item["path"])
				if err != nil {
					return "", nil, err
				}
				item["path"] = digest
			}
		}
		digest, err := asac.Digest("portable-component", document)
		if err != nil {
			return "", nil, err
		}
		components[key] = digest
	}
	entry := packagefmt.Text(packagefmt.Object(packagefmt.Object(bundle.Graph.Manifest["spec"])["entrypoint"])["ref"])
	values := []string{}
	for _, digest := range components {
		values = append(values, digest)
	}
	sort.Strings(values)
	public := clone(requirements)
	delete(public, "credentials")
	digest, err := asac.Digest("portable-definition", map[string]any{"entrypoint": components[entry], "components": values, "requires": public})
	return digest, components, err
}

// PortableProviderDigests follows the verified public requirements, excluding
// destination credentials and providers outside the captured closure.
func PortableProviderDigests(bundle *packagebundle.Bundle) (map[string]string, error) {
	result := map[string]string{}
	requirements := packagefmt.Object(packagefmt.Object(bundle.Graph.Manifest["spec"])["requires"])
	for _, raw := range packagefmt.List(requirements["credentials"]) {
		item := packagefmt.Object(raw)
		value := clone(packagefmt.Object(item["metadata"]))
		if value == nil {
			value = map[string]any{}
		}
		value["provider"] = item["provider"]
		digest, err := asac.Digest("portable-provider", value)
		if err != nil {
			return nil, err
		}
		result[packagefmt.Text(item["ref"])] = digest
	}
	return result, nil
}
