package example

import (
	"errors"
	"fmt"

	"github.com/lovromazgon/orion"
)

var ErrKeyNotFound = errors.New("key not found")

type Database interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

type contract struct {
	db Database
}

func newContract(db Database) *contract {
	return &contract{db: db}
}

func (c *contract) AfterSet(key, value string) orion.Breach {
	val, err := c.db.Get(key)
	if err != nil {
		return orion.NewBreach(fmt.Errorf("unexpected error :%w", err))
	}
	if val != value {
		return orion.NewBreach(fmt.Errorf("expected %q, got %q", value, val))
	}
	return orion.NoBreach
}

func (c *contract) AfterDelete(key string) orion.Breach {
	val, err := c.db.Get(key)
	if err != ErrKeyNotFound {
		return orion.NewBreach(fmt.Errorf("expected %q, got %q", ErrKeyNotFound, err))
	}
	if val != "" {
		return orion.NewBreach(fmt.Errorf("expected empty value, got %q", val))
	}
	return orion.NoBreach
	// TODO assert returning value
}
