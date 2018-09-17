package benchmark

import (
	"database/sql/driver"
	"os"
	"testing"

	"github.com/go-xorm/core"
	"github.com/ulule/makroud/benchmark/mimic"
)

func jetQuery() mimic.QueryResult {
	return mimic.QueryResult{
		Query: &mimic.Query{
			Cols: []string{"id", "pilot_id", "airport_id", "name", "color", "uuid", "identifier", "cargo", "manifest"},
			Vals: [][]driver.Value{
				[]driver.Value{
					int64(1), int64(1), int64(1), "test", nil, "test", "test", []byte("test"), []byte("test"),
				},
				[]driver.Value{
					int64(2), int64(2), int64(2), "test", nil, "test", "test", []byte("test"), []byte("test"),
				},
				[]driver.Value{
					int64(3), int64(3), int64(3), "test", nil, "test", "test", []byte("test"), []byte("test"),
				},
				[]driver.Value{
					int64(4), int64(4), int64(4), "test", nil, "test", "test", []byte("test"), []byte("test"),
				},
				[]driver.Value{
					int64(5), int64(5), int64(5), "test", nil, "test", "test", []byte("test"), []byte("test"),
				},
			},
		},
	}
}

func pilotQuery() mimic.QueryResult {
	return mimic.QueryResult{
		Query: &mimic.Query{
			Cols: []string{"id", "name"},
			Vals: [][]driver.Value{
				[]driver.Value{
					int64(1), "test",
				},
				[]driver.Value{
					int64(2), "test",
				},
				[]driver.Value{
					int64(3), "test",
				},
				[]driver.Value{
					int64(4), "test",
				},
				[]driver.Value{
					int64(5), "test",
				},
			},
		},
	}
}

func languageQuery() mimic.QueryResult {
	return mimic.QueryResult{
		Query: &mimic.Query{
			Cols: []string{"id", "name"},
			Vals: [][]driver.Value{
				[]driver.Value{
					int64(1), "test",
				},
				[]driver.Value{
					int64(2), "test",
				},
				[]driver.Value{
					int64(3), "test",
				},
				[]driver.Value{
					int64(4), "test",
				},
				[]driver.Value{
					int64(5), "test",
				},
				[]driver.Value{
					int64(6), "test",
				},
				[]driver.Value{
					int64(7), "test",
				},
				[]driver.Value{
					int64(8), "test",
				},
				[]driver.Value{
					int64(9), "test",
				},
				[]driver.Value{
					int64(10), "test",
				},
			},
		},
	}
}

func jetQueryUpdate() mimic.QueryResult {
	return mimic.QueryResult{
		Query: &mimic.Query{
			Cols: []string{"id", "pilot_id", "airport_id", "name", "color", "uuid", "identifier", "cargo", "manifest"},
			Vals: [][]driver.Value{
				[]driver.Value{
					int64(1), int64(1), int64(1), "test", nil, "test", "test", []byte("test"), []byte("test"),
				},
			},
		},
	}
}

func jetQueryInsert() mimic.QueryResult {
	return mimic.QueryResult{
		Query: &mimic.Query{
			Cols: []string{"id"},
			Vals: [][]driver.Value{
				[]driver.Value{
					int64(1),
				},
			},
		},
	}
}

func jetExec() mimic.QueryResult {
	return mimic.QueryResult{
		Result: &mimic.Result{
			NumRows: 5,
		},
	}
}

func jetExecUpdate() mimic.QueryResult {
	return mimic.QueryResult{
		Result: &mimic.Result{
			NumRows: 1,
		},
	}
}

func TestMain(m *testing.M) {
	// Register the mimic driver for Xorm
	core.RegisterDriver("mimic", &mimic.XormDriver{})
	code := m.Run()
	os.Exit(code)
}
