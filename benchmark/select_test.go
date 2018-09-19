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
	"github.com/ulule/loukoum"

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
			query := loukoum.Select("*").From("jets")
			err = makroud.Exec(ctx, driver, query, &store)
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
			err = dbx.SelectContext(ctx, &store, "SELECT * FROM jets")
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
			_, err = gorpdb.Select(&store, "SELECT * FROM jets")
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

func BenchmarkMakroud_SelectSubset(b *testing.B) {
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
			query := loukoum.Select("id", "name", "color", "uuid", "identifier", "cargo", "manifest").From("jets")
			err = makroud.Exec(ctx, driver, query, &store)
			if err != nil {
				b.Fatal(err)
			}
			store = nil
		}
	})
}

func BenchmarkSQLX_SelectSubset(b *testing.B) {
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
			err = dbx.SelectContext(ctx, &store, "SELECT id, name, color, uuid, identifier, cargo, manifest FROM jets")
			if err != nil {
				b.Fatal(err)
			}
			store = nil
		}
	})
}

func BenchmarkGORM_SelectSubset(b *testing.B) {
	query := jetQuery()
	dsn := mimic.NewQuery(query)

	gormdb, err := gorm.Open("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("gorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetGorm
			err := gormdb.Select("id, name, color, uuid, identifier, cargo, manifest").Find(&store).Error
			if err != nil {
				b.Fatal(err)
			}
			store = nil
		}
	})
}

func BenchmarkGORP_SelectSubset(b *testing.B) {
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
			_, err = gorpdb.Select(&store, "SELECT id, name, color, uuid, identifier, cargo, manifest FROM jets")
			if err != nil {
				b.Fatal(err)
			}
			store = nil
		}
	})
}

func BenchmarkXORM_SelectSubset(b *testing.B) {
	query := jetQuery()
	dsn := mimic.NewQuery(query)

	xormdb, err := xorm.NewEngine("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("xorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetXorm
			err = xormdb.Select("id, name, color, uuid, identifier, cargo, manifest").Find(&store)
			if err != nil {
				b.Fatal(err)
			}
			store = nil
		}
	})
}

func BenchmarkMakroud_SelectComplex(b *testing.B) {
	query := jetQuery()
	query.NumInput = -1
	dsn := mimic.NewQuery(query)

	driver, err := makroud.NewDebugClient("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()

	b.Run("makroud", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetMakroud
			query := loukoum.
				Select("id", "name", "color", "uuid", "identifier", "cargo", "manifest").
				From("jets").
				Where(loukoum.Condition("id").GreaterThan(1)).
				Where(loukoum.Condition("name").NotEqual("thing")).
				Limit(1).
				Offset(1).
				GroupBy("id")

			err = makroud.Exec(ctx, driver, query, &store)
			if err != nil {
				b.Fatal(err)
			}
			store = nil
		}
	})
}

func BenchmarkSQLX_SelectComplex(b *testing.B) {
	query := jetQuery()
	query.NumInput = -1
	dsn := mimic.NewQuery(query)

	dbx, err := sqlx.Connect("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()

	stmt := fmt.Sprint(
		"SELECT id, name, color, uuid, identifier, cargo, manifest FROM jets ",
		"WHERE id > ? AND name != ? GROUP BY id OFFSET 1 LIMIT 1",
	)

	b.Run("sqlx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetSQLX
			err = dbx.SelectContext(ctx, &store, stmt, 1, "thing")
			if err != nil {
				b.Fatal(err)
			}
			store = nil
		}
	})
}

func BenchmarkGORM_SelectComplex(b *testing.B) {
	query := jetQuery()
	query.NumInput = -1
	dsn := mimic.NewQuery(query)

	gormdb, err := gorm.Open("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("gorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetGorm
			err = gormdb.Where("id > ?", 1).
				Where("name <> ?", "thing").
				Limit(1).
				Group("id").
				Offset(1).
				Select("id, name, color, uuid, identifier, cargo, manifest").
				Find(&store).Error
			if err != nil {
				b.Fatal(err)
			}
			store = nil
		}
	})
}

func BenchmarkGORP_SelectComplex(b *testing.B) {
	query := jetQuery()
	query.NumInput = -1
	dsn := mimic.NewQuery(query)

	db, err := sql.Open("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	gorpdb := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	if err != nil {
		b.Fatal(err)
	}

	stmt := fmt.Sprint(
		"SELECT id, name, color, uuid, identifier, cargo, manifest FROM jets ",
		"WHERE id > $1 AND name != $2 GROUP BY id OFFSET $3 LIMIT $4",
	)

	b.Run("gorp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetGorp
			_, err = gorpdb.Select(&store, stmt, 1, "thing", 1, 1)
			if err != nil {
				b.Fatal(err)
			}
			store = nil
		}
	})
}

func BenchmarkXORM_SelectComplex(b *testing.B) {
	query := jetQuery()
	query.NumInput = -1
	dsn := mimic.NewQuery(query)

	xormdb, err := xorm.NewEngine("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("xorm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetXorm
			err = xormdb.
				Select("id, name, color, uuid, identifier, cargo, manifest").
				Where("id > ?", 1).
				Where("name <> ?", "thing").
				Limit(1, 1).
				GroupBy("id").
				Find(&store)
			if err != nil {
				b.Fatal(err)
			}
			store = nil
		}
	})
}
