package graft

import (
	"fmt"
	"reflect"
	"strconv"
)

type opRemove struct {
	path pointer
}

func (op opRemove) Apply(v reflect.Value) (func(), error) {
	par := op.path.parent()
	if par == nil {
		return nil, fmt.Errorf("Invalid path")
	}
	parent, err := par.get(v)
	if err != nil {
		return nil, err
	}

	var fn func()
	key := op.path.base()

	switch parent.Type().Kind() {
	case reflect.Map:
		fn = func() {
			// delete the key
			parent.SetMapIndex(reflect.ValueOf(key), reflect.Value{})
		}
	case reflect.Slice:
		set, err := op.path.parent().genSetter(v)
		if err != nil {
			return nil, err
		}

		i, err := strconv.Atoi(key)
		if err != nil {
			return nil, err
		} else if i >= parent.Len() {
			return nil, fmt.Errorf("Index out of range")
		}

		fn = func() {
			if i == parent.Len()-1 {
				set(parent.Slice(0, i))
			} else {
				// arr = append(arr[:i], arr[i+1:])
				set(reflect.AppendSlice(parent.Slice(0, i), parent.Slice(i+1, parent.Len())))
			}
		}
	default:
		return nil, fmt.Errorf("Cannot remove from type: %s", parent.Type().Kind())
	}
	return fn, nil
}
