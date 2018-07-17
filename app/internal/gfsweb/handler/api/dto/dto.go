package dto

import (
	"errors"
	"reflect"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user"
)

var ErrUnknownTypeToTransform = errors.New("dto: unknown type to transform")

var typeMap map[reflect.Type]TransformerFunc

type TransformerFunc func(interface{}) interface{}

func init() {
	typeMap = map[reflect.Type]TransformerFunc{
		reflect.TypeOf([]*user.User{}): Users,
		reflect.TypeOf(&user.User{}):   User,
	}
}

func Transform(v interface{}) (interface{}, error) {
	transformerFunc, ok := typeMap[reflect.TypeOf(v)]
	if !ok {
		return nil, ErrUnknownTypeToTransform
	}

	return transformerFunc(v), nil
}
