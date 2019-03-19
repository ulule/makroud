package snaker_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/makroud/snaker"
)

func TestSnaker_CamelToSnake(t *testing.T) {
	is := require.New(t)

	scenarios := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "",
		},
		{
			input:    "One",
			expected: "one",
		},
		{
			input:    "ONE",
			expected: "o_n_e",
		},
		{
			input:    "ID",
			expected: "id",
		},
		{
			input:    "i",
			expected: "i",
		},
		{
			input:    "ThisHasToBeConvertedCorrectlyID",
			expected: "this_has_to_be_converted_correctly_id",
		},
		{
			input:    "ThisIDIsFine",
			expected: "this_id_is_fine",
		},
		{
			input:    "ThisHTTPSConnection",
			expected: "this_https_connection",
		},
		{
			input:    "HelloHTTPSConnectionID",
			expected: "hello_https_connection_id",
		},
		{
			input:    "HTTPSID",
			expected: "https_id",
		},
		{
			input:    "OAuthClient",
			expected: "oauth_client",
		},
		{
			input:    "UserID",
			expected: "user_id",
		},
	}

	for _, scenario := range scenarios {

		actual := snaker.CamelToSnake(scenario.input)
		is.Equal(scenario.expected, actual)
	}
}

func TestSnaker_SnakeToCamel(t *testing.T) {
	is := require.New(t)

	scenarios := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "",
		},
		{
			input:    "potato_",
			expected: "Potato",
		},
		{
			input:    "this_has_to_be_uppercased",
			expected: "ThisHasToBeUppercased",
		},
		{
			input:    "this_is_an_id",
			expected: "ThisIsAnID",
		},
		{
			input:    "this_is_an_identifier",
			expected: "ThisIsAnIdentifier",
		},
		{
			input:    "id",
			expected: "ID",
		},
		{
			input:    "oauth_client",
			expected: "OAuthClient",
		},
		{
			input:    "id_me_plz",
			expected: "IDMePlz",
		},
		{
			input:    "user_id",
			expected: "UserID",
		},
	}

	for _, scenario := range scenarios {
		actual := snaker.SnakeToCamel(scenario.input)
		is.Equal(scenario.expected, actual)
	}
}

func TestSnaker_SnakeToCamelLower(t *testing.T) {
	is := require.New(t)

	scenarios := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "",
		},
		{
			input:    "potato_",
			expected: "potato",
		},
		{
			input:    "this_has_to_be_uppercased",
			expected: "thisHasToBeUppercased",
		},
		{
			input:    "this_is_an_id",
			expected: "thisIsAnID",
		},
		{
			input:    "this_is_an_identifier",
			expected: "thisIsAnIdentifier",
		},
		{
			input:    "id",
			expected: "id",
		},
		{
			input:    "id_me_plz",
			expected: "idMePlz",
		},
		{
			input:    "user_id",
			expected: "userID",
		},
	}

	for _, scenario := range scenarios {
		actual := snaker.SnakeToCamelLower(scenario.input)
		is.Equal(scenario.expected, actual)
	}
}
