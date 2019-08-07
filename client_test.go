package makroud_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum/v3"

	"github.com/ulule/makroud"
)

func TestClient_New(t *testing.T) {
	is := require.New(t)

	driver, err := makroud.New(Options()...)
	is.NoError(err)
	is.NotEmpty(driver)

	err = driver.Ping()
	is.NoError(err)

	is.Equal(makroud.ClientDriver, driver.DriverName())

	driver, err = makroud.New(Options(makroud.Host("10.0.0.1"))...)
	is.Error(err)
	is.Empty(driver)
}

func TestClient_Exec(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		human := &Human{Name: "Maria"}
		err := makroud.Save(ctx, driver, human)
		is.NoError(err)

		stmt1 := `UPDATE ztp_human SET name = $1 WHERE id = $2`
		stmt2 := `UPDATE ztp_human SAT name = $1 WHERE id = $2`

		err = driver.Exec(ctx, stmt1, "Delora", human.ID)
		is.NoError(err)

		err = driver.Exec(ctx, stmt2, "Veselko", human.ID)
		is.Error(err)

		is.NotPanics(func() {
			driver.MustExec(ctx, stmt1, "Liv", human.ID)
		})

		is.Panics(func() {
			driver.MustExec(ctx, stmt2, "Thorsten", human.ID)
		})
	})
}

func TestClient_Query(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		human1 := &Human{Name: "Scotty"}
		err := makroud.Save(ctx, driver, human1)
		is.NoError(err)

		human2 := &Human{Name: "Svend"}
		err = makroud.Save(ctx, driver, human2)
		is.NoError(err)

		stmt1 := `SELECT name FROM ztp_human WHERE id IN ($1, $2)`
		stmt2 := `SELECT name FROM ztp_human WHERE id`

		verifyQuery := func(rows makroud.Rows, errOpt ...error) {
			if len(errOpt) > 1 {
				is.NoError(errOpt[0])
			}
			is.NotEmpty(rows)

			names := []string{}
			for rows.Next() {
				name := ""
				err = rows.Scan(&name)
				is.NoError(err)
				names = append(names, name)
			}

			err = rows.Err()
			is.NoError(err)

			err = rows.Close()
			is.NoError(err)

			is.Len(names, 2)
			is.Contains(names, human1.Name)
			is.Contains(names, human2.Name)
		}

		verifyQuery(driver.Query(ctx, stmt1, human1.ID, human2.ID))

		rows, err := driver.Query(ctx, stmt2, human1.ID, human2.ID)
		is.Error(err)
		is.Empty(rows)

		is.NotPanics(func() {
			verifyQuery(driver.MustQuery(ctx, stmt1, human1.ID, human2.ID))
		})

		is.Panics(func() {
			rows := driver.MustQuery(ctx, stmt2, human1.ID, human2.ID)
			is.Empty(rows)
		})
	})
}

func TestClient_PrepareExec(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		human1 := &Human{Name: "Lucie"}
		err := makroud.Save(ctx, driver, human1)
		is.NoError(err)

		human2 := &Human{Name: "Louise"}
		err = makroud.Save(ctx, driver, human2)
		is.NoError(err)

		stmt, err := driver.Prepare(ctx, `UPDATE ztp_human SET deleted_at = NOW() WHERE id`)
		is.Error(err)
		is.Empty(stmt)

		stmt, err = driver.Prepare(ctx, `UPDATE ztp_human SET deleted_at = NOW() WHERE id = $1`)
		is.NoError(err)
		is.NotEmpty(stmt)

		err = stmt.Exec(ctx, human1.ID)
		is.NoError(err)

		err = stmt.Exec(ctx, human2.ID)
		is.NoError(err)

		err = stmt.Exec(ctx, struct{}{})
		is.Error(err)

		err = stmt.Close()
		is.NoError(err)

		query := loukoum.Select(loukoum.Count("*")).From("ztp_human").
			Where(loukoum.Condition("id").In(human1.ID, human2.ID)).
			And(loukoum.Condition("deleted_at").IsNull(true))

		count, err := makroud.Count(ctx, driver, query)
		is.NoError(err)
		is.Equal(int64(0), count)
	})
}

func TestClient_PrepareQueryRow(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		human1 := &Human{Name: "Adam"}
		err := makroud.Save(ctx, driver, human1)
		is.NoError(err)

		human2 := &Human{Name: "Adrien"}
		err = makroud.Save(ctx, driver, human2)
		is.NoError(err)

		stmt, err := driver.Prepare(ctx, `SELECT id, name, created_at FROM ztp_human WHERE id`)
		is.Error(err)
		is.Empty(stmt)

		stmt, err = driver.Prepare(ctx, `SELECT id, name, created_at FROM ztp_human WHERE id = $1`)
		is.NoError(err)
		is.NotEmpty(stmt)

		verifyColumns := func(row makroud.Row) {
			columns, err := row.Columns()
			is.NoError(err)
			is.NotEmpty(columns)
			is.Len(columns, 3)
			is.Equal(columns[0], "id")
			is.Equal(columns[1], "name")
			is.Equal(columns[2], "created_at")
		}

		verifyHuman := func(expected *Human, id string, name string, createdAt time.Time) {
			is.Equal(expected.ID, id)
			is.Equal(expected.Name, name)
			is.Equal(expected.CreatedAt.UnixNano(), createdAt.UnixNano())
		}

		{
			row, err := stmt.QueryRow(ctx, human1.ID)
			is.NoError(err)
			is.NotEmpty(row)

			verifyColumns(row)

			id := ""
			name := ""
			createdAt := time.Time{}

			err = row.Scan(&id, &name, &createdAt)
			is.NoError(err)
			verifyHuman(human1, id, name, createdAt)
		}

		{
			row, err := stmt.QueryRow(ctx, human2.ID)
			is.NoError(err)
			is.NotEmpty(row)

			verifyColumns(row)

			dest := map[string]interface{}{}
			err = row.Write(dest)
			is.NoError(err)

			id, ok := dest["id"].(string)
			is.True(ok, `expecting dest["id"] to be a string`)
			name, ok := dest["name"].(string)
			is.True(ok, `expecting dest["name"] to be a string`)
			createdAt, ok := dest["created_at"].(time.Time)
			is.True(ok, `expecting dest["created_at"] to be a time.Time`)

			verifyHuman(human2, id, name, createdAt)
		}

		{
			row, err := stmt.QueryRow(ctx, struct{}{})
			is.Error(err)
			is.Empty(row)
		}

		err = stmt.Close()
		is.NoError(err)
	})
}

func TestClient_PrepareQueryRows(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		human1 := &Human{Name: "Laura"}
		err := makroud.Save(ctx, driver, human1)
		is.NoError(err)

		human2 := &Human{Name: "Lana"}
		err = makroud.Save(ctx, driver, human2)
		is.NoError(err)

		stmt, err := driver.Prepare(ctx, `SELECT id, name, created_at, updated_at, deleted_at, cat_id FOM ztp_human`)
		is.Error(err)
		is.Empty(stmt)

		stmt, err = driver.Prepare(ctx, `SELECT id, name, created_at, updated_at, deleted_at, cat_id FROM ztp_human`)
		is.NoError(err)
		is.NotEmpty(stmt)

		verifyColumns := func(rows makroud.Rows) {
			columns, err := rows.Columns()
			is.NoError(err)
			is.NotEmpty(columns)
			is.Len(columns, 6)
			is.Equal(columns[0], "id")
			is.Equal(columns[1], "name")
			is.Equal(columns[2], "created_at")
			is.Equal(columns[3], "updated_at")
			is.Equal(columns[4], "deleted_at")
			is.Equal(columns[5], "cat_id")
		}

		verifyHuman := func(expected *Human, id string, name string,
			createdAt time.Time, updatedAt time.Time, deletedAt pq.NullTime, catID sql.NullString) {

			is.Equal(expected.ID, id)
			is.Equal(expected.Name, name)
			is.Equal(expected.CreatedAt.UnixNano(), createdAt.UnixNano())
			is.Equal(expected.UpdatedAt.UnixNano(), updatedAt.UnixNano())
			is.Equal(expected.DeletedAt.Valid, deletedAt.Valid)
			is.Equal(expected.DeletedAt.Time.UnixNano(), deletedAt.Time.UnixNano())
			is.Equal(expected.CatID.Valid, catID.Valid)
			is.Equal(expected.CatID.String, catID.String)
		}

		{
			rows, err := stmt.QueryRows(ctx)
			is.NoError(err)
			is.NotEmpty(rows)

			verifyColumns(rows)

			count := 0
			for rows.Next() {
				count++

				id := ""
				name := ""
				createdAt := time.Time{}
				updatedAt := time.Time{}
				deletedAt := pq.NullTime{}
				catID := sql.NullString{}

				err = rows.Scan(&id, &name, &createdAt, &updatedAt, &deletedAt, &catID)
				is.NoError(err)

				switch id {
				case human1.ID:
					verifyHuman(human1, id, name, createdAt, updatedAt, deletedAt, catID)
				case human2.ID:
					verifyHuman(human2, id, name, createdAt, updatedAt, deletedAt, catID)
				default:
					is.Fail("unexpected id", id)
				}
			}
			is.Equal(2, count)

			err = rows.Err()
			is.NoError(err)

			err = rows.Close()
			is.NoError(err)
		}

		{
			rows, err := stmt.QueryRows(ctx)
			is.NoError(err)
			is.NotEmpty(rows)

			verifyColumns(rows)

			count := 0
			for rows.Next() {
				count++

				dest := map[string]interface{}{}
				err = rows.Write(dest)
				is.NoError(err)

				id, ok := dest["id"].(string)
				is.True(ok, `expecting dest["id"] to be a string`)
				name, ok := dest["name"].(string)
				is.True(ok, `expecting dest["name"] to be a string`)
				createdAt, ok := dest["created_at"].(time.Time)
				is.True(ok, `expecting dest["created_at"] to be a time.Time`)
				updatedAt, ok := dest["updated_at"].(time.Time)
				is.True(ok, `expecting dest["updated_at"] to be a time.Time`)
				is.Nil(dest["deleted_at"])
				is.Nil(dest["cat_id"])

				switch id {
				case human1.ID:
					verifyHuman(human1, id, name, createdAt, updatedAt, pq.NullTime{}, sql.NullString{})
				case human2.ID:
					verifyHuman(human2, id, name, createdAt, updatedAt, pq.NullTime{}, sql.NullString{})
				default:
					is.Fail("unexpected id", id)
				}
			}
			is.Equal(2, count)

			err = rows.Err()
			is.NoError(err)

			err = rows.Close()
			is.NoError(err)
		}

		{
			rows, err := stmt.QueryRows(ctx, struct{}{})
			is.Error(err)
			is.Empty(rows)
		}
	})
}
