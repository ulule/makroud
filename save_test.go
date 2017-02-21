package sqlxx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	is := assert.New(t)

	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := User{Username: "thoas"}
	require.NoError(t, Save(db, &user))

	is.NotZero(user.ID)
	is.Equal(true, user.IsActive)
	is.NotZero(user.UpdatedAt)

	user.Username = "gilles"
	require.NoError(t, Save(db, &user))

	m := map[string]interface{}{"username": "gilles"}

	query := `
	SELECT count(*)
	FROM %s
	WHERE username = :username
	`

	stmt, err := db.PrepareNamed(fmt.Sprintf(query, user.TableName()))
	require.NoError(t, err)

	var count int

	err = stmt.Get(&count, m)
	require.NoError(t, err)

	is.Equal(1, count)
}
