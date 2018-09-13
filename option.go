package sqlxx

import (
	"github.com/pkg/errors"
)

// Option is used to define Client configuration.
type Option func(*ClientOptions) error

// Host will configure the Client to use the given server host.
func Host(host string) Option {
	return func(options *ClientOptions) error {
		options.Host = host
		return nil
	}
}

// Port will configure the Client to use the given server port.
func Port(port int) Option {
	return func(options *ClientOptions) error {
		options.Port = port
		return nil
	}
}

// User will configure the Client to use the given username.
func User(user string) Option {
	return func(options *ClientOptions) error {
		options.User = user
		return nil
	}
}

// Password will configure the Client to use the given username.
func Password(password string) Option {
	return func(options *ClientOptions) error {
		options.Password = password
		return nil
	}
}

// Database will configure the Client to use the given database name.
func Database(dbname string) Option {
	return func(options *ClientOptions) error {
		options.DBName = dbname
		return nil
	}
}

// EnableSSL will configure the Client to enable SSL mode.
func EnableSSL() Option {
	return func(options *ClientOptions) error {
		options.SSLMode = "require"
		return nil
	}
}

// DisableSSL will configure the Client to disable SSL mode.
func DisableSSL() Option {
	return func(options *ClientOptions) error {
		options.SSLMode = "disable"
		return nil
	}
}

// SSLMode will configure the Client with given SSL mode.
func SSLMode(mode string) Option {
	return func(options *ClientOptions) error {
		options.SSLMode = mode
		return nil
	}
}

// Timezone will configure the Client to use given timezone.
func Timezone(timezone string) Option {
	return func(options *ClientOptions) error {
		options.Timezone = timezone
		return nil
	}
}

// MaxOpenConnections will configure the Client to use this maximum number of open connections to the database.
func MaxOpenConnections(maximum int) Option {
	return func(options *ClientOptions) error {
		options.MaxOpenConnections = maximum
		return nil
	}
}

// MaxIdleConnections will configure the Client to keep this maximum number of idle connections in the
// connection pool.
func MaxIdleConnections(maximum int) Option {
	return func(options *ClientOptions) error {
		options.MaxIdleConnections = maximum
		return nil
	}
}

// Cache will configure if the Client should use a cache.
func Cache(enabled bool) Option {
	return func(options *ClientOptions) error {
		options.WithCache = enabled
		return nil
	}
}

// WithLogger will attach a logger on Client.
func WithLogger(logger Logger) Option {
	return func(options *ClientOptions) error {
		if logger == nil {
			return errors.New("sqlxx: a logger instance is required")
		}
		options.Logger = logger
		return nil
	}
}

// EnableSavepoint will enable the SAVEPOINT postgresql feature.
func EnableSavepoint() Option {
	return func(options *ClientOptions) error {
		options.SavepointEnabled = true
		return nil
	}
}
