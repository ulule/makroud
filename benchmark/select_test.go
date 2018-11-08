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
	exec := jetExecSelect()
	dsn := mimic.NewQuery(exec)

	driver, err := makroud.NewDebugClient("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	query := loukoum.Select("*").From("jets")

	b.Run("makroud", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetMakroud

			err = makroud.Exec(ctx, driver, query, &store)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkSQLX_SelectAll(b *testing.B) {
	exec := jetExecSelect()
	dsn := mimic.NewQuery(exec)

	dbx, err := sqlx.Connect("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	query := "SELECT * FROM jets"
	args := []interface{}{}

	b.Run("sqlx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetSQLX

			stmt, err := dbx.PreparexContext(ctx, query)
			if err != nil {
				b.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.SelectContext(ctx, &store, args...)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkGORM_SelectAll(b *testing.B) {
	exec := jetExecSelect()
	dsn := mimic.NewQuery(exec)

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
		}
	})
}

func BenchmarkGORP_SelectAll(b *testing.B) {
	exec := jetExecSelect()
	dsn := mimic.NewQuery(exec)

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
		}
	})
}

func BenchmarkXORM_SelectAll(b *testing.B) {
	exec := jetExecSelect()
	dsn := mimic.NewQuery(exec)

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
		}
	})
}

func BenchmarkMakroud_SelectSubset(b *testing.B) {
	exec := jetExecSelect()
	dsn := mimic.NewQuery(exec)

	driver, err := makroud.NewDebugClient("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	query := loukoum.Select("id", "name", "color", "uuid", "identifier", "cargo", "manifest").From("jets")

	b.Run("makroud", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetMakroud

			err = makroud.Exec(ctx, driver, query, &store)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkSQLX_SelectSubset(b *testing.B) {
	exec := jetExecSelect()
	dsn := mimic.NewQuery(exec)

	dbx, err := sqlx.Connect("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	query := "SELECT id, name, color, uuid, identifier, cargo, manifest FROM jets"
	args := []interface{}{}

	b.Run("sqlx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetSQLX

			stmt, err := dbx.PreparexContext(ctx, query)
			if err != nil {
				b.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.SelectContext(ctx, &store, args...)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkGORM_SelectSubset(b *testing.B) {
	exec := jetExecSelect()
	dsn := mimic.NewQuery(exec)

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
		}
	})
}

func BenchmarkGORP_SelectSubset(b *testing.B) {
	exec := jetExecSelect()
	dsn := mimic.NewQuery(exec)

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
		}
	})
}

func BenchmarkXORM_SelectSubset(b *testing.B) {
	exec := jetExecSelect()
	dsn := mimic.NewQuery(exec)

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
		}
	})
}

func BenchmarkMakroud_SelectComplex(b *testing.B) {
	exec := jetExecSelect()
	exec.NumInput = -1
	dsn := mimic.NewQuery(exec)

	driver, err := makroud.NewDebugClient("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	query := loukoum.
		Select("id", "name", "color", "uuid", "identifier", "cargo", "manifest").
		From("jets").
		Where(loukoum.Condition("id").GreaterThan(1)).
		Where(loukoum.Condition("name").NotEqual("thing")).
		Limit(1).
		Offset(1).
		GroupBy("id")

	b.Run("makroud", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetMakroud

			err = makroud.Exec(ctx, driver, query, &store)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkSQLX_SelectComplex(b *testing.B) {
	exec := jetExecSelect()
	exec.NumInput = -1
	dsn := mimic.NewQuery(exec)

	dbx, err := sqlx.Connect("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()

	query := fmt.Sprint(
		"SELECT id, name, color, uuid, identifier, cargo, manifest FROM jets ",
		"WHERE id > :id AND name != :name GROUP BY id OFFSET 1 LIMIT 1",
	)
	args := []interface{}{1, "thing"}

	b.Run("sqlx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetSQLX

			stmt, err := dbx.PreparexContext(ctx, query)
			if err != nil {
				b.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.SelectContext(ctx, &store, args...)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkGORM_SelectComplex(b *testing.B) {
	exec := jetExecSelect()
	exec.NumInput = -1
	dsn := mimic.NewQuery(exec)

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
		}
	})
}

func BenchmarkGORP_SelectComplex(b *testing.B) {
	exec := jetExecSelect()
	exec.NumInput = -1
	dsn := mimic.NewQuery(exec)

	db, err := sql.Open("mimic", dsn)
	if err != nil {
		b.Fatal(err)
	}

	gorpdb := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	if err != nil {
		b.Fatal(err)
	}

	query := fmt.Sprint(
		"SELECT id, name, color, uuid, identifier, cargo, manifest FROM jets ",
		"WHERE id > $1 AND name != $2 GROUP BY id OFFSET $3 LIMIT $4",
	)

	b.Run("gorp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var store []JetGorp

			_, err = gorpdb.Select(&store, query, 1, "thing", 1, 1)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkXORM_SelectComplex(b *testing.B) {
	exec := jetExecSelect()
	exec.NumInput = -1
	dsn := mimic.NewQuery(exec)

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
		}
	})
}
