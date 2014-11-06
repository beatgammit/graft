package graft

import (
	"encoding/json"
	"testing"
)

func TestRemove(t *testing.T) {
	testsJson := `
    [
        {
            "comment": "A.3. Removing an Object Member",
            "doc": {
                "baz": "qux",
                "foo": "bar"
            },
            "patch": [
            { "op": "remove", "path": "/baz" }
            ],
            "expected": {
                "foo": "bar"
            }
        },
        {
            "comment": "A.4. Removing an Array Element",
            "doc": {
                "foo": [ "bar", "qux", "baz" ]
            },
            "patch": [
            { "op": "remove", "path": "/foo/1" }
            ],
            "expected": {
                "foo": [ "bar", "baz" ]
            }
        },
        {
            "comment": "A.13 Invalid JSON Patch Document",
            "doc": {
                "foo": "bar"
            },
            "patch": [
            { "op": "add", "path": "/baz", "value": "qux", "op": "remove" }
            ],
            "error": "operation has two 'op' members",
            "disabled": true
        }
    ]
    `

	var tests []jsonPatchTest
	if err := json.Unmarshal([]byte(testsJson), &tests); err != nil {
		t.Fatal(err)
	}
	doTest(t, tests)
}
