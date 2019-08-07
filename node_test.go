package makroud_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum/v3"

	"github.com/ulule/makroud"
)

func TestNode_Connect(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		ctx := context.Background()

		node, err := makroud.Connect(makroud.ClientDriver, ClientOptions().String())
		is.NoError(err)
		is.NotEmpty(node)

		is.Equal(makroud.ClientDriver, node.DriverName())

		err = node.Ping()
		is.NoError(err)

		err = node.PingContext(ctx)
		is.NoError(err)

		err = node.Close()
		is.NoError(err)
	})
}

func TestNode_Exec(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		ctx := context.Background()

		node, err := makroud.Connect(makroud.ClientDriver, ClientOptions().String())
		is.NoError(err)
		is.NotEmpty(node)

		catName := "Oreo"
		catStmt := "DELETE FROM ztp_cat WHERE name = $1"

		saveCat := func(driver makroud.Driver) {
			cat := &Cat{Name: catName}
			err := makroud.Save(ctx, driver, cat)
			is.NoError(err)
		}

		hasCat := func(driver makroud.Driver) bool {
			query := loukoum.Select(loukoum.Count("*")).From("ztp_cat").
				Where(loukoum.Condition("name").Equal(catName))
			count, err := makroud.Count(ctx, driver, query)
			is.NoError(err)
			return count >= 1
		}

		verifyCatDeletion := func(resp sql.Result, err error) {
			is.NoError(err)
			is.NotEmpty(resp)
			n, err := resp.RowsAffected()
			is.NoError(err)
			is.Equal(int64(1), n)
		}

		is.False(hasCat(driver))

		{
			// Simple without context

			saveCat(driver)
			is.True(hasCat(driver))

			resp, err := node.Exec(catStmt, catName)
			verifyCatDeletion(resp, err)

			is.False(hasCat(driver))
		}

		{
			// Transaction without context

			saveCat(driver)
			is.True(hasCat(driver))

			tx, err := node.Begin()
			is.NoError(err)
			is.NotEmpty(tx)

			resp, err := tx.Exec(catStmt, catName)
			verifyCatDeletion(resp, err)

			err = tx.Commit()
			is.NoError(err)

			is.False(hasCat(driver))
		}

		{
			// Simple with context

			saveCat(driver)
			is.True(hasCat(driver))

			resp, err := node.ExecContext(ctx, catStmt, catName)
			verifyCatDeletion(resp, err)

			is.False(hasCat(driver))
		}

		{
			// Transaction without context

			saveCat(driver)
			is.True(hasCat(driver))

			tx, err := node.Begin()
			is.NoError(err)
			is.NotEmpty(tx)

			resp, err := tx.ExecContext(ctx, catStmt, catName)
			verifyCatDeletion(resp, err)

			err = tx.Commit()
			is.NoError(err)

			is.False(hasCat(driver))
		}
	})
}

func TestNode_Query(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		ctx := context.Background()

		node, err := makroud.Connect(makroud.ClientDriver, ClientOptions().String())
		is.NoError(err)
		is.NotEmpty(node)

		catName := "Barcode"
		catStmt := "SELECT id FROM ztp_cat WHERE name = $1"

		cat := &Cat{Name: catName}
		err = makroud.Save(ctx, driver, cat)
		is.NoError(err)

		verifyCatSelection := func(rows *sql.Rows, err error) {
			is.NoError(err)
			is.NotEmpty(rows)
			count := 0

			for rows.Next() {
				count++
				id := ""
				err = rows.Scan(&id)
				is.NoError(err)
				is.Equal(cat.ID, id)
			}
			is.Equal(1, count)

			err = rows.Err()
			is.NoError(err)

			err = rows.Close()
			is.NoError(err)
		}

		{
			// Simple without context

			rows, err := node.Query(catStmt, catName)
			verifyCatSelection(rows, err)
		}

		{
			// Transaction without context

			tx, err := node.Begin()
			is.NoError(err)
			is.NotEmpty(tx)

			rows, err := tx.Query(catStmt, catName)
			verifyCatSelection(rows, err)

			err = tx.Commit()
			is.NoError(err)
		}

		{
			// Simple with context

			rows, err := node.QueryContext(ctx, catStmt, catName)
			verifyCatSelection(rows, err)
		}

		{
			// Transaction without context

			tx, err := node.Begin()
			is.NoError(err)
			is.NotEmpty(tx)

			rows, err := tx.QueryContext(ctx, catStmt, catName)
			verifyCatSelection(rows, err)

			err = tx.Commit()
			is.NoError(err)
		}
	})
}

func TestNode_QueryRow(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		ctx := context.Background()

		node, err := makroud.Connect(makroud.ClientDriver, ClientOptions().String())
		is.NoError(err)
		is.NotEmpty(node)

		catName := "Stripes"
		catStmt := "SELECT id FROM ztp_cat WHERE name = $1"

		cat := &Cat{Name: catName}
		err = makroud.Save(ctx, driver, cat)
		is.NoError(err)

		verifyCatSelection := func(row *sql.Row) {
			is.NotEmpty(row)

			id := ""
			err := row.Scan(&id)
			is.NoError(err)
			is.Equal(cat.ID, id)
		}

		{
			// Simple without context

			row := node.QueryRow(catStmt, catName)
			verifyCatSelection(row)
		}

		{
			// Transaction without context

			tx, err := node.Begin()
			is.NoError(err)
			is.NotEmpty(tx)

			row := tx.QueryRow(catStmt, catName)
			verifyCatSelection(row)

			err = tx.Commit()
			is.NoError(err)
		}

		{
			// Simple with context

			row := node.QueryRowContext(ctx, catStmt, catName)
			verifyCatSelection(row)
		}

		{
			// Transaction without context

			tx, err := node.Begin()
			is.NoError(err)
			is.NotEmpty(tx)

			row := tx.QueryRowContext(ctx, catStmt, catName)
			verifyCatSelection(row)

			err = tx.Commit()
			is.NoError(err)
		}
	})
}

func TestNode_Prepare(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		ctx := context.Background()

		node, err := makroud.Connect(makroud.ClientDriver, ClientOptions().String())
		is.NoError(err)
		is.NotEmpty(node)

		catName1 := "Silhouette"
		catName2 := "Silk"
		catName3 := "Silja"
		catName4 := "Silver"
		catName5 := "Silvia"

		catStmt := "UPDATE ztp_cat SET name = $1 WHERE id = $2"

		cat := &Cat{Name: catName1}
		err = makroud.Save(ctx, driver, cat)
		is.NoError(err)

		verifyCatStatement := func(stmt *sql.Stmt, err error) {
			is.NoError(err)
			is.NotEmpty(stmt)
		}

		execCatStatement := func(stmt *sql.Stmt, name string) {
			resp, err := stmt.Exec(name, cat.ID)
			is.NoError(err)
			is.NotEmpty(resp)

			n, err := resp.RowsAffected()
			is.NoError(err)
			is.Equal(int64(1), n)
		}

		verifyCatName := func(driver makroud.Driver, expected string) {
			name := ""
			query := loukoum.Select("name").From("ztp_cat").
				Where(loukoum.Condition("id").Equal(cat.ID))
			err := makroud.Exec(ctx, driver, query, &name)
			is.NoError(err)
			is.Equal(expected, name)
		}

		{
			// Simple without context

			stmt, err := node.Prepare(catStmt)
			verifyCatStatement(stmt, err)
			execCatStatement(stmt, catName2)
			verifyCatName(driver, catName2)
		}

		{
			// Transaction without context

			tx, err := node.Begin()
			is.NoError(err)
			is.NotEmpty(tx)

			stmt, err := tx.Prepare(catStmt)
			verifyCatStatement(stmt, err)
			execCatStatement(stmt, catName3)
			verifyCatName(driver, catName2)

			err = tx.Commit()
			is.NoError(err)
			verifyCatName(driver, catName3)
		}

		{
			// Simple with context

			stmt, err := node.PrepareContext(ctx, catStmt)
			verifyCatStatement(stmt, err)
			execCatStatement(stmt, catName4)
			verifyCatName(driver, catName4)
		}

		{
			// Transaction with context

			tx, err := node.Begin()
			is.NoError(err)
			is.NotEmpty(tx)

			stmt, err := tx.PrepareContext(ctx, catStmt)
			verifyCatStatement(stmt, err)
			execCatStatement(stmt, catName5)
			verifyCatName(driver, catName4)

			err = tx.Commit()
			is.NoError(err)
			verifyCatName(driver, catName5)
		}
	})
}

func TestNode_Transaction(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		ctx := context.Background()

		node, err := makroud.Connect(makroud.ClientDriver, ClientOptions().String())
		is.NoError(err)
		is.NotEmpty(node)

		cat := &Cat{Name: "Yoshi_0"}
		err = makroud.Save(ctx, driver, cat)
		is.NoError(err)

		setCatName := func(node makroud.Node, name string) {
			_, err := node.Exec("UPDATE ztp_cat SET name = $1 WHERE id = $2", name, cat.ID)
			is.NoError(err)
		}

		getCatName := func(node makroud.Node) string {
			name := ""
			row := node.QueryRow("SELECT name FROM ztp_cat WHERE id = $1", cat.ID)
			is.NotEmpty(row)
			err := row.Scan(&name)
			is.NoError(err)
			return name
		}

		tx1, err := node.Begin()
		is.NoError(err)
		is.NotEmpty(tx1)

		tx2, err := tx1.Begin()
		is.NoError(err)
		is.NotEmpty(tx2)

		tx3, err := tx2.Begin()
		is.NoError(err)
		is.NotEmpty(tx3)

		is.NotEmpty(tx1.DB())
		is.Empty(node.Tx())
		is.NotEmpty(tx1.DB())
		is.NotEmpty(tx1.Tx())
		is.NotEmpty(tx2.DB())
		is.NotEmpty(tx2.Tx())
		is.NotEmpty(tx3.DB())
		is.NotEmpty(tx3.Tx())

		is.Equal("Yoshi_0", getCatName(tx1))
		is.Equal("Yoshi_0", getCatName(tx2))
		is.Equal("Yoshi_0", getCatName(tx3))
		is.Equal("Yoshi_0", getCatName(node))

		setCatName(tx3, "Yoshi_1")

		is.Equal("Yoshi_1", getCatName(tx1))
		is.Equal("Yoshi_1", getCatName(tx2))
		is.Equal("Yoshi_1", getCatName(tx3))
		is.Equal("Yoshi_0", getCatName(node))

		err = tx3.Commit()
		is.NoError(err)

		is.Equal("Yoshi_1", getCatName(tx1))
		is.Equal("Yoshi_1", getCatName(tx2))
		is.Equal("Yoshi_0", getCatName(node))

		err = tx2.Commit()
		is.NoError(err)

		is.Equal("Yoshi_1", getCatName(tx1))
		is.Equal("Yoshi_0", getCatName(node))

		err = tx1.Commit()
		is.NoError(err)

		is.Equal("Yoshi_1", getCatName(node))

		err = tx1.Commit()
		is.Error(err)
		is.Equal(makroud.ErrCommitNotInTransaction, errors.Cause(err))

		err = tx2.Rollback()
		is.NoError(err)

		is.NotEmpty(tx1.DB())
		is.Empty(node.Tx())
		is.NotEmpty(tx1.DB())
		is.Empty(tx1.Tx())
		is.NotEmpty(tx2.DB())
		is.Empty(tx2.Tx())
		is.NotEmpty(tx3.DB())
		is.Empty(tx3.Tx())
	})
}

func TestNode_Stats(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)

		node, err := makroud.Connect(makroud.ClientDriver, ClientOptions().String())
		is.NoError(err)
		is.NotEmpty(node)

		stats := node.Stats()
		is.NotEmpty(stats)
	})
}
