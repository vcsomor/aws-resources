package lister

import (
	"fmt"
	"reflect"
)

type Flattener interface {
	Flatten(data map[string]any) map[string]any
}

type objectFlattener struct {
}

var _ Flattener = (*objectFlattener)(nil)

func NewObjectFlattener() Flattener {
	return &objectFlattener{}
}

func (r *objectFlattener) Flatten(data map[string]any) map[string]any {
	return flattenMap("", data)
}

func flattenSlice(root string, data []any) map[string]any {
	if len(data) == 0 {
		return map[string]any{
			keyOf(root, "ARRAY"): nil,
		}
	}

	ret := map[string]any{}
	for i, dt := range data {
		indexTag := fmt.Sprintf("ARRAY$%d", i)
		t := reflect.TypeOf(dt)
		if t == nil {
			ret[keyOf(root, indexTag)] = nil
			continue
		}

		switch t.Kind() {
		case reflect.Slice:
			sl := dt.([]any)
			flattened := flattenSlice(indexTag, sl)
			for childKey, childValue := range flattened {
				ret[keyOf(root, childKey)] = childValue
			}

		case reflect.Map:
			flattened := flattenMap(indexTag, dt.(map[string]any))
			for childKey, childValue := range flattened {
				ret[keyOf(root, childKey)] = childValue
			}

		default:
			ret[keyOf(root, indexTag)] = dt
		}
	}
	return ret
}

func flattenMap(root string, data map[string]any) map[string]any {
	if data == nil {
		return nil
	}

	ret := map[string]any{}

	for k, v := range data {
		t := reflect.TypeOf(v)
		if t == nil {
			ret[keyOf(root, k)] = nil
			continue
		}

		switch t.Kind() {
		case reflect.Slice:
			sl := v.([]any)
			flattened := flattenSlice(k, sl)
			for childKey, childValue := range flattened {
				ret[keyOf(root, childKey)] = childValue
			}

		case reflect.Map:
			flattened := flattenMap(k, v.(map[string]any))
			for childKey, childValue := range flattened {
				ret[keyOf(root, childKey)] = childValue
			}

		default:
			ret[keyOf(root, k)] = v
		}
	}

	return ret
}

func keyOf(root, current string) string {
	if root == "" {
		return current
	}
	return root + "." + current
}
