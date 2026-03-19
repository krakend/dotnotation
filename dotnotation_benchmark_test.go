package dotnotation

import (
	"strings"
	"testing"
)

func BenchmarkDotNotationStrategiesApplier(b *testing.B) { // skipcq: GO-R1005
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
		m, err := CompileApplier("a.*.0.0.d.c.a.b.c.a.*.*", op)
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
															"d": "tupu",
															"e": "tupu",
															"f": "tupu",
														},
														[]interface{}{
															"tupu",
															"supu",
															"sapu",
															"tapu",
														},
														map[string]interface{}{
															"c": "tupu",
															"d": "tupu",
															"e": "tupu",
															"f": "tupu",
														},
														[]interface{}{
															"tupu",
															"supu",
															"sapu",
															"tapu",
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

var result []interface{}

func BenchmarkDotNotationStrategiesExtractor(b *testing.B) { // skipcq: GO-R1005
	b.Run("simple maps extractor", func(b *testing.B) {
		m, err := CompileExtractor("a.b.c")
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		var tmp []interface{}
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
			tmp = m.Extract(data)
		}
		result = tmp
	})
	b.Run("long map structure extractor", func(b *testing.B) {
		m, err := CompileExtractor("a.b.c.a.b.d.a.b.c")
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		var tmp []interface{}
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
			tmp = m.Extract(data)
		}
		result = tmp
	})
	b.Run("simple maps with wildcard extractor", func(b *testing.B) {
		m, err := CompileExtractor("a.*.c")
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		var tmp []interface{}
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
			tmp = m.Extract(data)
		}
		result = tmp
	})
	b.Run("long map structure with wildcard extractor", func(b *testing.B) {
		m, err := CompileExtractor("a.*.c.a.d.c.a.b.c.a.*.c")
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		var tmp []interface{}
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
													"a": []interface{}{
														map[string]interface{}{
															"c": "tupu",
															"d": "tupu",
															"e": "tupu",
															"f": "tupu",
														},
														[]interface{}{
															"tupu",
															"supu",
															"sapu",
															"tapu",
														},
														map[string]interface{}{
															"c": "tupu",
															"d": "tupu",
															"e": "tupu",
															"f": "tupu",
														},
														[]interface{}{
															"tupu",
															"supu",
															"sapu",
															"tapu",
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
			tmp = m.Extract(data)
		}
		result = tmp
	})
	b.Run("simple maps+slice extractor", func(b *testing.B) {
		m, err := CompileExtractor("a.1.c")
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		var tmp []interface{}
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
			tmp = m.Extract(data)
		}
		result = tmp
	})
	b.Run("simple maps+slice with wildcard extractor", func(b *testing.B) {
		m, err := CompileExtractor("a.*.c")
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		var tmp []interface{}
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
			tmp = m.Extract(data)
		}
		result = tmp
	})
	b.Run("long maps+slice structure with wildcards extractor", func(b *testing.B) {
		m, err := CompileExtractor("a.*.0.0.d.c.a.b.c.a.*.*")
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		var tmp []interface{}
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
			tmp = m.Extract(data)
		}
		result = tmp
	})
}
