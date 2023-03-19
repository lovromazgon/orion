package example

import (
	"errors"
	"testing"

	"github.com/lovromazgon/orion"
	"github.com/matryer/is"
)

type InMemoryDB struct {
	values map[string]string
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{values: make(map[string]string)}
}

func (db *InMemoryDB) Get(key string) (string, error) {
	val, ok := db.values[key]
	if ok {
		return val, nil
	}
	return "", errors.New("oops")
}

func (db *InMemoryDB) Set(key, value string) error {
	db.values[key] = value
	return nil
}

func (db *InMemoryDB) Delete(key string) error {
	delete(db.values, key)
	return nil
}

func TestSimple(t *testing.T) {
	is := is.New(t)
	db := NewDatabaseWithContract(NewInMemoryDB(), orion.TestBreachHandler(t))
	err := db.Set("a", "foo")
	is.NoErr(err)
	val, err := db.Get("a")
	is.NoErr(err)
	is.Equal(val, "foo")

	err = db.Delete("a")
	is.NoErr(err)
}
