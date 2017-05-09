package sqlxx_test

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestDelete_Delete(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := User{Username: "thoas"}

	_, err := sqlxx.SaveWithQueries(db, &user)
	assert.NoError(t, err)

	queries, err := sqlxx.DeleteWithQueries(db, &user)
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 1)
	assert.Contains(t, queries[0].Query, "DELETE FROM users WHERE id = :id")

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

func TestDelete_SoftDelete(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := User{Username: "thoas"}

	_, err := sqlxx.SaveWithQueries(db, &user)
	assert.NoError(t, err)

	queries, err := sqlxx.SoftDeleteWithQueries(db, &user, "DeletedAt")
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 1)
	assert.Contains(t, queries[0].Query, "UPDATE users SET deleted_at = :deleted_at WHERE id = :id")

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
