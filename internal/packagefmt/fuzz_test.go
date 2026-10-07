package packagefmt

import "testing"

func FuzzDecodePackage(f *testing.F) {
	for _, seed := range []string{"{}", "format: woobe-package\nschema_version: \"1.0\"\n", "a: &x [*x]", "{\"a\":1,\"a\":2}", "schema_version: 1.0\n"} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, source []byte) {
		if len(source) > MaxDescriptorBytes {
			t.Skip()
		}
		document, err := Decode(source, "fuzz.yaml")
		if err == nil {
			_ = Validate(document)
		}
	})
}
