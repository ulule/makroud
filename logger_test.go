package makroud_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum/v3/format"

	"github.com/ulule/makroud"
)

type logger struct {
	logs chan string
}

func (e *logger) Log(ctx context.Context, query string, duration time.Duration) {
	e.logs <- query
}

var ErrLogTimeout = fmt.Errorf("logger timeout")

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
	Setup(t, makroud.WithLogger(logger))(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		owl := &Owl{
			Name:         "Guacamowle",
			FeatherColor: "lavender",
			FavoriteFood: "Shrimps",
		}

		err := makroud.Save(ctx, driver, owl)
		is.NoError(err)
		expected := fmt.Sprint(
			`INSERT INTO ztp_owl (favorite_food, feather_color, group_id, name) VALUES `,
			`('Shrimps', 'lavender', NULL, 'Guacamowle') RETURNING id`,
		)

		log, err := logger.read()
		is.NoError(err)
		is.Equal(expected, log)

		owl.Name = "Nibbles"
		err = makroud.Save(ctx, driver, owl)
		is.NoError(err)
		expected = fmt.Sprint(
			`UPDATE ztp_owl SET favorite_food = 'Shrimps', feather_color = 'lavender', group_id = NULL, `,
			`name = 'Nibbles' WHERE (id = `, format.Int(owl.ID), `)`,
		)

		log, err = logger.read()
		is.NoError(err)
		is.Equal(expected, log)

		err = makroud.Delete(ctx, driver, owl)
		is.NoError(err)
		expected = fmt.Sprint(`DELETE FROM ztp_owl WHERE (id = `, format.Int(owl.ID), `)`)

		log, err = logger.read()
		is.NoError(err)
		is.Equal(expected, log)

	})
}
