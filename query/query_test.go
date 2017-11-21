package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		TestName string
		Prepare  func() Query
		Expected Query
	}{
		{
			"nil",
			func() (q Query) {
				return q.GetFields()
			},
			Query{},
		},
		{
			"empty",
			func() Query {
				return New()
			},
			Query{},
		},
		{
			"fields and order",
			func() Query {
				data := New().
					Add("a", true).
					Add("b", false).
					Add("c", []string{"a", "b", "c"}).
					OrderBy("c", true)
				data.GetFields()
				return data
			},
			Query{
				"a":        true,
				"b":        false,
				"c":        []string{"a", "b", "c"},
				orderByKey: Ordering{"c", true},
			},
		},
		{
			"fields and order return only fields",
			func() Query {
				return New().
					Add("a", true).
					Add("b", false).
					Add("c", []string{"a", "b", "c"}).
					OrderBy("c", true).GetFields()
			},
			Query{
				"a": true,
				"b": false,
				"c": []string{"a", "b", "c"},
			},
		},
		{
			"adding parameters as well",
			func() Query {
				return New().
					Add("a", true).
					Add("b", false).
					Add("c", []string{"a", "b", "c"}).
					OrderBy("c", true).
					AddParam("param-1", 1).
					AddParam("param-2", 2).
					GetFields()
			},
			Query{
				"a": true,
				"b": false,
				"c": []string{"a", "b", "c"},
			},
		},
	}

	for _, test := range testTable {
		assert.EqualValues(test.Expected, test.Prepare(), test.TestName)
	}
}
