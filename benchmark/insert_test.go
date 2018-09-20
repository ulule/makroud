package benchmark

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/go-gorp/gorp"
	"github.com/go-xorm/xorm"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"

	"github.com/ulule/makroud"
	"github.com/ulule/makroud/benchmark/mimic"
)

func BenchmarkMakroud_Insert(b *testing.B) {
	row := JetMakroud{
		PilotID:    1,
		AirportID:  1,
		Name:       "test",
		Color:      sql.NullString{},
		UUID:       "test",
		Identifier: "test",
		Cargo:      []byte("test"),
		Manifest:   []byte("test"),
	}

	exec := jetExecInsert()
	exec.NumInput = -1
	dsn := mimic.NewResult(exec)

	driver, err := makroud.NewDebugClient("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()

	b.Run("makroud", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err = makroud.Save(ctx, driver, &row)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkSQLX_Insert(b *testing.B) {
	exec := jetExecInsert()
	exec.NumInput = -1
	dsn := mimic.NewResult(exec)

	dbx, err := sqlx.Connect("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	query := fmt.Sprint(
		"INSERT INTO jets SET (pilot_id, airport_id, name, color, uuid, identifier, cargo, manifest) ",
		"VALUES (:pilot_id, :airport_id, :name, :color, :uuid, :identifier, :cargo, :manifest) RETURNING id",
	)
	args := map[string]interface{}{
		"pilot_id":   1,
		"airport_id": 1,
		"name":       "test",
		"color":      sql.NullString{},
		"uuid":       "test",
		"identifier": "test",
		"cargo":      []byte("test"),
		"manifest":   []byte("test"),
	}

	b.Run("sqlx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var jet JetSQLX

			stmt, err := dbx.PrepareNamedContext(ctx, query)
			if err != nil {
				b.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.GetContext(ctx, &jet, args)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkGORM_Insert(b *testing.B) {
	row := JetGorm{
		PilotID:    1,
		AirportID:  1,
		Name:       "test",
		Color:      sql.NullString{},
		UUID:       "test",
		Identifier: "test",
		Cargo:      []byte("test"),
		Manifest:   []byte("test"),
	}

	exec := jetExecInsert()
	exec.NumInput = -1
	dsn := mimic.NewResult(exec)

	gormdb, err := gorm.Open("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("gorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := gormdb.Create(&row).Error
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkGORP_Insert(b *testing.B) {
	row := JetGorp{
		PilotID:    1,
		AirportID:  1,
		Name:       "test",
		Color:      sql.NullString{},
		UUID:       "test",
		Identifier: "test",
		Cargo:      []byte("test"),
		Manifest:   []byte("test"),
	}

	exec := jetExecInsert()
	exec.NumInput = -1
	dsn := mimic.NewResult(exec)

	db, err := sql.Open("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	gorpdb := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	if err != nil {
		panic(err)
	}
	gorpdb.AddTable(JetGorp{}).SetKeys(true, "ID")

	b.Run("gorp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := gorpdb.Insert(&row)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkXORM_Insert(b *testing.B) {
	row := JetXorm{
		ID:         1,
		PilotID:    1,
		AirportID:  1,
		Name:       "test",
		Color:      sql.NullString{},
		UUID:       "test",
		Identifier: "test",
		Cargo:      []byte("test"),
		Manifest:   []byte("test"),
	}

	exec := jetExecInsert()
	exec.NumInput = -1
	dsn := mimic.NewResult(exec)

	xormdb, err := xorm.NewEngine("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("xorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := xormdb.Insert(&row)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
