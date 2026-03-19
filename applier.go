package dotnotation

import (
	"sync"
)

// CompileApplier compiles the dotnotation expression and returns an Applier
func CompileApplier(expr string, op func(interface{}) interface{}) (*Applier, error) {
	e, err := compile(expr, op)
	return &Applier{e: e}, err
}

var applyPool = sync.Pool{
	New: func() any {
		p := [2][]interface{}{make([]interface{}, 0, 8), make([]interface{}, 0, 8)}
		return &p
	},
}

// Applier holds the compiled instructions to apply an operation to the compiled dotnotation path
type Applier struct {
	e *dotNotation
}

// Apply applies an operation to the compiled dotnotation path of the given data
func (e *Applier) Apply(data interface{}) {
	if e.e.wc == 0 {
		e.applyLinear(data)
		return
	}

	pair := applyPool.Get().(*[2][]interface{})
	current := append(pair[0], data) // skipcq: CRT-D0001
	next := pair[1]

	for _, step := range e.e.extractSteps {
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

	switch e.e.applyStep.ft {
	case stringType:
		stringApply(current, e.e.applyStep, e.e.op)
	case numericType:
		numericApply(current, e.e.applyStep, e.e.op)
	case wildcardType:
		wildcardApply(current, e.e.op)
	}

	pair[0] = current[:0]
	pair[1] = next
	applyPool.Put(pair)
}

func (e *Applier) applyLinear(data interface{}) { // skipcq: GO-R1005
	current := data

	for _, step := range e.e.extractSteps {
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
				continue
			}
			if arr, ok := current.([]interface{}); ok {
				if step.index >= len(arr) {
					return
				}
				current = arr[step.index]
			}
		}
	}

	switch e.e.applyStep.ft {
	case stringType:
		if m, ok := current.(map[string]interface{}); ok {
			if v, exists := m[e.e.applyStep.key]; exists {
				m[e.e.applyStep.key] = e.e.op(v)
			}
		}
	case numericType:
		if m, ok := current.(map[string]interface{}); ok {
			if v, exists := m[e.e.applyStep.key]; exists {
				m[e.e.applyStep.key] = e.e.op(v)
			}
			return
		}
		if arr, ok := current.([]interface{}); ok {
			if e.e.applyStep.index < len(arr) {
				arr[e.e.applyStep.index] = e.e.op(arr[e.e.applyStep.index])
			}
		}
	case wildcardType:
		switch v := current.(type) {
		case []interface{}:
			for i, vv := range v {
				v[i] = e.e.op(vv)
			}
		case map[string]interface{}:
			for k, vv := range v {
				v[k] = e.e.op(vv)
			}
		}
	}
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

func numericApply(current []interface{}, step step, op func(interface{}) interface{}) {
	for _, n := range current {
		if m, ok := n.(map[string]interface{}); ok {
			if v, exists := m[step.key]; exists {
				m[step.key] = op(v)
			}
			continue
		}
		if arr, ok := n.([]interface{}); ok {
			if step.index < len(arr) {
				arr[step.index] = op(arr[step.index])
			}
		}
	}
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
