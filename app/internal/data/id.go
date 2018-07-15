package data

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrWrongVersion = errors.New("data: wrong version for Id")
)

//Id is the type that we are going to use as an identifier throughout the application.
//
//For storing, for instance in a database, that code and the database will need
//to know that an Id is a Version 4 UUID.
//However, through the NewId() function and its methods, client code
//can think of the Id type as some opaque identifier.
type Id uuid.UUID

//NewId returns a new Id or an error is one could not be created.
func NewId() (Id, error) {
	return fromUUID(uuid.NewRandom())
}

//ParseId attempts to parse s and return the Id parsed from it.
func ParseId(s string) (Id, error) {
	return parseResult(uuid.Parse(s))
}

//ParseIdBytes attempts to parse b and return the Id parsed from it.
func ParseIdBytes(b []byte) (Id, error) {
	return parseResult(uuid.ParseBytes(b))
}

func parseResult(uuid uuid.UUID, err error) (Id, error) {
	if err == nil {
		if uuid.Version() == 4 {
			err = ErrWrongVersion
		}
	}
	return fromUUID(uuid, err)
}

//fromUUID is a helper function to convert the library's type to our type.
func fromUUID(uuid uuid.UUID, err error) (Id, error) {
	return Id(uuid), err
}

//MustId panics if err if not nil.
//It returns id otherwise.
//
//This can be used as a utility, for instance:
// id := MustId(NewId())
func MustId(id Id, err error) Id {
	if err != nil {
		panic(err)
	}
	return id
}

//EmptyId returns an Id without any data in it.
//This can be used instead of nil to indicate a zero value.
func EmptyId() Id {
	return Id(uuid.UUID([16]byte{}))
}

//IsEmpty returns whether or not id is an empty Id.
func (id Id) IsEmpty() bool {
	return id.Equal(EmptyId())
}

//Equal returns whether or not id and other represent that same value.
func (id Id) Equal(other Id) bool {
	return bytes.Equal([]byte(id[:]), []byte(other[:]))
}

//String is the fmt.Stringer implementation.
func (id Id) String() string {
	return uuid.UUID(id).String()
}

//MarshalJSON is the json.Marshaler implementation that allows transparent marshaling
//throughout the application.
func (id Id) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

//UnmarshalJSON is the json.Unmarshaler implementation that allows for transparent
//unmarshaling throughtout the application.
func (id *Id) UnmarshalJSON(p []byte) error {
	s := new(string)
	err := json.Unmarshal(p, s)
	if err != nil {
		return err
	}
	*id, err = ParseId(*s)
	return err
}

//Scan is the sql.Scanner implementation that allows retrieving an Id transparently
//from a sql database.
func (id *Id) Scan(src interface{}) error {
	toScan := uuid.UUID([16]byte{})
	err := toScan.Scan(src)
	if err != nil {
		return err
	}
	*id = Id(uuid.UUID(toScan))
	return nil
}

//Value is the driver.Value implementation that allows storing an Id transparently
//in a sql database.
func (id Id) Value() (driver.Value, error) {
	return uuid.UUID(id).Value()
}
