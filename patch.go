package graft

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	OP_ADD     = "add"
	OP_REMOVE  = "remove"
	OP_REPLACE = "replace"
	OP_MOVE    = "move"
	OP_COPY    = "copy"
	OP_TEST    = "test"
)

var (
	PATH_NO_EXIST = errors.New("Path does not exist")
)

type pointer []string

func (p pointer) String() string {
	return strings.Join(p, "/")
}

// generates a setter function for this path
// this gets around issues of unaddressable values
func (p pointer) genSetter(v reflect.Value) (func(reflect.Value), error) {
	par, err := p.get(v)
	if err != nil {
		return nil, err
	} else if par.CanSet() {
		return par.Set, nil
	}

	var set func(reflect.Value)
	parPar, _ := p.parent().get(v)
	key := p.base()
	switch parPar.Type().Kind() {
	case reflect.Map:
		// TODO: convert types
		if parPar.Type().Key().Kind() != reflect.String {
			return nil, fmt.Errorf("Only string keys is currently supported")
		}
		set = func(val reflect.Value) {
			parPar.SetMapIndex(reflect.ValueOf(key), val)
		}
	case reflect.Array, reflect.Slice:
		i, err := strconv.Atoi(key)
		if err != nil {
			return nil, err
		}
		set = func(val reflect.Value) {
			parPar.Index(i).Set(val)
		}
	default:
		return nil, fmt.Errorf("Value isn't addressable")
	}
	return set, nil
}

func (p pointer) parent() pointer {
	if len(p) > 0 {
		return p[:len(p)-1]
	} else {
		return nil
	}
}

func (p pointer) base() string {
	if len(p) == 0 {
		return ""
	}
	return p[len(p)-1]
}

func (p pointer) get(v reflect.Value) (reflect.Value, error) {
	for _, part := range p {
		if v.Type().Kind() == reflect.Ptr || v.Type().Kind() == reflect.Interface {
			v = v.Elem()
		}
		switch v.Type().Kind() {
		case reflect.Struct:
			v = v.FieldByName(part)
		case reflect.Slice, reflect.Array:
			i, err := strconv.Atoi(part)
			if err != nil {
				return reflect.Value{}, PATH_NO_EXIST
			}
			v = v.Index(i)
		case reflect.Map:
			v = v.MapIndex(reflect.ValueOf(part))
		default:
			return reflect.Value{}, PATH_NO_EXIST
		}

		if !v.IsValid() {
			return v, PATH_NO_EXIST
		}
	}
	if v.Type().Kind() == reflect.Ptr || v.Type().Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v, nil
}

func newPointer(p string) (pointer, error) {
	ptr := pointer(strings.Split(p, "/")[1:])
	if len(ptr) == 0 && len(p) > 0 {
		return nil, fmt.Errorf("Invalid pointer")
	}
	return ptr, nil
}

type Operation interface {
	Apply(reflect.Value) (func(), error)
}

type Patch []Operation

func Parse(b []byte) (Patch, error) {
	var m []interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	var p Patch
	for i, oper := range m {
		if op, ok := oper.(map[string]interface{}); !ok {
			return nil, fmt.Errorf("Not a valid operation: %d", i)
		} else if typ, ok := op["op"]; !ok {
			return nil, fmt.Errorf("Missing field 'op' in operation")
		} else {
			switch typ {
			case OP_TEST:
				path, ok := op["path"].(string)
				if !ok {
					return nil, fmt.Errorf("Missing path in '%s': %d", typ, i)
				}
				ptr, err := newPointer(path)
				if err != nil {
					return nil, err
				}
				value, ok := op["value"]
				if !ok {
					return nil, fmt.Errorf("Missig value in '%s': %d", typ, i)
				}
				p = append(p, opTest{path: ptr, value: value})
			case OP_REMOVE:
				path, ok := op["path"].(string)
				if !ok {
					return nil, fmt.Errorf("Missing path in '%s': %d", typ, i)
				}
				ptr, err := newPointer(path)
				if err != nil {
					return nil, err
				}
				p = append(p, opRemove{path: ptr})
			case OP_ADD:
				path, ok := op["path"].(string)
				if !ok {
					return nil, fmt.Errorf("Missing path in '%s': %d", typ, i)
				}
				ptr, err := newPointer(path)
				if err != nil {
					return nil, err
				}
				value, ok := op["value"]
				if !ok {
					return nil, fmt.Errorf("Missig value in '%s': %d", typ, i)
				}
				p = append(p, opAdd{path: ptr, value: value})
			case OP_REPLACE:
				path, ok := op["path"].(string)
				if !ok {
					return nil, fmt.Errorf("Missing path in '%s': %d", typ, i)
				}
				ptr, err := newPointer(path)
				if err != nil {
					return nil, err
				}
				value, ok := op["value"]
				if !ok {
					return nil, fmt.Errorf("Missig value in '%s': %d", typ, i)
				}
				p = append(p, opReplace{path: ptr, value: value})
			default:
				return nil, fmt.Errorf("Operation '%s' not implemented: %d", typ, i)
			}
		}
	}

	return p, nil
}

// Apply applies this patch to the given object.
func (p Patch) Apply(dst interface{}) error {
	v := reflect.ValueOf(dst)
	var fns []func()

	// generate functions that will apply this patch
	for _, op := range p {
		if fn, err := op.Apply(v); err != nil {
			return err
		} else if fn != nil {
			// fn may be nil if there is nothing to do
			fns = append(fns, fn)
		}
	}

	// apply changes
	for _, fn := range fns {
		fn()
	}
	return nil
}
