package benchmark

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-gorp/gorp"
	"github.com/go-xorm/xorm"
	"github.com/jinzhu/gorm"
	"github.com/ulule/sqlx"

	"github.com/ulule/makroud"
	"github.com/ulule/makroud/benchmark/mimic"
)

func BenchmarkMakroud_SelectAll(b *testing.B) {
	query := jetQuery()
	dsn := mimic.NewQuery(query)

	driver, err := makroud.NewDebugClient("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()

	b.Run("makroud", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetMakroud
			err = makroud.RawExec(ctx, driver, "select * from jets", &store)
			if err != nil {
				b.Fatal(err)
			}
			store = nil
		}
	})
}

func BenchmarkSQLX_SelectAll(b *testing.B) {
	query := jetQuery()
	dsn := mimic.NewQuery(query)

	dbx, err := sqlx.Connect("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()

	b.Run("sqlx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetSQLX
			err = dbx.SelectContext(ctx, &store, "select * from jets")
			if err != nil {
				b.Fatal(err)
			}
			store = nil
		}
	})
}

func BenchmarkGORM_SelectAll(b *testing.B) {
	query := jetQuery()
	dsn := mimic.NewQuery(query)

	gormdb, err := gorm.Open("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("gorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetGorm
			err := gormdb.Find(&store).Error
			if err != nil {
				b.Fatal(err)
			}
			store = nil
		}
	})
}

func BenchmarkGORP_SelectAll(b *testing.B) {
	query := jetQuery()
	dsn := mimic.NewQuery(query)

	db, err := sql.Open("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	gorpdb := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	if err != nil {
		b.Fatal(err)
	}

	b.Run("gorp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetGorp
			_, err = gorpdb.Select(&store, "select * from jets")
			if err != nil {
				b.Fatal(err)
			}
			store = nil
		}
	})
}

func BenchmarkXORM_SelectAll(b *testing.B) {
	query := jetQuery()
	dsn := mimic.NewQuery(query)

	xormdb, err := xorm.NewEngine("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("xorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetXorm
			err = xormdb.Find(&store)
			if err != nil {
				b.Fatal(err)
			}
			store = nil
		}
	})
}
