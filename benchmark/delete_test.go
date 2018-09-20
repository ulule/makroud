package benchmark

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-gorp/gorp"
	"github.com/go-xorm/xorm"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"

	"github.com/ulule/makroud"
	"github.com/ulule/makroud/benchmark/mimic"
)

func BenchmarkMakroud_Delete(b *testing.B) {
	row := JetMakroud{
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

	exec := jetExec()
	exec.NumInput = -1
	dsn := mimic.NewResult(exec)

	driver, err := makroud.NewDebugClient("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()

	b.Run("makroud", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err = makroud.Delete(ctx, driver, &row)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkSQLX_Delete(b *testing.B) {
	exec := jetExec()
	exec.NumInput = -1
	dsn := mimic.NewResult(exec)

	dbx, err := sqlx.Connect("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	query := "DELETE FROM jets WHERE id = :id"
	args := map[string]interface{}{
		"id": 1,
	}

	b.Run("sqlx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := dbx.NamedExecContext(ctx, query, args)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkGORM_Delete(b *testing.B) {
	row := JetGorm{
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

	exec := jetExec()
	exec.NumInput = -1
	dsn := mimic.NewResult(exec)

	gormdb, err := gorm.Open("mimic", dsn)
	if err != nil {
		panic(err)
	}

	b.Run("gorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := gormdb.Delete(&row).Error
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkGORP_Delete(b *testing.B) {
	row := JetGorp{
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

	exec := jetExec()
	exec.NumInput = -1
	dsn := mimic.NewResult(exec)

	db, err := sql.Open("mimic", dsn)
	if err != nil {
		panic(err)
	}

	gorpdb := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	if err != nil {
		panic(err)
	}
	gorpdb.AddTable(JetGorp{}).SetKeys(true, "ID")

	b.Run("gorp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := gorpdb.Delete(&row)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkXORM_Delete(b *testing.B) {
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

	exec := jetExec()
	exec.NumInput = -1
	dsn := mimic.NewResult(exec)

	xormdb, err := xorm.NewEngine("mimic", dsn)
	if err != nil {
		panic(err)
	}

	b.Run("xorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := xormdb.Delete(&row)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
