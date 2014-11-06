package graft

import (
	"fmt"
	"reflect"
	"strconv"
)

type opAdd struct {
	path  pointer
	value interface{}
}

func (op opAdd) Apply(v reflect.Value) (func(), error) {
	par, err := op.path.parent().get(v)
	if err != nil {
		return nil, err
	}

	val := reflect.ValueOf(op.value)
	if !val.Type().AssignableTo(par.Type().Elem()) {
		// TODO: make a good-faith effort to convert val to the appropriate type
		return nil, fmt.Errorf("Invalid type for add: '%s' != '%s'", val.Type(), par.Type().Elem())
	}

	var fn func()
	switch par.Type().Kind() {
	case reflect.Slice:
		set, err := op.path.parent().genSetter(v)
		if err != nil {
			return nil, err
		}

		if op.path.base() == "-" {
			fn = func() {
				set(reflect.Append(par, val))
			}
		} else if i, err := strconv.Atoi(op.path.base()); err != nil {
			return nil, err
		} else {
			fn = func() {
				newSlice := reflect.New(par.Type()).Elem()
				newSlice.Set(reflect.AppendSlice(newSlice, par.Slice(0, i)))
				newSlice.Set(reflect.Append(newSlice, val))
				newSlice.Set(reflect.AppendSlice(newSlice, par.Slice(i, par.Len())))
				set(newSlice)
			}
		}
	case reflect.Map:
		fn = func() {
			par.SetMapIndex(reflect.ValueOf(op.path.base()), val)
		}
	default:
		return nil, fmt.Errorf("Cannot add to '%s'", par.Type().Kind())
	}
	return fn, nil
}
