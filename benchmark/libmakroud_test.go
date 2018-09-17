package benchmark

import (
	"database/sql"
)

type PilotMakroud struct {
	ID        int64  `makroud:"column:id,pk"`
	Name      string `makroud:"column:name"`
	Languages []PilotLanguageMakroud
}

func (PilotMakroud) TableName() string {
	return "pilots"
}

type JetMakroud struct {
	ID         int64          `makroud:"column:id,pk"`
	PilotID    int64          `makroud:"column:pilot_id,fk:pilots"`
	AirportID  int64          `makroud:"column:airport_id,fk:airports"`
	Name       string         `makroud:"column:name"`
	Color      sql.NullString `makroud:"column:color"`
	UUID       string         `makroud:"column:uuid"`
	Identifier string         `makroud:"column:identifier"`
	Cargo      []byte         `makroud:"column:cargo"`
	Manifest   []byte         `makroud:"column:manifest"`
	Pilot      *PilotMakroud
	Airport    *AirportMakroud
}

func (JetMakroud) TableName() string {
	return "jets"
}

type AirportMakroud struct {
	ID   int64         `makroud:"column:id,pk"`
	Size sql.NullInt64 `makroud:"column:size"`
}

func (AirportMakroud) TableName() string {
	return "airports"
}

type LicenseMakroud struct {
	ID      int64 `makroud:"column:id,pk"`
	PilotID int64 `makroud:"column:pilot_id,fk:pilots"`
	Pilot   *PilotMakroud
}

func (LicenseMakroud) TableName() string {
	return "licenses"
}

type HangarMakroud struct {
	ID   int64  `makroud:"column:id,pk"`
	Name string `makroud:"column:name"`
}

func (HangarMakroud) TableName() string {
	return "hangars"
}

type LanguageMakroud struct {
	ID       int64  `makroud:"column:id,pk"`
	Language string `makroud:"column:language"`
}

func (LanguageMakroud) TableName() string {
	return "languages"
}

type PilotLanguageMakroud struct {
	PilotID    int64 `makroud:"column:pilot_id,fk:pilots"`
	LanguageID int64 `makroud:"column:language_id,fk:languages"`
	Pilot      *PilotMakroud
	Language   *LanguageMakroud
}

func (PilotLanguageMakroud) TableName() string {
	return "pilot_languages"
}
