package makroud

import (
	"context"
	"time"
)

// Logger is an observer that collect queries executed in makroud.
type Logger interface {
	// Log push what query was executed and its duration.
	Log(ctx context.Context, query string, duration time.Duration)
}

// Log will emmit given query on driver's attached Logger.
// nolint: interfacer
func Log(ctx context.Context, driver Driver, query Query, duration time.Duration) {
	if driver == nil || !driver.hasLogger() {
		return
	}
	driver.logger().Log(ctx, query.String(), duration)
}
