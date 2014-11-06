package graft

import (
	"encoding/json"
	"testing"
)

func TestTest(t *testing.T) {
	type test struct {
		Comment string
		Doc     string
		Patch   string
	}

	tests := []test{
		{
			Comment: "A.8. Testing a Value: Success",
			Doc:     `{ "baz": "qux", "foo": [ "a", 2, "c" ] }`,
			Patch:   `[ { "op": "test", "path": "/baz", "value": "qux" }, { "op": "test", "path": "/foo/1", "value": 2 } ]`,
		},
	}

	for _, tst := range tests {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(tst.Doc), &m); err != nil {
			t.Fatal(err)
		}
		patch, err := Parse([]byte(tst.Patch))
		if err != nil {
			t.Fatal(err)
		}
		err = patch.Apply(m)
		if err != nil {
			t.Error(err)
		}
	}
}
