package benchmark

import (
	"database/sql"
)

type PilotXorm struct {
	Id        int            `xorm:"pk"`
	Name      string         `xorm:"not null"`
	Languages []LanguageXorm `xorm:"extends"`
}

type JetXorm struct {
	Id int `xorm:"pk"`

	PilotId int `xorm:"not null"`

	AirportId int `xorm:"not null"`

	Name       string `xorm:"not null"`
	Color      sql.NullString
	Uuid       string `xorm:"not null"`
	Identifier string `xorm:"not null"`
	Cargo      []byte `xorm:"not null"`
	Manifest   []byte `xorm:"not null"`
}

type AirportXorm struct {
	Id   int `xorm:"pk"`
	Size sql.NullInt64
}

type LicenseXorm struct {
	Id int `xorm:"pk"`

	Pilot   PilotXorm
	PilotId int
}

type HangarXorm struct {
	Id   int    `xorm:"pk"`
	Name string `xorm:"not null"`
}

type LanguageXorm struct {
	Id       int    `xorm:"pk"`
	Language string `xorm:"index not null"`
}

type PilotLanguageXorm struct {
	PilotId    int `xorm:"pk"`
	LanguageId int `xorm:"pk"`
}
