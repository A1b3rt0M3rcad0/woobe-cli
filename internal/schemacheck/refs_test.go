package schemacheck

import "testing"

func TestLocalReferencesAndComposition(t *testing.T) {
	doc := decode(`{"components":{"schemas":{"Count":{"type":"integer","minimum":1}}}}`)
	s := decode(`{"anyOf":[{"$ref":"#/components/schemas/Count"},{"type":"null"}]}`)
	for _, v := range []string{`null`, `2`} {
		if e := Check(s, decode(v), doc); e != nil {
			t.Fatal(e)
		}
	}
	if Check(s, decode(`0`), doc) == nil {
		t.Fatal("invalid")
	}
	if Check(decode(`{"$ref":"https://foreign/schema"}`), nil, doc) == nil {
		t.Fatal("external fetched")
	}
	if Check(decode(`{"anyOf":[{}, {"format":"email"}]}`), nil, doc) == nil {
		t.Fatal("unsupported hidden")
	}
}

func TestJSONPointerReferences(t *testing.T) {
	doc := decode(`{"a/b":{"~x":[{"type":"integer"}]},"space name":true}`)
	for _, ref := range []string{"#/a~1b/~0x/0", "#/space%20name"} {
		if _, e := Resolve(doc, ref); e != nil {
			t.Fatal(ref, e)
		}
	}
	if _, e := Resolve(doc, "#"); e != nil {
		t.Fatal(e)
	}
	for _, ref := range []string{"#/a~2b", "#/a~1b/~0x/00", "#/space%xxname", "#anchor"} {
		if _, e := Resolve(doc, ref); e == nil {
			t.Fatal(ref)
		}
	}
}
