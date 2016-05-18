package reflector

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Address struct {
	Street string `tag:"be"`
	Number int    `tag:"bi"`
}

type Person struct {
	Name string `tag:"bu"`
	Address
}

type CustomType int

func (p Person) Hi(name string) string {
	return fmt.Sprintf("Hi %s my name is %s", name, p.Name)
}

func TestListFieldsFlattened(t *testing.T) {
	p := Person{}
	obj := New(p)

	fields := obj.FieldFlattened()
	assert.Equal(t, len(fields), 3)
	assert.Equal(t, fields[0].Name(), "Name")
	assert.Equal(t, fields[1].Name(), "Street")
	assert.Equal(t, fields[2].Name(), "Number")
}

func TestListFields(t *testing.T) {
	p := Person{}
	obj := New(p)

	fields := obj.Fields()
	assert.Equal(t, len(fields), 2)
	assert.Equal(t, fields[0].Name(), "Name")
	assert.Equal(t, fields[1].Name(), "Address")
}

func TestListFieldsOnPointer(t *testing.T) {
	p := &Person{}
	obj := New(p)

	fields := obj.Fields()
	assert.Equal(t, len(fields), 2)
	assert.Equal(t, fields[0].Name(), "Name")
	assert.Equal(t, fields[1].Name(), "Address")
}

func TestListFieldsFlattenedOnPointer(t *testing.T) {
	p := &Person{}
	obj := New(p)

	fields := obj.FieldFlattened()
	assert.Equal(t, len(fields), 3)
	assert.Equal(t, fields[0].Name(), "Name")
	assert.Equal(t, fields[1].Name(), "Street")
	assert.Equal(t, fields[2].Name(), "Number")
}

func TestNoFieldsNoCustomType(t *testing.T) {
	assert.Equal(t, len(New(CustomType(1)).Fields()), 0)
	ct := CustomType(2)
	assert.Equal(t, len(New(&ct).Fields()), 0)
}

func TestFieldValidity(t *testing.T) {
	assert.False(t, New(CustomType(1)).Field("jkljkl").Valid())
	assert.False(t, New(Person{}).Field("street").Valid())
	assert.True(t, New(Person{}).Field("Street").Valid())
	assert.True(t, New(Person{}).Field("Number").Valid())
	assert.True(t, New(Person{}).Field("Name").Valid())
}

func TestSetFieldNonPointer(t *testing.T) {
	p := Person{}
	obj := New(p)
	assert.False(t, obj.IsPtr())

	err := obj.Field("Street").Set("ulica")
	assert.Error(t, err)
	assert.NotEqual(t, "ulica", p.Street)

	street, err := obj.Field("Street").Get()
	assert.Nil(t, err)

	// This actually don't work because p is a struct and reflector is working on it's own copy:
	assert.Equal(t, "", street)

}

func TestSetField(t *testing.T) {
	p := Person{}
	obj := New(&p)
	assert.True(t, obj.IsPtr())

	err := obj.Field("Street").Set("ulica")
	assert.Nil(t, err)
	assert.Equal(t, "ulica", p.Street)
}
