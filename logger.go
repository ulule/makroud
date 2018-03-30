package sqlxx

import (
	"time"
)

// Logger is an observer that collect queries executed in sqlxx.
type Logger interface {
	Log(query string, duration time.Duration)
}

// Log will emmit given queries on driver's attached Logger.
func Log(driver Driver, queries Queries, duration time.Duration) {
	if driver == nil || len(queries) == 0 {
		return
	}
	go func() {
		query := queries.String()
		driver.logger().Log(query, duration)
	}()
}

// EmptyLogger is a no-op Logger.
type EmptyLogger struct{}

func (EmptyLogger) Log(query string, duration time.Duration) {}
