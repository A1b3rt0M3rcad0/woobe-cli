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

func TestUnsupportedInactiveProperty(t *testing.T) {
	e := Check(decode(`{"type":"object","properties":{"absent":{"format":"email"}}}`), decode(`{}`), nil)
	if x, ok := e.(*Error); !ok || !x.Unsupported {
		t.Fatal(e)
	}
}

func TestObjectPropertyBounds(t *testing.T) {
	s := decode(`{"minProperties":1,"maxProperties":2}`)
	for _, v := range []string{`{"a":1}`, `{"a":null,"b":2}`} {
		if e := Check(s, decode(v), nil); e != nil {
			t.Fatal(v, e)
		}
	}
	for _, v := range []string{`{}`, `{"a":1,"b":2,"c":3}`} {
		if Check(s, decode(v), nil) == nil {
			t.Fatal(v)
		}
	}
}

func TestDependentRequired(t *testing.T) {
	s := decode(`{"dependentRequired":{"card":["billing"]}}`)
	for _, v := range []string{`{}`, `{"card":null,"billing":"A"}`} {
		if e := Check(s, decode(v), nil); e != nil {
			t.Fatal(v, e)
		}
	}
	for _, v := range []string{`{"card":"C"}`} {
		if Check(s, decode(v), nil) == nil {
			t.Fatal(v)
		}
	}
}

func TestDependentSchemas(t *testing.T) {
	s := decode(`{"dependentSchemas":{"card":{"required":["billing"]}}}`)
	for _, v := range []string{`{}`, `{"card":1,"billing":"B"}`} {
		if e := Check(s, decode(v), nil); e != nil {
			t.Fatal(v, e)
		}
	}
	for _, v := range []string{`{"card":1}`} {
		if Check(s, decode(v), nil) == nil {
			t.Fatal(v)
		}
	}
}

func TestPropertyNames(t *testing.T) {
	s := decode(`{"propertyNames":{"pattern":"^[a-z]+$"}}`)
	for _, v := range []string{`{}`, `{"name":1}`} {
		if e := Check(s, decode(v), nil); e != nil {
			t.Fatal(v, e)
		}
	}
	for _, v := range []string{`{"NAME":1}`, `{"name1":1}`} {
		if Check(s, decode(v), nil) == nil {
			t.Fatal(v)
		}
	}
}
