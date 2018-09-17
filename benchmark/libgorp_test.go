package benchmark

import (
	"database/sql"
)

type PilotGorp struct {
	ID   int
	Name string
}

type JetGorp struct {
	ID         int
	PilotID    int `db:"pilot_id"`
	AirportID  int `db:"airport_id"`
	Name       string
	Color      sql.NullString
	UUID       string
	Identifier string
	Cargo      []byte
	Manifest   []byte
}

type AirportGorp struct {
	ID   int
	Size sql.NullInt64
}

type LicenseGorp struct {
	ID      int
	PilotID int `db:"pilot_id"`
}

type HangarGorp struct {
	ID   int
	Name string
}

type LanguageGorp struct {
	ID       int
	Language string
}
