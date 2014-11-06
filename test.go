package graft

import (
	"fmt"
	"reflect"
)

type opTest struct {
	path  pointer
	value interface{}
}

func (op opTest) Apply(v reflect.Value) (func(), error) {
	val, err := op.path.get(v)
	if err != nil {
		return nil, err
	}

	// TODO: implement DeepEqual to compare different types
	if !reflect.DeepEqual(val.Interface(), op.value) {
		return nil, fmt.Errorf("test(%s): '%v' != '%v'", op.path, val.Interface(), op.value)
	}
	return nil, nil
}
