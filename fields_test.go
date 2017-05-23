package sqlxx_test

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestFields_IsForeignKey(t *testing.T) {
	env := setup(t)
	defer env.teardown()

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
		schema, err := sqlxx.GetSchema(env.driver, tt.model)
		field, err := sqlxx.NewField(env.driver, &schema, tt.model, tt.field)
		assert.Nil(t, err)
		assert.Equal(t, tt.result, field.IsForeignKey, fmt.Sprintf("index: %d", i))
	}
}

func TestFields_IsExcludedField(t *testing.T) {
	env := setup(t)
	defer env.teardown()

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
		schema, err := sqlxx.GetSchema(env.driver, tt.model)
		field, err := sqlxx.NewField(env.driver, &schema, tt.model, tt.field)
		assert.Nil(t, err)
		assert.Equal(t, tt.result, field.IsExcluded, fmt.Sprintf("index: %d", i))
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
