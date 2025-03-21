package main

import (
	"bytes"
	"testing"

	"github.com/matryer/is"
)

func TestRender(t *testing.T) {
	cases := map[string]struct {
		input    Data
		expected string
	}{
		"basic events": {
			input: Data{
				UnionName: "eventTypes",
				TypeName:  "EventType",
				Fields: []Field{
					{
						Name: "CREATED",
						Doc:  "Thing was created",
					},
				},
			},
			expected: `const eventTypes = [
	// Thing was created
	"CREATED",
] as const;
type EventType = typeof eventTypes[number];
`,
		},
		"no fields": {
			input: Data{
				UnionName: "eventTypes",
				TypeName:  "EventType",
				Fields:    []Field{},
			},
			expected: `const eventTypes = [
] as const;
type EventType = typeof eventTypes[number];
`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			result := new(bytes.Buffer)
			r := NewRenderer(result)

			err := r.Render(tc.input)
			is.NoErr(err)
			is.Equal(tc.expected, result.String())
		})
	}
}
