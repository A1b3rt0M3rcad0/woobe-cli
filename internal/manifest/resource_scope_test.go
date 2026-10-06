package manifest

import "testing"

func TestKindsDeclareScope(t *testing.T) {
	for _, k := range Kinds() {
		if k.Scope != "project" {
			t.Fatalf("unexpected scope for %s: %s", k.Name, k.Scope)
		}
	}
}
