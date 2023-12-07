package models

import "reflect"

type ComparatorFunc func(actual interface{}) bool

type Expectation struct {
	Actual     interface{}
	Comparator ComparatorFunc
	SuccessMsg string
	FailMsg    string
}

type ExpectationBuilder struct {
	control *Control
}

// func (eb *ExpectationBuilder) To(comparator ComparatorFunc, successMsg string, failMsg string) *Control {
// 	eb.control.Expectation = &Expectation{
// 		Comparator: comparator,
// 		SuccessMsg: successMsg,
// 		FailMsg:    failMsg,
// 	}
// 	return eb.control
// }

// NotEqual comparator
func NotEqual(expected interface{}) ComparatorFunc {
	return func(actual interface{}) bool {
		return actual != expected
	}
}

// Equal Comparator
func Equal(expected interface{}) ComparatorFunc {
	return func(actual interface{}) bool {
		return actual == expected
	}
}

// GreaterThan Comparator
func GreaterThan(expected interface{}) ComparatorFunc {
	return func(actual interface{}) bool {
		return actual.(int) > expected.(int)
	}
}

// GreaterThanOrEqual Comparator
func GreaterThanOrEqual(expected interface{}) ComparatorFunc {
	return func(actual interface{}) bool {
		return actual.(int) >= expected.(int)
	}
}

// LessThan Comparator
func LessThan(expected interface{}) ComparatorFunc {
	return func(actual interface{}) bool {
		return actual.(int) < expected.(int)
	}
}

// LessThanOrEqual Comparator
func LessThanOrEqual(expected interface{}) ComparatorFunc {
	return func(actual interface{}) bool {
		return actual.(int) <= expected.(int)
	}
}

// Contains checks if a slice contains a specific value.
func Contains(expected interface{}) ComparatorFunc {
	return func(actual interface{}) bool {
		// Type assertion to check if actual is a slice
		actualSlice, ok := toSlice(actual)
		if !ok {
			// Actual is not a slice, so it cannot contain the expected value
			return false
		}

		for _, value := range actualSlice {
			if reflect.DeepEqual(value, expected) {
				return true
			}
		}
		return false
	}
}

// NotContains checks if a slice does not contain a specific value.
func NotContains(expected interface{}) ComparatorFunc {
	return func(actual interface{}) bool {
		// Type assertion to check if actual is a slice
		actualSlice, ok := toSlice(actual)
		if !ok {
			// Actual is not a slice, so it cannot contain the expected value
			return true
		}

		for _, value := range actualSlice {
			if reflect.DeepEqual(value, expected) {
				return false
			}
		}
		return true
	}
}

// toSlice attempts to convert an interface{} to a []interface{}.
func toSlice(slice interface{}) ([]interface{}, bool) {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		return nil, false
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, true
}
