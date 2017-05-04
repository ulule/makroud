package sqlxx_test

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestSave_Save(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := User{Username: "thoas"}

	_, err := sqlxx.Save(db, &user)
	assert.NoError(t, err)

	assert.NotZero(t, user.ID)
	assert.Equal(t, true, user.IsActive)
	assert.NotZero(t, user.UpdatedAt)

	user.Username = "gilles"
	_, err = sqlxx.Save(db, &user)
	assert.NoError(t, err)

	m := map[string]interface{}{"username": "gilles"}

	query := `
	SELECT count(*)
	FROM %s
	WHERE username = :username
	`

	stmt, err := db.PrepareNamed(fmt.Sprintf(query, user.TableName()))
	assert.Nil(t, err)

	var count int
	err = stmt.Get(&count, m)
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}
