package graft

import (
	"encoding/json"
	"testing"
)

func TestAdd(t *testing.T) {
	testsJson := `
    [
        {
            "comment": "4.1. add with missing object",
            "doc": { "q": { "bar": 2 } },
            "patch": [ {"op": "add", "path": "/a/b", "value": 1} ],
            "error": "path /a does not exist -- missing objects are not created recursively"
        },
        {
            "comment": "A.1. Adding an Object Member",
            "doc": {
                "foo": "bar"
            },
            "patch": [
            { "op": "add", "path": "/baz", "value": "qux" }
            ],
            "expected": {
                "baz": "qux",
                "foo": "bar"
            }
        },
        {
            "comment": "A.2. Adding an Array Element",
            "doc": {
                "foo": [ "bar", "baz" ]
            },
            "patch": [
            { "op": "add", "path": "/foo/1", "value": "qux" }
            ],
            "expected": {
                "foo": [ "bar", "qux", "baz" ]
            }
        }
    ]
    `

	var tests []jsonPatchTest
	if err := json.Unmarshal([]byte(testsJson), &tests); err != nil {
		t.Fatal(err)
	}

	doTest(t, tests)
}
