package makroud

import (
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// DefaultSelector defines the default selector alias.
const DefaultSelector = MasterSelector

// MasterSelector defines the master alias.
const MasterSelector = "master"

// SlaveSelector defines the slave alias.
const SlaveSelector = "slave"

// Selector contains a pool of drivers indexed by their name.
type Selector struct {
	mutex          sync.RWMutex
	store          *cache
	configurations map[string]*ClientOptions
	connections    map[string]Driver
}

func (selector *Selector) cache() *cache {
	return selector.store
}

func (selector *Selector) hasCache() bool {
	return selector.store != nil
}

// NewSelector returns a new selector containing a pool of drivers with given configuration.
func NewSelector(configurations map[string]*ClientOptions) (*Selector, error) {
	connections := map[string]Driver{}

	selector := &Selector{
		configurations: configurations,
		connections:    connections,
		store:          newCache(),
	}

	return selector, nil
}

// NewSelectorWithDriver returns a new selector containing the given connection.
func NewSelectorWithDriver(driver Driver) (*Selector, error) {
	selector := &Selector{
		configurations: map[string]*ClientOptions{},
		store:          driver.cache(),
		connections: map[string]Driver{
			DefaultSelector: driver,
		},
	}

	return selector, nil
}

// Using returns the underlying drivers if it's alias exists.
func (selector *Selector) Using(alias string) (Driver, error) {
	alias = strings.ToLower(alias)

	selector.mutex.RLock()
	connection, found := selector.connections[alias]
	selector.mutex.RUnlock()

	if found {
		return connection, nil
	}

	selector.mutex.Lock()
	defer selector.mutex.Unlock()

	connection, found = selector.connections[alias]
	if found {
		return connection, nil
	}

	for name, configuration := range selector.configurations {
		if alias == strings.ToLower(name) {

			connection, err := NewWithOptions(configuration)
			if err != nil {
				return nil, err
			}

			if selector.hasCache() {
				connection.setCache(selector.cache())
			}

			selector.connections[alias] = connection

			return connection, nil
		}
	}

	return nil, errors.Wrapf(ErrSelectorNotFoundConnection, "connection alias '%s' not found", alias)
}

// RetryAliases is an helper calling Retry with a list of aliases.
func (selector *Selector) RetryAliases(handler func(Driver) error, aliases ...string) error {
	drivers := []Driver{}

	for _, alias := range aliases {
		connection, err := selector.Using(alias)
		if err != nil {
			continue
		}

		drivers = append(drivers, connection)
	}

	return Retry(handler, drivers...)
}

// RetryMaster is an helper calling RetryAliases with a slave then a master connection.
func (selector *Selector) RetryMaster(handler func(Driver) error) error {
	return selector.RetryAliases(handler, SlaveSelector, MasterSelector)
}

// Close closes all drivers connections.
func (selector *Selector) Close() error {
	selector.mutex.Lock()
	defer selector.mutex.Unlock()

	failures := []error{}

	for alias, connection := range selector.connections {
		err := connection.Close()
		if err != nil {
			failures = append(failures, errors.Wrapf(err, "cannot close drivers connection for %s", alias))
		}
	}

	selector.connections = map[string]Driver{}

	if len(failures) > 0 {
		// TODO (novln): Add an observer to collect these errors.
		return failures[0]
	}

	return nil
}

// Ping checks if a connection is available.
func (selector *Selector) Ping() error {
	return selector.RetryMaster(func(driver Driver) error {
		return driver.Ping()
	})
}

// Retry execute given handler on several drivers until it succeeds on a connection.
func Retry(handler func(Driver) error, drivers ...Driver) (err error) {
	if len(drivers) == 0 {
		return errors.WithStack(ErrSelectorMissingRetryConnection)
	}

	for _, driver := range drivers {
		err = handler(driver)
		if err == nil {
			return nil
		}
	}

	return err
}
