package dotnotation

import (
	"reflect"
	"strings"
	"testing"
)

func TestWildcardStruct(t *testing.T) {
	expected := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "supu",
			},
			"d": map[string]interface{}{
				"c": "supu",
			},
		},
	}
	m, err := CompileApplier("a.*.c", func(data interface{}) interface{} {
		if v, ok := data.(string); ok {
			return strings.ReplaceAll(v, "tupu", "supu")
		}
		return data
	})
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "tupu",
			},
			"d": map[string]interface{}{
				"c": "tupu",
			},
		},
	}
	m.Apply(data)
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("%v is not %v", data, expected)
	}
}

func TestWildcardSlice(t *testing.T) {
	expected := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"c": "supu",
			},
			map[string]interface{}{
				"c": "supu",
			},
		},
	}
	m, err := CompileApplier("a.*.c", func(data interface{}) interface{} {
		if v, ok := data.(string); ok {
			return strings.ReplaceAll(v, "tupu", "supu")
		}
		return data
	})
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"c": "tupu",
			},
			map[string]interface{}{
				"c": "tupu",
			},
		},
	}
	m.Apply(data)
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("%v is not %v", data, expected)
	}
}

func TestWildcardSliceExtract(t *testing.T) {
	expected := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"c": "tupu1",
			},
			map[string]interface{}{
				"c": "tupu2",
			},
		},
	}
	m, err := CompileExtractor("a.*.c")
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"c": "tupu1",
			},
			map[string]interface{}{
				"c": "tupu2",
			},
		},
	}
	res := m.Extract(data)
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("%v is not %v", data, expected)
	}
	if !reflect.DeepEqual(res, []interface{}{"tupu1", "tupu2"}) {
		t.Errorf("%v is not %v", data, expected)
	}
}

func TestIndexSlice(t *testing.T) {
	expected := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"c": "tupu",
			},
			map[string]interface{}{
				"c": "supu",
			},
		},
	}
	m, err := CompileApplier("a.1.c", func(data interface{}) interface{} {
		if v, ok := data.(string); ok {
			return strings.ReplaceAll(v, "tupu", "supu")
		}
		return data
	})
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"c": "tupu",
			},
			map[string]interface{}{
				"c": "tupu",
			},
		},
	}
	m.Apply(data)
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("%v is not %v", data, expected)
	}
}

func TestIndexMap(t *testing.T) {
	expected := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "tupu",
			},
			"1": map[string]interface{}{
				"c": "tupu",
			},
		},
	}
	m, err := CompileApplier("patata", func(data interface{}) interface{} {
		if v, ok := data.(string); ok {
			return strings.ReplaceAll(v, "tupu", "supu")
		}
		return data
	})
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "tupu",
			},
			"1": map[string]interface{}{
				"c": "tupu",
			},
		},
	}
	m.Apply(data)
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("%v is not %v", data, expected)
	}
}

func TestMapStruct(t *testing.T) {
	expected := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "supu",
			},
			"d": map[string]interface{}{
				"c": "tupu",
			},
		},
	}
	m, err := CompileApplier("a.b.c", func(data interface{}) interface{} {
		if v, ok := data.(string); ok {
			return strings.ReplaceAll(v, "tupu", "supu")
		}
		return data
	})
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "tupu",
			},
			"d": map[string]interface{}{
				"c": "tupu",
			},
		},
	}
	m.Apply(data)
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("%v is not %v", data, expected)
	}
}
