package sqlxx_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestFields_IsForeignKey(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

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
		is.NoError(err)
		is.Equal(tt.result, field.IsForeignKey, fmt.Sprintf("index: %d", i))
	}
}

func TestFields_IsExcludedField(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

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
		is.NoError(err)
		is.Equal(tt.result, field.IsExcluded, fmt.Sprintf("index: %d", i))
	}
}

type StructValidField struct {
	ID int
}

func (StructValidField) TableName() string {
	return "structvalidfield"
}

func (StructValidField) PrimaryKeyType() sqlxx.PrimaryKeyType {
	return sqlxx.PrimaryKeyInteger
}

type StructUnexportedField struct {
	unexported int
}

func (StructUnexportedField) TableName() string {
	return "structunexportedfield"
}

func (StructUnexportedField) PrimaryKeyType() sqlxx.PrimaryKeyType {
	return sqlxx.PrimaryKeyInteger
}

type StructDBExcludedField struct {
	ID int `db:"-"`
}

func (StructDBExcludedField) TableName() string {
	return "structdbexcludedfield"
}

func (StructDBExcludedField) PrimaryKeyType() sqlxx.PrimaryKeyType {
	return sqlxx.PrimaryKeyInteger
}

type FKTest1 struct {
	excluded int
	ID       int
}

func (FKTest1) TableName() string {
	return "fktest1"
}

func (FKTest1) PrimaryKeyType() sqlxx.PrimaryKeyType {
	return sqlxx.PrimaryKeyInteger
}

type FKTest2 struct {
	MyFieldID int
}

func (FKTest2) TableName() string {
	return "fktest2"
}

func (FKTest2) PrimaryKeyType() sqlxx.PrimaryKeyType {
	return sqlxx.PrimaryKeyInteger
}

type FKTest3 struct {
	CustomField int `sqlxx:"fk"`
}

func (FKTest3) TableName() string {
	return "fktest3"
}

func (FKTest3) PrimaryKeyType() sqlxx.PrimaryKeyType {
	return sqlxx.PrimaryKeyInteger
}
