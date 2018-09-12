package reflectx_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx/reflectx"
)

func TestReflectx_Walker(t *testing.T) {
	is := require.New(t)

	type F struct {
		Foobar int64
	}

	type E struct {
		F F
	}

	type D struct {
		E []E
	}

	type C struct {
		D D
	}

	type B struct {
		C *[]C
	}

	type A struct {
		B []B
	}

	values := []A{
		{
			B: []B{
				{
					C: &[]C{
						{
							D: D{
								E: []E{
									{
										F: F{
											Foobar: 1,
										},
									},
									{
										F: F{
											Foobar: 3,
										},
									},
									{
										F: F{
											Foobar: 5,
										},
									},
								},
							},
						},
						{
							D: D{
								E: []E{
									{
										F: F{
											Foobar: 2,
										},
									},
									{
										F: F{
											Foobar: 4,
										},
									},
								},
							},
						},
						{
							D: D{
								E: []E{},
							},
						},
					},
				},
				{
					C: &[]C{
						{
							D: D{},
						},
						{
							D: D{
								E: []E{
									{
										F: F{
											Foobar: 10,
										},
									},
									{
										F: F{
											Foobar: 11,
										},
									},
								},
							},
						},
						{
							D: D{
								E: []E{
									{
										F: F{
											Foobar: 12,
										},
									},
									{
										F: F{
											Foobar: 13,
										},
									},
									{
										F: F{
											Foobar: 14,
										},
									},
								},
							},
						},
						{
							D: D{
								E: []E{
									{
										F: F{
											Foobar: 15,
										},
									},
									{
										F: F{
											Foobar: 16,
										},
									},
									{
										F: F{
											Foobar: 17,
										},
									},
									{
										F: F{
											Foobar: 18,
										},
									},
									{
										F: F{
											Foobar: 19,
										},
									},
									{
										F: F{
											Foobar: 20,
										},
									},
								},
							},
						},
					},
				},
				{
					C: nil,
				},
			},
		},
		{
			B: []B{
				{
					C: &[]C{
						{
							D: D{
								E: []E{
									{
										F: F{
											Foobar: 21,
										},
									},
									{
										F: F{
											Foobar: 23,
										},
									},
									{
										F: F{
											Foobar: 25,
										},
									},
								},
							},
						},
						{
							D: D{
								E: []E{
									{
										F: F{
											Foobar: 22,
										},
									},
									{
										F: F{
											Foobar: 24,
										},
									},
								},
							},
						},
						{
							D: D{
								E: []E{},
							},
						},
					},
				},
				{
					C: &[]C{
						{
							D: D{},
						},
						{
							D: D{
								E: []E{
									{
										F: F{
											Foobar: 30,
										},
									},
									{
										F: F{
											Foobar: 31,
										},
									},
								},
							},
						},
						{
							D: D{
								E: []E{
									{
										F: F{
											Foobar: 32,
										},
									},
									{
										F: F{
											Foobar: 33,
										},
									},
									{
										F: F{
											Foobar: 34,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	expected := []F{
		{
			Foobar: 1,
		},
		{
			Foobar: 3,
		},
		{
			Foobar: 5,
		},
		{
			Foobar: 2,
		},
		{
			Foobar: 4,
		},
		{
			Foobar: 10,
		},
		{
			Foobar: 11,
		},
		{
			Foobar: 12,
		},
		{
			Foobar: 13,
		},
		{
			Foobar: 14,
		},
		{
			Foobar: 15,
		},
		{
			Foobar: 16,
		},
		{
			Foobar: 17,
		},
		{
			Foobar: 18,
		},
		{
			Foobar: 19,
		},
		{
			Foobar: 20,
		},
		{
			Foobar: 21,
		},
		{
			Foobar: 22,
		},
		{
			Foobar: 23,
		},
		{
			Foobar: 24,
		},
		{
			Foobar: 25,
		},
		{
			Foobar: 30,
		},
		{
			Foobar: 31,
		},
		{
			Foobar: 32,
		},
		{
			Foobar: 33,
		},
		{
			Foobar: 34,
		},
	}

	walker := reflectx.NewWalker(&values)
	defer walker.Close()

	err := walker.Find("B.C.D.E.F", func(values interface{}) error {
		array, ok := (values).(*[]*F)
		is.True(ok)
		is.NotNil(array)
		is.NotEmpty((*array))
		is.Len((*array), len(expected))

		for i, value := range *array {
			message := fmt.Sprintf("#%d value: %v", i, value)
			is.NotNil(value, message)
			is.Contains(expected, (*value), message)
		}

		return nil
	})
	is.NoError(err)

}
