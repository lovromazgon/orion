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

type DatabaseWithContract interface {
	Database
	contract() *contract
}

type contract struct {
	Database
	o *orion.O
}

func NewDatabaseWithContract(db Database, handler orion.BreachHandler) DatabaseWithContract {
	return &contract{
		Database: db,
		o:        orion.New(handler),
	}
}

func (c *contract) contract() *contract {
	return c
}

func (c *contract) Set(key, value string) error {
	err := c.Database.Set(key, value)
	if err != nil {
		return err
	}

	// Set successful, check contract
	val, err := c.Database.Get(key)
	if err != nil {
		c.o.AddBreach(orion.NewBreach(fmt.Errorf("unexpected error :%w", err)))
	}
	if val != value {
		c.o.AddBreach(orion.NewBreach(fmt.Errorf("expected %q, got %q", value, val)))
	}

	return nil
}

func (c *contract) Delete(key string) error {
	err := c.Database.Delete(key)
	if err != nil {
		return err
	}

	// Delete successful, check contract
	val, err := c.Database.Get(key)
	if err != ErrKeyNotFound {
		c.o.AddBreach(orion.NewBreach(fmt.Errorf("expected %q, got %q", ErrKeyNotFound, err)))
	}
	if val != "" {
		c.o.AddBreach(orion.NewBreach(fmt.Errorf("expected empty value, got %q", val)))
	}

	return nil
}
