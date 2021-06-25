package comparer_test

import (
	"fmt"
	"testing"

	"github.com/gum-dev-ar/comparer"
)

type es1 struct {
	a int
	b string
}

type es2 struct {
	a int
	b string
}

type tc struct {
	a interface{}
	b interface{}
}

var groups = map[string][]tc{
	"Array": {
		{nil, [0]int{}},
		{nil, [2]int{1, 2}},
		{[2]int{1, 2}, [2]int{3, 4}},
		{[2]int{1, 2}, [2]int{2, 1}},
		{[2]int{1}, [2]int{1, 2}},
	},
	"Bool": {
		{nil, true},
		{nil, false},
		{true, false},
	},
	"Float32": {
		{nil, float32(0.2)},
		{float32(0.2), float32(1.7)},
		{float32(-3.2), float32(3.2)},
	},
	"Float64": {
		{nil, float64(0.2)},
		{float64(0.2), float64(1.7)},
		{float64(-3.2), float64(3.2)},
	},
	"Int": {
		{nil, int(1)},
		{int(1), int(2)},
		{int(-1), int(1)},
	},
	"Int8": {
		{nil, int8(1)},
		{int8(1), int8(2)},
		{int8(-1), int8(1)},
	},
	"Int16": {
		{nil, int16(1)},
		{int16(1), int16(2)},
		{int16(-1), int16(1)},
	},
	"Int32": {
		{nil, int32(1)},
		{int32(1), int32(2)},
		{int32(-1), int32(1)},
	},
	"Int64": {
		{nil, int64(1)},
		{int64(1), int64(2)},
		{int64(-1), int64(1)},
	},
	"Map": {
		{nil, map[string]int(nil)},
		{nil, map[string]int{}},
		{nil, map[string]int{"A": 1, "B": 2}},
		{map[string]int(nil), map[string]int{}}, //CHECK
		{map[string]int(nil), map[string]int{"A": 1, "B": 2}},
		{map[string]int{}, map[string]int{"A": 1, "B": 2}},
		{map[string]int{"A": 1, "B": 2}, map[string]int{"C": 1, "D": 2}},
		{map[string]int{"A": 1, "B": 2}, map[string]int{"A": 3, "B": 4}},
		{map[string]int{"A": 1, "B": 2}, map[string]int{"A": 1}},
	},
	"Slice": {
		{nil, []int(nil)},
		{nil, []int{}},
		{nil, []int{1, 2}},
		{[]int(nil), []int{}},
		{[]int(nil), []int{1, 2}},
		{[]int{}, []int{1, 2}},
		{[]int{1, 2}, []int{3, 4}},
		{[]int{1, 2}, []int{2, 1}},
		{[]int{1}, []int{1, 2}},
	},
	"String": {
		{nil, "test1"},
		{nil, ""},
		{"", "test1"},
		{"test1", "test2"},
		{"test1", "TEST1"},
	},
	"Struct": {
		{nil, es1{1, "A"}},
		{es1{1, "A"}, es2{1, "A"}},
		{es1{1, "A"}, es1{2, "B"}},
	},
	"Uint": {
		{nil, uint(1)},
		{uint(1), uint(2)},
	},
	"Uint8": {
		{nil, uint8(1)},
		{uint8(1), uint8(2)},
	},
	"Uint16": {
		{nil, uint16(1)},
		{uint16(1), uint16(2)},
	},
	"Uint32": {
		{nil, uint32(1)},
		{uint32(1), uint32(2)},
	},
	"Uint64": {
		{nil, uint64(1)},
		{uint64(1), uint64(2)},
	},
}

func TestEqual(t *testing.T) {
	c := comparer.Default()

	for name, cases := range groups {
		t.Run(name, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(fmt.Sprintf("comparer.Equal(%+v, %+v)", tc.a, tc.a), func(t *testing.T) {
					if !c.Equal(tc.a, tc.a) {
						t.Errorf("The values should be equal")
					}
				})
				t.Run(fmt.Sprintf("comparer.Equal(%+v, %+v)", tc.a, tc.b), func(t *testing.T) {
					if c.Equal(tc.a, tc.b) {
						t.Errorf("The values should not be equal")
					}
				})
				t.Run(fmt.Sprintf("comparer.Equal(%+v, %+v)", tc.b, tc.a), func(t *testing.T) {
					if c.Equal(tc.b, tc.a) {
						t.Errorf("The values should not be equal")
					}
				})
			}
		})
	}
}
