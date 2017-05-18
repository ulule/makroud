package sqlxx

// Option is used to define Client configuration.
type Option interface {
	apply(*clientOptions) error
}

type option func(*clientOptions) error

func (o option) apply(client *clientOptions) error {
	return o(client)
}

// Host will configure the Client to use the given server host.
func Host(host string) Option {
	return option(func(options *clientOptions) error {
		options.host = host
		return nil
	})
}

// Port will configure the Client to use the given server port.
func Port(port int) Option {
	return option(func(options *clientOptions) error {
		options.port = port
		return nil
	})
}

// User will configure the Client to use the given username.
func User(user string) Option {
	return option(func(options *clientOptions) error {
		options.user = user
		return nil
	})
}

// Password will configure the Client to use the given username.
func Password(password string) Option {
	return option(func(options *clientOptions) error {
		options.password = password
		return nil
	})
}

// Database will configure the Client to use the given database name.
func Database(dbname string) Option {
	return option(func(options *clientOptions) error {
		options.dbName = dbname
		return nil
	})
}

// EnableSSL will configure the Client to enable SSL mode.
func EnableSSL() Option {
	return option(func(options *clientOptions) error {
		// NOTE Some refactoring may be required to allow further options like CA certificate, etc...
		options.sslMode = "require"
		return nil
	})
}

// DisableSSL will configure the Client to disable SSL mode.
func DisableSSL() Option {
	return option(func(options *clientOptions) error {
		options.sslMode = "disable"
		return nil
	})
}

// Timezone will configure the Client to use given timezone.
func Timezone(timezone string) Option {
	return option(func(options *clientOptions) error {
		options.timezone = timezone
		return nil
	})
}

// MaxOpenConnections will configure the Client to use this maximum number of open connections to the database.
func MaxOpenConnections(maximum int) Option {
	return option(func(options *clientOptions) error {
		options.maxOpenConnections = maximum
		return nil
	})
}

// MaxIdleConnections will configure the Client to keep this maximum number of idle connections in the
// connection pool.
func MaxIdleConnections(maximum int) Option {
	return option(func(options *clientOptions) error {
		options.maxIdleConnections = maximum
		return nil
	})
}
