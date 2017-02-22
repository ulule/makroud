package sqlxx_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ulule/sqlxx"
)

func TestIsForeignKey(t *testing.T) {
	testers := []struct {
		model  sqlxx.Model
		field  string
		result bool
	}{
		{FKTest1{}, "ID", false},
		{FKTest2{}, "MyFieldID", true},
		{FKTest3{}, "CustomField", true},
	}

	for i, tt := range testers {
		st, found := reflect.TypeOf(tt.model).FieldByName(tt.field)
		assert.True(t, found, fmt.Sprintf("index: %d", i))

		field, err := sqlxx.NewField(st, tt.model)
		assert.Nil(t, err)
		assert.Equal(t, tt.result, sqlxx.IsForeignKey(field), fmt.Sprintf("index: %d", i))
	}
}

func TestIsExcludedField(t *testing.T) {
	testers := []struct {
		model  sqlxx.Model
		field  string
		result bool
	}{
		{StructUnexportedField{}, "unexported", true},
		{StructDBExcludedField{}, "ID", true},
		{StructValidField{}, "ID", false},
	}

	for i, tt := range testers {
		st, found := reflect.TypeOf(tt.model).FieldByName(tt.field)
		assert.True(t, found, fmt.Sprintf("index: %d", i))

		field, err := sqlxx.NewField(st, tt.model)
		assert.Nil(t, err)
		assert.Equal(t, tt.result, sqlxx.IsExcludedField(st, field.Tags), fmt.Sprintf("index: %d", i))
	}
}

type StructValidField struct{ ID int }

func (s StructValidField) TableName() string { return "structvalidfield" }

type StructUnexportedField struct{ unexported int }

func (s StructUnexportedField) TableName() string { return "structunexportedfield" }

type StructDBExcludedField struct {
	ID int `db:"-"`
}

func (s StructDBExcludedField) TableName() string { return "structdbexcludedfield" }

type FKTest1 struct {
	excluded int
	ID       int
}

func (f FKTest1) TableName() string {
	return "fktest1"
}

type FKTest2 struct {
	MyFieldID int
}

func (f FKTest2) TableName() string {
	return "fktest2"
}

type FKTest3 struct {
	CustomField int `sqlxx:"fk"`
}

func (f FKTest3) TableName() string { return "fktest3" }
