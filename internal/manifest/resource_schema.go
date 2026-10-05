package manifest

func ResourceSchema() map[string]any {
	kinds := []string{}
	for _, k := range Kinds() {
		kinds = append(kinds, k.Name)
	}
	return map[string]any{"$schema": "https://json-schema.org/draft/2020-12/schema", "$id": "urn:woobe:manifest:resources:2", "type": "object", "additionalProperties": false, "required": []string{"schema_version", "project_id", "resources"}, "properties": map[string]any{"schema_version": map[string]any{"const": "2"}, "project_id": map[string]any{"type": "string", "minLength": 1}, "workspace_id": map[string]any{"type": "string"}, "resources": map[string]any{"type": "array", "minItems": 1, "maxItems": 1000, "items": map[string]any{"type": "object", "additionalProperties": false, "required": []string{"key", "kind", "action", "spec"}, "properties": map[string]any{"key": map[string]any{"type": "string", "pattern": "^[a-zA-Z0-9_-]{1,128}$"}, "kind": map[string]any{"enum": kinds}, "action": map[string]any{"enum": []string{"create", "update"}}, "resource_id": map[string]any{"type": "string", "minLength": 1}, "parents": map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "string", "minLength": 1}}, "spec": map[string]any{"type": "object"}, "depends_on": map[string]any{"type": "array", "uniqueItems": true, "items": map[string]any{"type": "string"}}, "if_match": map[string]any{"type": "string"}}}}}, "description": "Explicit resource intents compiled to steps; semantic action/parent/reference rules are additionally validated by manifest validate"}
}
