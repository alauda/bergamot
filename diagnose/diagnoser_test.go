package diagnose_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/alauda/bergamot/diagnose"
)

func TestComponentReportMarshal(t *testing.T) {
	// obj := &diagnose.ComponentReport{
	// 	Status:     diagnose.StatusOK,
	// 	Name:       "some",
	// 	Message:    "msg",
	// 	Suggestion: "sug",
	// 	Latency:    time.Second * 10,
	// }

	type TestCase struct {
		Name   string
		Obj    *diagnose.ComponentReport
		Result []byte
		Err    error
	}

	table := []TestCase{
		{
			"100000 microseconds",
			&diagnose.ComponentReport{
				Status:     diagnose.StatusOK,
				Name:       "some",
				Message:    "msg",
				Suggestion: "sug",
				Latency:    time.Microsecond * 100000,
			},
			[]byte(
				`{` +
					`"status":"OK",` +
					`"name":"some",` +
					`"message":"msg",` +
					`"suggestion":"sug",` +
					`"latency":"100ms"` +
					`}`,
			),
			nil,
		},
		{
			"100 microseconds",
			&diagnose.ComponentReport{
				Status:     diagnose.StatusOK,
				Name:       "some",
				Message:    "msg",
				Suggestion: "sug",
				Latency:    time.Microsecond * 100,
			},
			[]byte(
				`{` +
					`"status":"OK",` +
					`"name":"some",` +
					`"message":"msg",` +
					`"suggestion":"sug",` +
					`"latency":"0.10ms"` +
					`}`,
			),
			nil,
		},
		{
			"150 milliseconds",
			&diagnose.ComponentReport{
				Status:     diagnose.StatusOK,
				Name:       "some",
				Message:    "msg",
				Suggestion: "sug",
				Latency:    time.Millisecond * 150,
			},
			[]byte(
				`{` +
					`"status":"OK",` +
					`"name":"some",` +
					`"message":"msg",` +
					`"suggestion":"sug",` +
					`"latency":"150ms"` +
					`}`,
			),
			nil,
		},
		{
			"10 seconds",
			&diagnose.ComponentReport{
				Status:     diagnose.StatusOK,
				Name:       "some",
				Message:    "msg",
				Suggestion: "sug",
				Latency:    time.Second * 10,
			},
			[]byte(
				`{` +
					`"status":"OK",` +
					`"name":"some",` +
					`"message":"msg",` +
					`"suggestion":"sug",` +
					`"latency":"10s"` +
					`}`,
			),
			nil,
		},
	}

	for _, test := range table {
		t.Run(test.Name, func(t *testing.T) {
			result, err := test.Obj.MarshalJSON()
			if err == nil && test.Err != nil {
				t.Error("should return error but didn't")
			}
			if err != nil && test.Err == nil {
				t.Error("should not return error but it did")
			}
			if !reflect.DeepEqual(test.Result, result) {
				t.Error("generated json is not equal:", string(test.Result), "!=", string(result))
			}
		})

	}
}
