package benchmark

import (
	"database/sql/driver"
	"os"
	"testing"

	"github.com/go-xorm/core"
	"github.com/ulule/makroud/benchmark/mimic"
)

func jetExecSelect() mimic.QueryResult {
	return mimic.QueryResult{
		Query: &mimic.Query{
			Cols: []string{"id", "pilot_id", "airport_id", "name", "color", "uuid", "identifier", "cargo", "manifest"},
			Vals: [][]driver.Value{
				{
					int64(1), int64(1), int64(1), "test", nil, "test", "test", []byte("test"), []byte("test"),
				},
				{
					int64(2), int64(2), int64(2), "test", nil, "test", "test", []byte("test"), []byte("test"),
				},
				{
					int64(3), int64(3), int64(3), "test", nil, "test", "test", []byte("test"), []byte("test"),
				},
				{
					int64(4), int64(4), int64(4), "test", nil, "test", "test", []byte("test"), []byte("test"),
				},
				{
					int64(5), int64(5), int64(5), "test", nil, "test", "test", []byte("test"), []byte("test"),
				},
			},
		},
	}
}

func jetExecInsert() mimic.QueryResult {
	return mimic.QueryResult{
		Result: &mimic.Result{
			NumRows: 1,
		},
		Query: &mimic.Query{
			Cols: []string{"id"},
			Vals: [][]driver.Value{
				{
					int64(1),
				},
			},
		},
	}
}

func jetExecUpdate() mimic.QueryResult {
	return mimic.QueryResult{
		Result: &mimic.Result{
			NumRows: 1,
		},
		Query: &mimic.Query{
			Cols: []string{"id", "pilot_id", "airport_id", "name", "color", "uuid", "identifier", "cargo", "manifest"},
			Vals: [][]driver.Value{
				{
					int64(1), int64(1), int64(1), "test", nil, "test", "test", []byte("test"), []byte("test"),
				},
			},
		},
	}
}

func jetExecDelete() mimic.QueryResult {
	return mimic.QueryResult{
		Result: &mimic.Result{
			NumRows: 5,
		},
	}
}

func TestMain(m *testing.M) {
	// Register the mimic driver for Xorm
	core.RegisterDriver("mimic", &mimic.XormDriver{})
	code := m.Run()
	os.Exit(code)
}
