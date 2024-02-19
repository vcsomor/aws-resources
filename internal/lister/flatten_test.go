package lister

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlatten(t *testing.T) {
	assert.Equal(t,
		map[string]any(nil),
		NewObjectFlattener().Flatten(nil))

	assert.Equal(t,
		map[string]any{},
		NewObjectFlattener().Flatten(map[string]any{}))

	assert.Equal(t,
		map[string]any{
			"field0":       "value0",
			"field1":       123,
			"field2":       nil,
			"field3.ARRAY": nil,

			"field4.ARRAY$0":               "value0",
			"field4.ARRAY$1":               123,
			"field4.ARRAY$2":               nil,
			"field4.ARRAY$3.ARRAY":         nil,
			"field4.ARRAY$4.ARRAY$0":       "value0",
			"field4.ARRAY$4.ARRAY$1":       123,
			"field4.ARRAY$4.ARRAY$2":       nil,
			"field4.ARRAY$4.ARRAY$3.ARRAY": nil,
			"field4.ARRAY$4.ARRAY$4.f0":    "value0",
			"field4.ARRAY$4.ARRAY$4.f1":    123,
			"field4.ARRAY$5.fx":            "value-x",
			"field4.ARRAY$5.fy":            10,

			"filed5.child-field0":       "value0",
			"filed5.child-field1":       123,
			"filed5.child-field2":       nil,
			"filed5.child-field3.ARRAY": nil,

			"filed5.child-field4.ARRAY$0":         "value0",
			"filed5.child-field4.ARRAY$1":         500,
			"filed5.child-field4.ARRAY$2":         nil,
			"filed5.child-field4.ARRAY$3.ARRAY":   nil,
			"filed5.child-field4.ARRAY$4.ARRAY$0": "child-value0",
			"filed5.child-field4.ARRAY$4.ARRAY$1": 1000,
			"filed5.child-field4.ARRAY$4.ARRAY$2": nil,

			"filed5.child-field5.child-field0":       "value0",
			"filed5.child-field5.child-field1":       123,
			"filed5.child-field5.child-field2":       nil,
			"filed5.child-field5.child-field3.ARRAY": nil,
		},
		NewObjectFlattener().Flatten(map[string]any{
			"field0": "value0",
			"field1": 123,
			"field2": nil,
			"field3": []any{},
			"field4": []any{
				"value0",
				123,
				nil,
				[]any{},
				[]any{
					"value0",
					123,
					nil,
					[]any{},
					map[string]any{
						"f0": "value0",
						"f1": 123,
					},
				},
				map[string]any{
					"fx": "value-x",
					"fy": 10,
				},
			},

			"filed5": map[string]any{
				"child-field0": "value0",
				"child-field1": 123,
				"child-field2": nil,

				"child-field3": []any{},
				"child-field4": []any{
					"value0",
					500,
					nil,
					[]any{},
					[]any{
						"child-value0",
						1000,
						nil,
					},
				},
				"child-field5": map[string]any{
					"child-field0": "value0",
					"child-field1": 123,
					"child-field2": nil,
					"child-field3": []any{},
				},
			},
		}))
}
