package makroud

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

// ClientOptions configure a Client instance.
type ClientOptions struct {
	Port               int
	Host               string
	User               string
	Password           string
	Database           string
	SSLMode            string
	Timezone           string
	MaxOpenConnections int
	MaxIdleConnections int
	WithCache          bool
	SavepointEnabled   bool
	Logger             Logger
	ApplicationName    string
	ConnectTimeout     int
}

func (e ClientOptions) String() string {
	uri := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s&timezone=%s",
		ClientDriver,
		url.QueryEscape(e.User),
		url.QueryEscape(e.Password),
		e.Host,
		e.Port,
		e.Database,
		e.SSLMode,
		e.Timezone,
	)
	if e.ApplicationName != "" {
		uri = fmt.Sprintf("%s&application_name=%s", uri, e.ApplicationName)
	}
	if e.ConnectTimeout > 0 {
		uri = fmt.Sprintf("%s&connect_timeout=%d", uri, e.ConnectTimeout)
	}
	return uri
}

// NewClientOptions creates a new ClientOptions instance with default options.
func NewClientOptions() *ClientOptions {
	return &ClientOptions{
		Host:               "localhost",
		Port:               5432,
		User:               "postgres",
		Password:           "",
		Database:           "postgres",
		SSLMode:            "disable",
		Timezone:           "UTC",
		MaxOpenConnections: 5,
		MaxIdleConnections: 2,
		WithCache:          true,
		SavepointEnabled:   false,
		ApplicationName:    "Makroud",
		ConnectTimeout:     10,
	}
}

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
		options.Database = dbname
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
		if maximum < 0 {
			return errors.New("makroud: the maximum number of open connections must be a positive number")
		}
		options.MaxOpenConnections = maximum
		return nil
	}
}

// MaxIdleConnections will configure the Client to keep this maximum number of idle connections in the
// connection pool.
func MaxIdleConnections(maximum int) Option {
	return func(options *ClientOptions) error {
		if maximum < 0 {
			return errors.New("makroud: the maximum number of idle connections must be a positive number")
		}
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
			return errors.New("makroud: a logger instance is required")
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

// ApplicationName will configure the Client to use given application name.
func ApplicationName(name string) Option {
	return func(options *ClientOptions) error {
		if len(name) >= 64 {
			return errors.New("makroud: application name must be less than 64 characters")
		}
		options.ApplicationName = name
		return nil
	}
}

// ConnectTimeout will configure the Client to wait this maximum number of seconds to obtain a connection.
// Zero or not specified means wait indefinitely.
func ConnectTimeout(timeout int) Option {
	return func(options *ClientOptions) error {
		if timeout < 0 {
			return errors.New("makroud: the maximum wait for a connection must be a positive number")
		}
		options.ConnectTimeout = timeout
		return nil
	}
}
