package sqlxx_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum/format"

	"github.com/ulule/sqlxx"
)

type logger struct {
	logs chan string
}

func (e *logger) Log(query string, duration time.Duration) {
	e.logs <- query
}

var ErrLogTimeout = fmt.Errorf("logger timeout")
var EOL = "\n"

func (e *logger) read() (string, error) {
	select {
	case log := <-e.logs:
		return log, nil
	case <-time.After(500 * time.Millisecond):
		return "", ErrLogTimeout
	}
}

func TestLogger(t *testing.T) {
	logger := &logger{
		logs: make(chan string, 10),
	}
	Setup(t, sqlxx.WithLogger(logger))(func(driver sqlxx.Driver) {
		is := require.New(t)

		owl := &Owl{
			Name:         "Guacamowle",
			FeatherColor: "lavender",
			FavoriteFood: "Shrimps",
		}

		err := sqlxx.Save(driver, owl)
		is.NoError(err)
		expected := fmt.Sprint(
			`INSERT INTO wp_owl (favorite_food, feather_color, name) VALUES `,
			`('Shrimps', 'lavender', 'Guacamowle') RETURNING id;`, EOL,
		)

		log, err := logger.read()
		is.NoError(err)
		is.Equal(expected, log)

		owl.Name = "Nibbles"
		err = sqlxx.Save(driver, owl)
		is.NoError(err)
		expected = fmt.Sprint(
			`UPDATE wp_owl SET favorite_food = 'Shrimps', feather_color = 'lavender', name = 'Nibbles' `,
			`WHERE (id = `, format.Int(owl.ID), `);`, EOL,
		)

		log, err = logger.read()
		is.NoError(err)
		is.Equal(expected, log)

		sqlxx.Delete(driver, owl)
		is.NoError(err)
		expected = fmt.Sprint(`DELETE FROM wp_owl WHERE (id = `, format.Int(owl.ID), `);`, EOL)

		log, err = logger.read()
		is.NoError(err)
		is.Equal(expected, log)

	})
}
