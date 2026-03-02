package dotnotation

import (
	"strings"
	"testing"
)

func BenchmarkDotNotationStrategies(b *testing.B) { // skipcq: GO-R1005
	op := func(data interface{}) interface{} {
		if v, ok := data.(string); ok {
			return strings.ReplaceAll(v, "tupu", "supu")
		}
		return data
	}
	b.Run("simple maps", func(b *testing.B) {
		m, err := CompileApplier("a.b.c", op)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
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
		}
	})
	b.Run("long map structure", func(b *testing.B) {
		m, err := CompileApplier("a.b.c.a.b.d.a.b.c", op)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			data := map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": map[string]interface{}{
							"a": map[string]interface{}{
								"b": map[string]interface{}{
									"c": "tupu",
								},
								"d": map[string]interface{}{
									"c": map[string]interface{}{
										"a": map[string]interface{}{
											"b": map[string]interface{}{
												"c": map[string]interface{}{
													"a": map[string]interface{}{
														"b": map[string]interface{}{
															"c": "tupu",
														},
														"d": map[string]interface{}{
															"c": "tupu",
														},
													},
												},
											},
											"d": map[string]interface{}{
												"c": "tupu",
											},
										},
									},
								},
							},
						},
					},
					"d": map[string]interface{}{
						"c": "tupu",
					},
				},
			}
			m.Apply(data)
		}
	})
	b.Run("simple maps with wildcard", func(b *testing.B) {
		m, err := CompileApplier("a.*.c", op)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
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
		}
	})
	b.Run("long map structure with wildcard", func(b *testing.B) {
		m, err := CompileApplier("a.*.c.a.d.c.a.b.c.a.*.c", op)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			data := map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": map[string]interface{}{
							"a": map[string]interface{}{
								"b": map[string]interface{}{
									"c": "tupu",
								},
								"d": map[string]interface{}{
									"c": map[string]interface{}{
										"a": map[string]interface{}{
											"b": map[string]interface{}{
												"c": map[string]interface{}{
													"a": map[string]interface{}{
														"b": map[string]interface{}{
															"c": "tupu",
														},
														"d": map[string]interface{}{
															"c": "tupu",
														},
													},
												},
											},
											"d": map[string]interface{}{
												"c": "tupu",
											},
										},
									},
								},
							},
						},
					},
					"d": map[string]interface{}{
						"c": "tupu",
					},
				},
			}
			m.Apply(data)
		}
	})
	b.Run("simple maps+slice", func(b *testing.B) {
		m, err := CompileApplier("a.1.c", op)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
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
		}
	})
	b.Run("simple maps+slice with wildcard", func(b *testing.B) {
		m, err := CompileApplier("a.*.c", op)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
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
		}
	})
	b.Run("long maps+slice structure with wildcards", func(b *testing.B) {
		m, err := CompileApplier("a.*.0.0.d.c.a.b.c.a.*.c", op)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			data := map[string]interface{}{
				"a": []interface{}{
					[]interface{}{
						[]interface{}{
							map[string]interface{}{
								"b": map[string]interface{}{
									"c": "tupu",
								},
								"d": map[string]interface{}{
									"c": map[string]interface{}{
										"a": map[string]interface{}{
											"b": map[string]interface{}{
												"c": map[string]interface{}{
													"a": []interface{}{
														map[string]interface{}{
															"c": "tupu",
														},
														map[string]interface{}{
															"c": "tupu",
														},
													},
												},
											},
											"d": map[string]interface{}{
												"c": "tupu",
											},
										},
									},
								},
							},
						},
					},
					map[string]interface{}{
						"c": "tupu",
					},
				},
			}
			m.Apply(data)
		}
	})
}
