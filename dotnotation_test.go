package dotnotation

import (
	"reflect"
	"strings"
	"testing"
)

func TestApply(t *testing.T) {
	t.Run("TestWildcardStruct", func(t *testing.T) {
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
	})

	t.Run("TestNumericMap", func(t *testing.T) {
		m, err := CompileApplier("a.1.1", func(data interface{}) interface{} {
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
					"1": "tupu",
				},
				map[string]interface{}{
					"1": "tupu",
				},
			},
		}
		m.Apply(data)
		expected := map[string]interface{}{
			"a": []interface{}{
				map[string]interface{}{
					"1": "tupu",
				},
				map[string]interface{}{
					"1": "supu",
				},
			},
		}
		if !reflect.DeepEqual(data, expected) {
			t.Errorf("%v is not %v", data, expected)
		}
	})

	t.Run("TestIndexSlice", func(t *testing.T) {
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
	})

	t.Run("TestIndexMap", func(t *testing.T) {
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
	})

	t.Run("TestMapStruct", func(t *testing.T) {
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
	})

	t.Run("TestNumeric", func(t *testing.T) {
		expected := map[string]interface{}{
			"a": []interface{}{
				[]interface{}{
					1,
					"supu",
				},
				map[string]interface{}{
					"c": "tupu",
				},
			},
		}
		m, err := CompileApplier("a.0.1", func(data interface{}) interface{} {
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
				[]interface{}{
					1,
					"tupu",
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
	})

	t.Run("TestWildcardNumeric", func(t *testing.T) {
		expected := map[string]interface{}{
			"a": []interface{}{
				[]interface{}{
					1,
					"supu",
				},
				map[string]interface{}{
					"c": "tupu",
				},
			},
		}
		m, err := CompileApplier("a.*.1", func(data interface{}) interface{} {
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
				[]interface{}{
					1,
					"tupu",
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
	})

	t.Run("TestWildcardSlice", func(t *testing.T) {
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
	})

	t.Run("TestMapWildcard", func(t *testing.T) {
		expected := map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"c": "tupu",
					"d": "tupu",
				},
				"d": map[string]interface{}{
					"c": "supu",
					"d": "supu",
				},
			},
		}
		m, err := CompileApplier("a.d.*", func(data interface{}) interface{} {
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
					"d": "tupu",
				},
				"d": map[string]interface{}{
					"c": "tupu",
					"d": "tupu",
				},
			},
		}
		m.Apply(data)
		if !reflect.DeepEqual(data, expected) {
			t.Errorf("%v is not %v", data, expected)
		}
	})

	t.Run("TestMapWildcardWildcard", func(t *testing.T) {
		expected := map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"c": "supu",
					"d": "supu",
				},
				"d": map[string]interface{}{
					"c": "supu",
					"d": "supu",
				},
			},
		}
		m, err := CompileApplier("a.*.*", func(data interface{}) interface{} {
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
					"d": "tupu",
				},
				"d": map[string]interface{}{
					"c": "tupu",
					"d": "tupu",
				},
			},
		}
		m.Apply(data)
		if !reflect.DeepEqual(data, expected) {
			t.Errorf("%v is not %v", data, expected)
		}
	})
}

func TestExtract(t *testing.T) {
	t.Run("TestWildcardStruct", func(t *testing.T) {
		m, err := CompileExtractor("a.*.c")
		if err != nil {
			t.Fatal(err)
		}

		data := map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"c": "tupu",
				},
				"d": map[string]interface{}{
					"c": "supu",
				},
			},
		}
		res := m.Extract(data)
		expected := []interface{}{"tupu", "supu"}
		if !equalNoOrder(res, expected) {
			t.Errorf("%v is not %v", res, expected)
		}
	})

	t.Run("TestNumericMap", func(t *testing.T) {
		m, err := CompileExtractor("a.1.1")
		if err != nil {
			t.Fatal(err)
		}

		data := map[string]interface{}{
			"a": []interface{}{
				map[string]interface{}{
					"1": "tupu",
				},
				map[string]interface{}{
					"1": "supu",
				},
			},
		}
		res := m.Extract(data)
		expected := []interface{}{"supu"}
		if !reflect.DeepEqual(res, expected) {
			t.Errorf("%v is not %v", data, expected)
		}
	})

	t.Run("TestWildcardNumericMap", func(t *testing.T) {
		m, err := CompileExtractor("a.*.1")
		if err != nil {
			t.Fatal(err)
		}

		data := map[string]interface{}{
			"a": []interface{}{
				map[string]interface{}{
					"0": "tupu",
				},
				map[string]interface{}{
					"1": "supu",
				},
			},
		}
		res := m.Extract(data)
		expected := []interface{}{"supu"}
		if !reflect.DeepEqual(res, expected) {
			t.Errorf("%v is not %v", data, expected)
		}
	})

	t.Run("TestNoMatch", func(t *testing.T) {
		m, err := CompileExtractor("patata")
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
		res := m.Extract(data)
		expected := []interface{}{}
		if len(res) != 0 {
			t.Errorf("%v is not %v", res, expected)
		}
	})

	t.Run("TestMapStruct", func(t *testing.T) {
		m, err := CompileExtractor("a.b.c")
		if err != nil {
			t.Fatal(err)
		}

		data := map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"c": "tupu",
				},
				"d": map[string]interface{}{
					"c": "supu",
				},
			},
		}
		res := m.Extract(data)
		expected := []interface{}{"tupu"}
		if !reflect.DeepEqual(res, expected) {
			t.Errorf("%v is not %v", res, expected)
		}
	})

	t.Run("TestNumeric", func(t *testing.T) {
		m, err := CompileExtractor("a.0.1")
		if err != nil {
			t.Fatal(err)
		}

		data := map[string]interface{}{
			"a": []interface{}{
				[]interface{}{
					1,
					"tupu",
				},
				map[string]interface{}{
					"c": "supu",
				},
			},
		}
		v := m.Extract(data)
		expected := []interface{}{"tupu"}
		if !reflect.DeepEqual(v, []interface{}{"tupu"}) {
			t.Errorf("%v is not %v", v, expected)
		}
	})

	t.Run("TestWildcardSlice", func(t *testing.T) {
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
		expected := []interface{}{"tupu1", "tupu2"}
		if !reflect.DeepEqual(res, expected) {
			t.Errorf("%v is not %v", data, expected)
		}
	})

	t.Run("TestWildcardLastStepMap", func(t *testing.T) {
		m, err := CompileExtractor("a.1.*")
		if err != nil {
			t.Fatal(err)
		}

		data := map[string]interface{}{
			"a": []interface{}{
				map[string]interface{}{
					"c": "tupu1",
				},
				map[string]interface{}{
					"c": "tupu",
					"b": "supu",
				},
			},
		}
		res := m.Extract(data)
		expected := []interface{}{"tupu", "supu"}
		if !equalNoOrder(res, expected) {
			t.Errorf("%v is not %v", data, expected)
		}
	})

	t.Run("TestWildcardLastStepSlice", func(t *testing.T) {
		m, err := CompileExtractor("a.1.*")
		if err != nil {
			t.Fatal(err)
		}

		data := map[string]interface{}{
			"a": []interface{}{
				map[string]interface{}{
					"c": "tupu1",
				},
				[]interface{}{
					"tupu",
					"supu",
					2,
				},
			},
		}
		res := m.Extract(data)
		expected := []interface{}{"tupu", "supu", 2}
		if !reflect.DeepEqual(res, expected) {
			t.Errorf("%v is not %v", data, expected)
		}
	})
}

func equalNoOrder(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	counts := make(map[interface{}]int, len(a))
	for _, v := range a {
		counts[v]++
	}
	for _, v := range b {
		counts[v]--
		if counts[v] < 0 {
			return false
		}
	}
	return true
}
