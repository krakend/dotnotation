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
