package gcode

import (
	"testing"

	"github.com/go-test/deep"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		description string
		line        string
		expected    instruction
	}{
		{
			"empty line",
			"",
			instruction{},
		},
		{
			"command only",
			"G90",
			instruction{
				cmd: "G90",
			},
		},
		{
			"comment only",
			"; test",
			instruction{
				comment: "; test",
			},
		},
		{
			"1 parameter",
			"G0 F3000",
			instruction{
				cmd:    "G0",
				params: []param{{"F", 3000}},
			},
		},
		{
			"2 parameters",
			"G0 F3000 X2.5",
			instruction{
				cmd:    "G0",
				params: []param{{"F", 3000}, {"X", 2.5}},
			},
		},
		{
			"2 negative parameters",
			"G0 F-3000 X-2.5",
			instruction{
				cmd:    "G0",
				params: []param{{"F", -3000}, {"X", -2.5}},
			},
		},
		{
			"2 negative parameters with comment",
			"G0 F-3000 X-2.5 ; commenty stuff",
			instruction{
				cmd:     "G0",
				params:  []param{{"F", -3000}, {"X", -2.5}},
				comment: "; commenty stuff",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			i := parse(testCase.line)

			if diff := deep.Equal(i, testCase.expected); diff != nil {
				t.Error(diff)
				return
			}
		})
	}
}

func init() {
	deep.CompareUnexportedFields = true
}
