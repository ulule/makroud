package sqlxx

import "testing"

func TestGetByParams(t *testing.T) {
	_, shutdown := dbConnection(t)
	shutdown()
}

func TestFindByParams(t *testing.T) {
	_, shutdown := dbConnection(t)
	shutdown()
}
