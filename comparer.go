package comparer

import (
	"fmt"
	"reflect"
	"strings"
)

type Comparator func(string, interface{}, interface{}) (int, bool)

type Comparer struct {
	c Comparator
}

func Default() *Comparer {
	return &Comparer{
		c: func(p string, l interface{}, r interface{}) (int, bool) {
			return 0, false
		},
	}
}

func Custom(comparator Comparator) *Comparer {
	return &Comparer{c: comparator}
}

func (c *Comparer) Compare(left interface{}, right interface{}) (int, bool) {
	if comparison, comparable := c.c("", left, right); comparable {
		return comparison, comparable
	} else {
		return c.compare(reflect.ValueOf(left), reflect.ValueOf(right))
	}
}

func (c *Comparer) Equal(left interface{}, right interface{}) bool {
	return c.equal("", reflect.ValueOf(left), reflect.ValueOf(right))
}

func (c *Comparer) compare(left reflect.Value, right reflect.Value) (int, bool) {
	if left.Type() != right.Type() {
		return 0, false
	}

	switch left.Kind() {
	case reflect.Bool:
		if left.Bool() == right.Bool() {
			return 0, true
		} else if left.Bool() {
			return 1, true
		} else {
			return -1, true
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if left.Int() < right.Int() {
			return -1, true
		} else if left.Int() > right.Int() {
			return 1, true
		} else {
			return 0, true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if left.Uint() < right.Uint() {
			return -1, true
		} else if left.Uint() > right.Uint() {
			return 1, true
		} else {
			return 0, true
		}
	case reflect.Float32, reflect.Float64:
		if left.Float() < right.Float() {
			return -1, true
		} else if left.Float() > right.Float() {
			return 1, true
		} else {
			return 0, true
		}
	case reflect.String:
		return strings.Compare(left.String(), right.String()), true
	default:
		return 0, false
	}
}

func (c *Comparer) equal(path string, left reflect.Value, right reflect.Value) bool {
	if !left.IsValid() || !right.IsValid() {
		return left.IsValid() == right.IsValid()
	} else if comparison, comparable := c.c(path, c.value(left), c.value(right)); comparable {
		return comparison == 0
	} else if left.Type() != right.Type() {
		return false
	}

	switch left.Kind() {
	case reflect.Array:
		for i := 0; i < left.Len(); i++ {
			child := path + "[" + fmt.Sprintf("%d", i) + "]"
			if !c.equal(child, left.Index(i), right.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Map:
		if left.IsNil() != right.IsNil() {
			return false
		}
		if left.Len() != right.Len() {
			return false
		}
		for _, k := range left.MapKeys() {
			child := path + "[" + fmt.Sprintf("%v", k.Interface()) + "]"
			if !c.equal(child, left.MapIndex(k), right.MapIndex(k)) {
				return false
			}
		}
		return true
	case reflect.Slice:
		if left.IsNil() != right.IsNil() {
			return false
		}
		if left.Len() != right.Len() {
			return false
		}
		for i := 0; i < left.Len(); i++ {
			child := path + "[" + fmt.Sprintf("%d", i) + "]"
			if !c.equal(child, left.Index(i), right.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Struct:
		for i := 0; i < left.Type().NumField(); i++ {
			child := path
			if child != "" {
				child += "."
			}
			child += left.Type().Field(i).Name
			if !c.equal(child, left.Field(i), right.Field(i)) {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(c.value(left), c.value(right))
	}
}

func (c *Comparer) value(v reflect.Value) interface{} {
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint()
	case reflect.Float32, reflect.Float64:
		return v.Float()
	case reflect.String:
		return v.String()
	default:
		return v.Interface()
	}
}
