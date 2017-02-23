package sqlxx_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ulule/sqlxx"
)

func TestSave_Save(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := User{Username: "thoas"}
	assert.Nil(t, sqlxx.Save(db, &user))

	assert.NotZero(t, user.ID)
	assert.Equal(t, true, user.IsActive)
	assert.NotZero(t, user.UpdatedAt)

	user.Username = "gilles"
	assert.Nil(t, sqlxx.Save(db, &user))

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
