package schemacheck

import "testing"

func TestSemanticNumericEquality(t *testing.T) {
	s := decode(`{"enum":[1,{"count":1}]}`)
	for _, v := range []string{`1.0`, `1e0`, `{"count":1.00}`} {
		if e := Check(s, decode(v), nil); e != nil {
			t.Fatal(v, e)
		}
	}
	for _, v := range []string{`2`, `"1"`, `{"count":2}`} {
		if Check(s, decode(v), nil) == nil {
			t.Fatal(v)
		}
	}
}

func TestUniqueNumericItems(t *testing.T) {
	s := decode(`{"uniqueItems":true}`)
	for _, v := range []string{`[1,2]`, `[1,"1"]`} {
		if e := Check(s, decode(v), nil); e != nil {
			t.Fatal(v, e)
		}
	}
	for _, v := range []string{`[1,1.0]`, `[{"n":1},{"n":1e0}]`} {
		if Check(s, decode(v), nil) == nil {
			t.Fatal(v)
		}
	}
}

func TestInvalidTypeDeclarations(t *testing.T) {
	for _, s := range []string{`{"type":[]}`, `{"type":"typo"}`, `{"type":["number",4]}`, `{"type":["null","null"]}`} {
		e := Check(decode(s), nil, nil)
		if x, ok := e.(*Error); !ok || !x.Unsupported {
			t.Fatal(s, e)
		}
	}
}

func TestMalformedAssertions(t *testing.T) {
	for _, s := range []string{`{"required":"name"}`, `{"enum":[]}`, `{"minItems":-1}`, `{"properties":[]}`, `{"uniqueItems":"yes"}`, `{"pattern":4}`, `{"$ref":4}`} {
		e := Check(decode(s), nil, nil)
		if x, ok := e.(*Error); !ok || !x.Unsupported {
			t.Fatal(s, e)
		}
	}
}
