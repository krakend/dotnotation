package dotnotation

import (
	"errors"
	"strconv"
	"strings"
	"sync"
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
}

// Applier holds the compiled instructions to apply an operation to the compiled dotnotation path
type Applier struct {
	e *dotNotation
}

// Apply applies an operation to the compiled dotnotation path of the given data
func (a *Applier) Apply(v interface{}) {
	a.e.Apply(v)
}

// Extractor holds the compiled instructions to extract the compiled dotnotation path values
type Extractor struct {
	e *dotNotation
}

// Extract extracts the values in the compiled dotnotation path of the given data
func (e *Extractor) Extract(v interface{}) []interface{} {
	return e.e.Extract(v)
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

	// if there's a wildcard on last step, we remove it from counter so it goes through Linear paths anyway
	// since it does not need to allocate a slice while traversing the main path
	if steps[len(steps)-1].ft == wildcardType {
		wc--
	}

	var applyStep step
	if op != nil {
		applyStep = steps[len(steps)-1]
		steps = steps[:len(steps)-1]
	}

	return &dotNotation{
		extractSteps: steps,
		applyStep:    applyStep,
		op:           op,
		wc:           wc,
	}, nil
}

type dotNotation struct {
	extractSteps []step
	applyStep    step
	op           func(interface{}) interface{}
	wc           int
}

var extractPool = sync.Pool{
	New: func() any {
		s := make([]interface{}, 0, 8)
		return &s
	},
}

var applyPool = sync.Pool{
	New: func() any {
		p := [2][]interface{}{make([]interface{}, 0, 8), make([]interface{}, 0, 8)}
		return &p
	},
}

func (e *dotNotation) Extract(data interface{}) []interface{} {
	if e.wc == 0 {
		return e.extractLinear(data)
	}

	current := make([]interface{}, 1, 4*e.wc+1)
	current[0] = data

	nextPtr := extractPool.Get().(*[]interface{})
	next := *nextPtr

	for _, step := range e.extractSteps {
		switch step.ft {
		case stringType:
			next = stringTraverse(current, step, next)
		case numericType:
			next = numericTraverse(current, step, next)
		case wildcardType:
			next = wildcardTraverse(current, next)
		}

		if len(next) == 0 {
			*nextPtr = next
			extractPool.Put(nextPtr)
			return next[:0:0]
		}

		current, next = next, current[:0]
	}

	*nextPtr = next
	extractPool.Put(nextPtr)
	return current
}

func (e *dotNotation) extractLinear(data interface{}) []interface{} {
	current := data

	for _, step := range e.extractSteps {
		switch step.ft {
		case stringType:
			m, ok := current.(map[string]interface{})
			if !ok {
				return nil
			}
			v, exists := m[step.key]
			if !exists {
				return nil
			}
			current = v
		case numericType:
			if m, ok := current.(map[string]interface{}); ok {
				v, exists := m[step.key]
				if !exists {
					return nil
				}
				current = v
			} else if arr, ok := current.([]interface{}); ok {
				if step.index >= len(arr) {
					return nil
				}
				current = arr[step.index]
			}
		// wildcard can only happen on last step
		case wildcardType:
			if m, ok := current.(map[string]interface{}); ok {
				res := make([]interface{}, 0, len(m))
				for _, v := range m {
					res = append(res, v)
				}
				return res
			} else if arr, ok := current.([]interface{}); ok {
				res := make([]interface{}, 0, len(arr))
				return append(res, arr...)
			}
		}
	}
	return []interface{}{current}
}

func (e *dotNotation) Apply(data interface{}) {
	if e.wc == 0 {
		e.applyLinear(data)
		return
	}

	pair := applyPool.Get().(*[2][]interface{})
	current := append(pair[0], data)
	next := pair[1]

	for _, step := range e.extractSteps {
		switch step.ft {
		case stringType:
			next = stringTraverse(current, step, next)
		case numericType:
			next = numericTraverse(current, step, next)
		case wildcardType:
			next = wildcardTraverse(current, next)
		}

		if len(next) == 0 {
			pair[0] = current[:0]
			pair[1] = next
			applyPool.Put(pair)
			return
		}

		current, next = next, current[:0]
	}

	switch e.applyStep.ft {
	case stringType:
		stringApply(current, e.applyStep, e.op)
	case numericType:
		numericApply(current, e.applyStep, e.op)
	case wildcardType:
		wildcardApply(current, e.op)
	}

	pair[0] = current[:0]
	pair[1] = next
	applyPool.Put(pair)
}

func (e *dotNotation) applyLinear(data interface{}) {
	current := data

	for _, step := range e.extractSteps {
		switch step.ft {
		case stringType:
			m, ok := current.(map[string]interface{})
			if !ok {
				return
			}
			v, exists := m[step.key]
			if !exists {
				return
			}
			current = v
		case numericType:
			if m, ok := current.(map[string]interface{}); ok {
				v, exists := m[step.key]
				if !exists {
					return
				}
				current = v
			} else if arr, ok := current.([]interface{}); ok {
				if step.index >= len(arr) {
					return
				}
				current = arr[step.index]
			}
		}
	}

	switch e.applyStep.ft {
	case stringType:
		if m, ok := current.(map[string]interface{}); ok {
			if v, exists := m[e.applyStep.key]; exists {
				m[e.applyStep.key] = e.op(v)
			}
		}
	case numericType:
		if m, ok := current.(map[string]interface{}); ok {
			if v, exists := m[e.applyStep.key]; exists {
				m[e.applyStep.key] = e.op(v)
			}
		} else if arr, ok := current.([]interface{}); ok {
			if e.applyStep.index < len(arr) {
				arr[e.applyStep.index] = e.op(arr[e.applyStep.index])
			}
		}
	case wildcardType:
		switch v := current.(type) {
		case []interface{}:
			for i, vv := range v {
				v[i] = e.op(vv)
			}
		case map[string]interface{}:
			for k, vv := range v {
				v[k] = e.op(vv)
			}
		}
	}
}

func wildcardTraverse(current []interface{}, next []interface{}) []interface{} {
	for _, n := range current {
		switch v := n.(type) {
		case []interface{}:
			next = append(next, v...)
		case map[string]interface{}:
			for _, vv := range v {
				next = append(next, vv)
			}
		}
	}
	return next
}

func wildcardApply(current []interface{}, op func(interface{}) interface{}) {
	for _, n := range current {
		switch v := n.(type) {
		case []interface{}:
			for i, vv := range v {
				v[i] = op(vv)
			}
		case map[string]interface{}:
			for i, vv := range v {
				v[i] = op(vv)
			}
		}
	}
}

func numericTraverse(current []interface{}, step step, next []interface{}) []interface{} {
	for _, n := range current {
		if m, ok := n.(map[string]interface{}); ok {
			if v, exists := m[step.key]; exists {
				next = append(next, v)
			}
			continue
		}
		if arr, ok := n.([]interface{}); ok {
			if step.index < len(arr) {
				next = append(next, arr[step.index])
			}
		}
	}
	return next
}

func numericApply(current []interface{}, step step, op func(interface{}) interface{}) {
	for _, n := range current {
		if m, ok := n.(map[string]interface{}); ok {
			if v, exists := m[step.key]; exists {
				m[step.key] = op(v)
			}
		} else if arr, ok := n.([]interface{}); ok {
			if step.index < len(arr) {
				arr[step.index] = op(arr[step.index])
			}
		}
	}
}

func stringTraverse(current []interface{}, step step, next []interface{}) []interface{} {
	for _, n := range current {
		if m, ok := n.(map[string]interface{}); ok {
			if v, exists := m[step.key]; exists {
				next = append(next, v)
			}
		}
	}
	return next
}

func stringApply(current []interface{}, step step, op func(interface{}) interface{}) {
	for _, n := range current {
		if m, ok := n.(map[string]interface{}); ok {
			if v, exists := m[step.key]; exists {
				m[step.key] = op(v)
			}
		}
	}
}
