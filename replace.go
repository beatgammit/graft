package graft

import (
	"fmt"
	"reflect"
)

type opReplace struct {
	path  pointer
	value interface{}
}

func (op opReplace) Apply(v reflect.Value) (func(), error) {
	val, err := op.path.get(v)
	if err != nil {
		return nil, err
	}

	set, err := op.path.genSetter(v)
	if err != nil {
		return nil, err
	}
	newVal := reflect.ValueOf(op.value)
	if !newVal.Type().AssignableTo(val.Type()) {
		return nil, fmt.Errorf("%s not assignable to %s", newVal.Type(), val.Type())
	}

	fn := func() {
		set(newVal)
	}
	return fn, nil
}
