package graft

import (
	"fmt"
	"reflect"
)

type opCopy struct {
	from, path pointer
}

func (op opCopy) Apply(v reflect.Value) (func(), error) {
	return nil, fmt.Errorf("Not implemented")
}
