package orion

import (
	"errors"
	"fmt"
)

var ErrKeyNotFound = errors.New("key not found")

type Database interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

type contract struct {
	contractBase
	o         *O
	db        Database
	deleteErr error
}

func (c *contract) Init(o *O, db Database) {
	c.o = o
	c.db = db
}

func (c *contract) AfterSet(key, value string) {
	val, err := c.db.Get(key)
	if err != nil {
		c.o.NewBreach(fmt.Errorf("unexpected error :%w", err))
	}
	if val != value {
		c.o.NewBreach(fmt.Errorf("expected %q, got %q", value, val))
	}
}

func (c *contract) BeforeDelete(key string) {
	_, err := c.db.Get(key)
	c.deleteErr = err
}

func (c *contract) AfterDelete(key string) {
	val, err := c.db.Get(key)
	if err != ErrKeyNotFound {
		c.o.NewBreach(fmt.Errorf("expected %q, got %q", ErrKeyNotFound, err))
	}
	if val != "" {
		c.o.NewBreach(fmt.Errorf("expected empty value, got %q", val))
	}
	// TODO assert returning value
}

// -------------------
// ---- GENERATED ----
// -------------------

type DatabaseWithContract interface {
	Database
	contract() *contract
}

func NewDatabaseWithContract(db Database, handler BreachHandler) DatabaseWithContract {
	o := New(handler)
	c := &contract{}
	c.Init(o, db)
	return &databaseWithContract{
		db: db,
		c:  c,
	}
}

type databaseWithContract struct {
	db Database
	c  *contract
}

func (d databaseWithContract) contract() *contract {
	return d.c
}

func (d databaseWithContract) Get(key string) (string, error) {
	d.c.BeforeGet(key)
	defer d.c.AfterGet(key)
	return d.db.Get(key)
}

func (d databaseWithContract) Set(key, value string) error {
	d.c.BeforeSet(key, value)
	defer d.c.AfterSet(key, value)
	return d.db.Set(key, value)
}

func (d databaseWithContract) Delete(key string) error {
	d.c.BeforeDelete(key)
	defer d.c.AfterDelete(key)
	return d.db.Delete(key)
}

type contractBase struct{}

func (contractBase) base()                       {}
func (contractBase) BeforeGet(key string)        {}
func (contractBase) AfterGet(key string)         {}
func (contractBase) BeforeSet(key, value string) {}
func (contractBase) AfterSet(key, value string)  {}
func (contractBase) BeforeDelete(key string)     {}
func (contractBase) AfterDelete(key string)      {}
