package sqlxx

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

// Logger is an observer that collect queries executed in sqlxx.
type Logger interface {
	Log(query string, duration time.Duration)
}

// Log will emmit given queries on driver's attached Logger.
func Log(driver Driver, queries Queries, duration time.Duration) {
	if driver == nil || queries == nil {
		return
	}
	go func() {
		buffer := &bytes.Buffer{}

		for i, query := range queries {
			statement := query.Query

			if i != 0 {
				buffer.WriteString("\n")
			}

			for k, v := range query.Params {
				key := fmt.Sprint(":", k)
				value := formatLog(v)
				statement = strings.Replace(statement, key, value, -1)
			}

			for _, v := range query.Args {
				value := formatLog(v)
				statement = strings.Replace(statement, "?", value, 1)
			}

			if !strings.HasSuffix(statement, ";") {
				statement = fmt.Sprint(statement, ";")
			}

			buffer.WriteString(statement)

		}

		driver.logger().Log(buffer.String(), duration)

	}()
}

// formatLog try to convert any type into a valid log statement.
func formatLog(value interface{}) string {
	if value == nil {
		return "NULL"
	}
	switch v := value.(type) {
	case sql.NullBool:
		if v.Valid {
			return formatLogRaw(v.Bool)
		}
	case *sql.NullBool:
		if v.Valid {
			return formatLogRaw(v.Bool)
		}
	case sql.NullFloat64:
		if v.Valid {
			return formatLogRaw(v.Float64)
		}
	case *sql.NullFloat64:
		if v.Valid {
			return formatLogRaw(v.Float64)
		}
	case sql.NullInt64:
		if v.Valid {
			return formatLogRaw(v.Int64)
		}
	case *sql.NullInt64:
		if v.Valid {
			return formatLogRaw(v.Int64)
		}
	case sql.NullString:
		if v.Valid {
			return formatLogString(v.String)
		}
	case *sql.NullString:
		if v.Valid {
			return formatLogString(v.String)
		}
	case pq.NullTime:
		if v.Valid {
			return formatLogTime(v.Time)
		}
	case *pq.NullTime:
		if v.Valid {
			return formatLogTime(v.Time)
		}
	case string:
		return formatLogString(v)
	case *string:
		return formatLogString(*v)
	case time.Time:
		return formatLogTime(v)
	case *time.Time:
		return formatLogTime(*v)
	default:
		return formatLogRaw(v)
	}
	return "NULL"
}

func formatLogRaw(v interface{}) string {
	return fmt.Sprint(v)
}

func formatLogString(v string) string {
	return fmt.Sprint("'", v, "'")
}

func formatLogTime(v time.Time) string {
	return fmt.Sprint("'", v.Format(time.RFC3339), "'")
}

// EmptyLogger is a no-op Logger.
type EmptyLogger struct{}

func (EmptyLogger) Log(query string, duration time.Duration) {}
