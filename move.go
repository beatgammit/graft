package graft

import (
	"fmt"
	"reflect"
)

type opMove struct {
	from, path pointer
}

func (op opMove) Apply(v reflect.Value) (func(), error) {
	return nil, fmt.Errorf("Not implemented")
}
