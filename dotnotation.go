package dotnotation

import (
	"errors"
	"strconv"
	"strings"
)

type fieldType uint8

const (
	stringType fieldType = iota
	numericType
	wildcardType
)

type step struct {
	ft    fieldType
	key   string
	index int
	o     func(interface{}) interface{}
}

// Applier holds the compiled instructions to apply an operation to the compiled dotnotation path
type Applier struct {
	e *dotNotation
}

// Apply applies an operation to the compiled dotnotation path of the given data
func (a *Applier) Apply(v interface{}) {
	a.e.Evaluate(v)
}

// Extractor holds the compiled instructions to extract the compiled dotnotation path values
type Extractor struct {
	e *dotNotation
}

// Extract extracts the values in the compiled dotnotation path of the given data
func (e *Extractor) Extract(v interface{}) []interface{} {
	return e.e.Evaluate(v)
}

// CompileExtractor compiles the dotnotation expression and returns an Extractor
func CompileExtractor(expr string) (*Extractor, error) {
	e, err := compile(expr, nil)
	return &Extractor{e: e}, err
}

// CompileApplier compiles the dotnotation expression and returns an Applier
func CompileApplier(expr string, op func(interface{}) interface{}) (*Applier, error) {
	e, err := compile(expr, op)
	return &Applier{e: e}, err
}

func compile(expr string, op func(interface{}) interface{}) (*dotNotation, error) {
	if expr == "" {
		return nil, errors.New("cannot compile empty expr")
	}

	parts := strings.Split(expr, ".")
	steps := make([]step, 0, len(parts))
	wc := 0
	for _, p := range parts {
		idx, err := strconv.Atoi(p)
		isIndex := err == nil && idx >= 0

		switch {
		case p == "":
			return nil, errors.New("found empty field on expresion " + expr)
		case p == "*":
			wc++
			steps = append(steps, step{ft: wildcardType})
		case isIndex:
			steps = append(steps, step{ft: numericType, index: idx, key: p})
		default:
			steps = append(steps, step{ft: stringType, key: p})
		}
	}

	steps[len(steps)-1].o = op

	return &dotNotation{
		steps: steps,
		wc:    wc,
	}, nil
}

type dotNotation struct {
	steps []step
	wc    int
}

func (e *dotNotation) Evaluate(data interface{}) []interface{} {
	current := []interface{}{data}
	next := make([]interface{}, 0, 4*e.wc+1)

	for _, step := range e.steps {
		next = next[:0]

		switch step.ft {
		case stringType:
			next = stringTraverse(current, step, next)
		case numericType:
			next = numericTraverse(current, step, next)
		case wildcardType:
			next = wildcardTraverse(current, step, next)
		}

		if len(next) == 0 {
			return next
		}

		current, next = next, current
	}
	return current
}

func wildcardTraverse(current []interface{}, step step, next []interface{}) []interface{} {
	for _, n := range current {
		switch v := n.(type) {
		case []interface{}:
			if step.o != nil {
				for i := range v {
					v[i] = step.o(v[i])
				}
				continue
			}
			next = append(next, v...)
		case map[string]interface{}:
			if step.o != nil {
				for i := range v {
					v[i] = step.o(v[i])
				}
				continue
			}
			for i := range v {
				next = append(next, v[i])
			}
		}
	}
	return next
}

func numericTraverse(current []interface{}, step step, next []interface{}) []interface{} {
	for _, n := range current {
		if m, ok := n.(map[string]interface{}); ok {
			if v, exists := m[step.key]; exists {
				if step.o != nil {
					m[step.key] = step.o(v)
					continue
				}
				next = append(next, v)
			}
			continue
		}
		if arr, ok := n.([]interface{}); ok {
			if step.index < len(arr) {
				if step.o != nil {
					arr[step.index] = step.o(arr[step.index])
					continue
				}
				next = append(next, arr[step.index])
			}
		}
	}
	return next
}

func stringTraverse(current []interface{}, step step, next []interface{}) []interface{} {
	for _, n := range current {
		if m, ok := n.(map[string]interface{}); ok {
			if v, exists := m[step.key]; exists {
				if step.o != nil {
					m[step.key] = step.o(v)
					continue
				}
				next = append(next, v)
			}
		}
	}
	return next
}
