// This package allows you to make comparisons between two generic values and override the default behavior of go comparison operators.
package comparer

import (
	"fmt"
	"reflect"
	"strings"
)

// A Config is the function that allows to apply a configuratión to Comparer.
// It allows creating new configurations in the future without modifying the New signature.
type Config func(*Comparer)

// A Comparator is the function that allows defining the behavior of the comparisons.
//
// Returns an integer comparing two values, and a boolean indicating if the two values are comparable. The result will be 0 if a == b, -1 if a < b, and +1 if a > b.
type Comparator func(path string, a interface{}, b interface{}) (int, bool)

// A Comparer holds the configurations of the comparison methods.
type Comparer struct {
	c Comparator
}

// CustomComparator returns a new Config that overrides the Comparator function.
func CustomComparator(c Comparator) Config {
	return func(comp *Comparer) {
		comp.c = c
	}
}

// New returns a new Comparer with the provided configuration.
func New(configs ...Config) *Comparer {
	c := Comparer{}
	for _, config := range configs {
		config(&c)
	}
	if c.c == nil {
		c.c = func(p string, l interface{}, r interface{}) (int, bool) {
			return 0, false
		}
	}
	return &c
}

// Compare returns an integer comparing two values, and a boolean indicating if the two values are comparable. The result will be 0 if a == b, -1 if a < b, and +1 if a > b.
func (c *Comparer) Compare(a interface{}, b interface{}) (int, bool) {
	return c.compare(reflect.ValueOf(a), reflect.ValueOf(b))
}

// Equal reports whether a and b are equal.
func (c *Comparer) Equal(a interface{}, b interface{}) bool {
	return c.equal("", reflect.ValueOf(a), reflect.ValueOf(b))
}

func (c *Comparer) compare(a reflect.Value, b reflect.Value) (int, bool) {
	if !a.IsValid() || !b.IsValid() {
		return 0, false
	} else if comparison, comparable := c.c("", c.value(a), c.value(b)); comparable {
		return comparison, comparable
	} else if a.Type() != b.Type() {
		return 0, false
	}

	switch a.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if a.Int() < b.Int() {
			return -1, true
		} else if a.Int() > b.Int() {
			return 1, true
		} else {
			return 0, true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if a.Uint() < b.Uint() {
			return -1, true
		} else if a.Uint() > b.Uint() {
			return 1, true
		} else {
			return 0, true
		}
	case reflect.Float32, reflect.Float64:
		if a.Float() < b.Float() {
			return -1, true
		} else if a.Float() > b.Float() {
			return 1, true
		} else {
			return 0, true
		}
	case reflect.String:
		return strings.Compare(a.String(), b.String()), true
	default:
		return 0, false
	}
}

func (c *Comparer) equal(path string, a reflect.Value, b reflect.Value) bool {
	if !a.IsValid() || !b.IsValid() {
		return a.IsValid() == b.IsValid()
	} else if comparison, comparable := c.c(path, c.value(a), c.value(b)); comparable {
		return comparison == 0
	} else if a.Type() != b.Type() {
		return false
	}

	switch a.Kind() {
	case reflect.Array:
		for i := 0; i < a.Len(); i++ {
			child := path + "[" + fmt.Sprintf("%d", i) + "]"
			if !c.equal(child, a.Index(i), b.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Interface:
		return c.equal(path, a.Elem(), b.Elem())
	case reflect.Map:
		if a.IsNil() != b.IsNil() {
			return false
		}
		if a.Len() != b.Len() {
			return false
		}
		for _, k := range a.MapKeys() {
			child := path + "[" + fmt.Sprintf("%v", k.Interface()) + "]"
			if !c.equal(child, a.MapIndex(k), b.MapIndex(k)) {
				return false
			}
		}
		return true
	case reflect.Ptr:
		return c.equal(path, a.Elem(), b.Elem())
	case reflect.Slice:
		if a.IsNil() != b.IsNil() {
			return false
		}
		if a.Len() != b.Len() {
			return false
		}
		for i := 0; i < a.Len(); i++ {
			child := path + "[" + fmt.Sprintf("%d", i) + "]"
			if !c.equal(child, a.Index(i), b.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Struct:
		for i := 0; i < a.Type().NumField(); i++ {
			child := path
			if child != "" {
				child += "."
			}
			child += a.Type().Field(i).Name
			if !c.equal(child, a.Field(i), b.Field(i)) {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(c.value(a), c.value(b))
	}
}

func (c *Comparer) value(v reflect.Value) interface{} {
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.Int:
		return int(v.Int())
	case reflect.Int8:
		return int8(v.Int())
	case reflect.Int16:
		return int16(v.Int())
	case reflect.Int32:
		return int32(v.Int())
	case reflect.Int64:
		return v.Int()
	case reflect.Uint:
		return uint(v.Uint())
	case reflect.Uint8:
		return uint8(v.Uint())
	case reflect.Uint16:
		return uint16(v.Uint())
	case reflect.Uint32:
		return uint32(v.Uint())
	case reflect.Uint64:
		return v.Uint()
	case reflect.Float32:
		return float32(v.Float())
	case reflect.Float64:
		return v.Float()
	case reflect.String:
		return v.String()
	default:
		//BUG(kupuka): Panic if try to compare nested structs with unexported fields.
		return v.Interface()
	}
}
