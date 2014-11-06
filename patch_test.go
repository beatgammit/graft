package graft

import (
	"encoding/json"
	"reflect"
	"testing"
)

type jsonPatchTest struct {
	Comment       string
	Doc, Expected json.RawMessage
	Patch         json.RawMessage
	Error         string
	Disabled      bool
}

func doTest(t *testing.T, tests []jsonPatchTest) {
	for _, tst := range tests {
		if tst.Disabled {
			continue
		}

		var m, exp map[string]interface{}
		if err := json.Unmarshal(tst.Doc, &m); err != nil {
			t.Fatal(err)
		}
		if len(tst.Expected) > 0 {
			json.Unmarshal([]byte(tst.Expected), &exp)
		}

		patch, err := Parse([]byte(tst.Patch))
		if err != nil {
			t.Fatal(err)
		}

		if err := patch.Apply(m); err != nil {
			if tst.Error == "" {
				t.Error(err)
			}
		} else if tst.Error != "" {
			t.Error(tst.Error)
		}

		if exp != nil && !reflect.DeepEqual(m, exp) {
			t.Errorf("%s: '%v' != '%v'", tst.Comment, m, exp)
		}
	}
}
