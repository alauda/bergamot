package log_test

import (
	"testing"

	"github.com/alauda/bergamot/log"
)

func TestFFunction(t *testing.T) {

	type TestCase struct {
		Args     []interface{}
		Expected map[string]interface{}
	}

	table := []TestCase{
		{
			// simple case
			[]interface{}{"key", 1},
			map[string]interface{}{"key": 1},
		},
		{
			// extra key without value will be ignored
			[]interface{}{"key", 1, "extra"},
			map[string]interface{}{"key": 1},
		},
		{
			// multiple keys
			[]interface{}{"key", 1, "extra", "some"},
			map[string]interface{}{"key": 1, "extra": "some"},
		},
		{
			// keys that are not string
			[]interface{}{1, 2, true, "some"},
			map[string]interface{}{"1": 2, "true": "some"},
		},
	}

	for i, test := range table {
		result := log.F(test.Args...)
		if len(result) != len(test.Expected) {
			t.Errorf("%d -- lens are different expected: %d got %d", i, len(test.Expected), len(result))
			t.Fail()
		}
		for k, v := range test.Expected {
			if result[k] != v {
				t.Errorf("%d -- key %s and value %v are different from expected: %v", i, k, result[k], v)
				t.Fail()
			}
		}
	}
}
