package graft

import (
	"encoding/json"
	"testing"
)

func TestReplace(t *testing.T) {
	testsJson := `
    [
        {
            "comment": "A.5. Replacing a Value",
            "doc": {
                "baz": "qux",
                "foo": "bar"
            },
            "patch": [
            { "op": "replace", "path": "/baz", "value": "boo" }
            ],
            "expected": {
                "baz": "boo",
                "foo": "bar"
            }
        }
    ]`

	var tests []jsonPatchTest
	if err := json.Unmarshal([]byte(testsJson), &tests); err != nil {
		t.Fatal(err)
	}
	doTest(t, tests)
}
