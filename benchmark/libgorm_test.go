package benchmark

import (
	"database/sql"
)

type PilotGorm struct {
	ID        int
	Name      string         `gorm:"not null"`
	Languages []LanguageGorm `gorm:"many2many:pilot_languages;"`
}

type JetGorm struct {
	ID int

	Pilot   PilotGorm `gorm:"ForeignKey:PilotID"`
	PilotID int       `gorm:"not null"`

	Airport   AirportGorm `gorm:"ForeignKey:Airport"`
	AirportID int         `gorm:"not null"`

	Name       string `gorm:"not null"`
	Color      sql.NullString
	UUID       string `gorm:"not null"`
	Identifier string `gorm:"not null"`
	Cargo      []byte `gorm:"not null"`
	Manifest   []byte `gorm:"not null"`
}

type AirportGorm struct {
	ID   int
	Size sql.NullInt64
}

type LicenseGorm struct {
	ID int

	Pilot   PilotGorm `gorm:"ForeignKey:PilotID"`
	PilotID int
}

type HangarGorm struct {
	ID   int
	Name string `gorm:"not null"`
}

type LanguageGorm struct {
	ID       int
	Language string `gorm:"index:idx_pilot_languages;not null"`
}
