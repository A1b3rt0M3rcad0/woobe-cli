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
