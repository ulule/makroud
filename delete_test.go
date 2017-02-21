package sqlxx_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ulule/sqlxx"
)

func TestDelete(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := User{Username: "thoas"}
	assert.Nil(t, sqlxx.Save(db, &user))
	assert.Nil(t, sqlxx.Delete(db, &user))

	m := map[string]interface{}{"username": "thoas"}

	query := `
	SELECT count(*)
	FROM %s
	WHERE username = :username
	`

	stmt, err := db.PrepareNamed(fmt.Sprintf(query, user.TableName()))
	assert.NoError(t, err)

	var count int
	err = stmt.Get(&count, m)
	assert.Nil(t, err)
	assert.Equal(t, 0, count)
}

func TestSoftDelete(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := User{Username: "thoas"}
	assert.Nil(t, sqlxx.Save(db, &user))
	assert.Nil(t, sqlxx.SoftDelete(db, &user, "DeletedAt"))

	m := map[string]interface{}{"username": "thoas"}

	query := `
	SELECT count(*)
	FROM %s
	WHERE username = :username
	AND deleted_at IS NULL
	`

	stmt, err := db.PrepareNamed(fmt.Sprintf(query, user.TableName()))
	assert.Nil(t, err)

	var count int
	err = stmt.Get(&count, m)
	assert.Nil(t, err)
	assert.Equal(t, 0, count)
}
