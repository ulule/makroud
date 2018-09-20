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
	ID         int    `xorm:"pk"`
	PilotID    int    `xorm:"not null"`
	AirportID  int    `xorm:"not null"`
	Name       string `xorm:"not null"`
	Color      sql.NullString
	UUID       string `xorm:"not null"`
	Identifier string `xorm:"not null"`
	Cargo      []byte `xorm:"not null"`
	Manifest   []byte `xorm:"not null"`
}

type AirportXorm struct {
	ID   int `xorm:"pk"`
	Size sql.NullInt64
}

type LicenseXorm struct {
	ID      int `xorm:"pk"`
	Pilot   PilotXorm
	PilotId int
}

type HangarXorm struct {
	ID   int    `xorm:"pk"`
	Name string `xorm:"not null"`
}

type LanguageXorm struct {
	ID       int    `xorm:"pk"`
	Language string `xorm:"index not null"`
}

type PilotLanguageXorm struct {
	PilotID    int `xorm:"pk"`
	LanguageId int `xorm:"pk"`
}
