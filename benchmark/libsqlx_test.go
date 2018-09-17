package benchmark

import (
	"database/sql"
)

type PilotSQLX struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type JetSQLX struct {
	ID         int64          `db:"id"`
	PilotID    int64          `db:"pilot_id"`
	AirportID  int64          `db:"airport_id"`
	Name       string         `db:"name"`
	Color      sql.NullString `db:"color"`
	UUID       string         `db:"uuid"`
	Identifier string         `db:"identifier"`
	Cargo      []byte         `db:"cargo"`
	Manifest   []byte         `db:"manifest"`
}

type AirportSQLX struct {
	ID   int64         `db:"id"`
	Size sql.NullInt64 `db:"size"`
}

type LicenseSQLX struct {
	ID      int64 `db:"id,pk"`
	PilotID int64 `db:"pilot_id"`
}

type HangarSQLX struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type LanguageSQLX struct {
	ID       int64  `db:"id"`
	Language string `db:"language"`
}

type PilotLanguageSQLX struct {
	PilotID    int64 `db:"pilot_id"`
	LanguageID int64 `db:"language_id"`
}
