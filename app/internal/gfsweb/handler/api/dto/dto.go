package dto

import (
	"errors"
	"reflect"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user"
)

//ErrUnknownTypeToTransform is a sentinel error indicating that a type requested
//to be transformed by this package is not known by this package.
var ErrUnknownTypeToTransform = errors.New("dto: unknown type to transform")

//typeMap is used to store a mapping from types to transform to their transformation
//functions.
var typeMap map[reflect.Type]TransformerFunc

//TransformerFunc is a type that can convert from type to another.
//Ideally it is used to convert domain types to output, marshalable types.
type TransformerFunc func(interface{}) interface{}

func init() {
	typeMap = map[reflect.Type]TransformerFunc{
		reflect.TypeOf([]*user.User{}): Users,
		reflect.TypeOf(&user.User{}):   User,

		reflect.TypeOf([]*client.Client{}): Clients,
		reflect.TypeOf(&client.Client{}):   Client,
	}
}

//Transform attempts to transform v via a known TransformerFunc and return the
//result.
//
//ErrUnknownTypeToTransform is returned if v is an unknown type.
func Transform(v interface{}) (interface{}, error) {
	transformerFunc, ok := typeMap[reflect.TypeOf(v)]
	if !ok {
		return nil, ErrUnknownTypeToTransform
	}

	return transformerFunc(v), nil
}
