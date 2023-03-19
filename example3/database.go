package example

import (
	"errors"
	"fmt"

	"github.com/lovromazgon/orion"
)

var ErrKeyNotFound = errors.New("key not found")

type Database interface {
	Get(key string) (value string, err error)
	Set(key, value string) error
	Delete(key string) error
}

type contract struct {
	db Database
}

func newContract(db Database) *contract {
	return &contract{db: db}
}

func (c *contract) Set(in setIn, out setOut) orion.Breach {
	if out.Error != nil {
		return orion.NoBreach
	}
	gotValue, err := c.db.Get(in.Key)
	if err != nil {
		return orion.NewBreach(fmt.Errorf("unexpected error :%w", err))
	}
	if gotValue != in.Value {
		return orion.NewBreach(fmt.Errorf("expected %+v, got %+v", in.Value, gotValue))
	}
	return orion.NoBreach
}

func (c *contract) Delete(in deleteIn, out deleteOut) orion.Breach {
	if out.Error != nil {
		return orion.NoBreach
	}
	val, err := c.db.Get(in.Key)
	if err != ErrKeyNotFound {
		return orion.NewBreach(fmt.Errorf("expected %+v, got %+v", ErrKeyNotFound, err))
	}
	if val != "" {
		return orion.NewBreach(fmt.Errorf("expected empty value, got %+v", val))
	}
	return orion.NoBreach
}
