package dotnotation

import (
	"sync"
)

// CompileExtractor compiles the dotnotation expression and returns an Extractor
func CompileExtractor(expr string) (*Extractor, error) {
	e, err := compile(expr, nil)
	return &Extractor{e: e}, err
}

// Extractor holds the compiled instructions to extract the compiled dotnotation path values
type Extractor struct {
	e *dotNotation
}

var extractPool = sync.Pool{
	New: func() any {
		s := make([]interface{}, 0, 8)
		return &s
	},
}

// Extract extracts the values in the compiled dotnotation path of the given data
func (e *Extractor) Extract(data interface{}) []interface{} {
	if e.e.wc == 0 {
		return e.extractLinear(data)
	}

	current := make([]interface{}, 1, 4*e.e.wc+1)
	current[0] = data

	nextPtr := extractPool.Get().(*[]interface{})
	next := *nextPtr

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

func (e *Extractor) extractLinear(data interface{}) []interface{} {
	current := data

	for _, step := range e.e.extractSteps {
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
