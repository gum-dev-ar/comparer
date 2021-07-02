package comparer_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gum-dev-ar/comparer"
)

type es1 struct {
	A int
	B string
}

type es2 struct {
	A int
	B string
}

type es3 struct {
	A es1
	B *es2
}

type dtc struct {
	min        interface{}
	max        interface{}
	comparable bool
}

type etc struct {
	a          interface{}
	b          interface{}
	comparable bool
}

var cdifferent = map[string][]dtc{
	"Array": {
		{nil, [0]int{}, false},
		{nil, [2]int{1, 2}, false},
		{[2]int{1, 2}, [2]int{3, 4}, false},
		{[2]int{1, 2}, [2]int{2, 1}, false},
		{[2]int{1}, [2]int{1, 2}, false},
		{[2]es1{{1, "test1"}, {2, "test2"}}, [2]es1{{1, "test1"}, {3, "test3"}}, false},
		{[2]es1{{1, "test1"}, {2, "test2"}}, [2]es2{{1, "test1"}, {2, "test2"}}, false},
		{[2]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, [2]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{5, "test5"}}}, false},
	},
	"Bool": {
		{nil, true, false},
		{nil, false, false},
		{true, false, false},
	},
	"Float32": {
		{nil, float32(0.2), false},
		{float32(0.2), float32(1.7), true},
		{float32(-3.2), float32(3.2), true},
	},
	"Float64": {
		{nil, float64(0.2), false},
		{float64(0.2), float64(1.7), true},
		{float64(-3.2), float64(3.2), true},
	},
	"Int": {
		{nil, int(1), false},
		{int(1), int(2), true},
		{int(-1), int(1), true},
	},
	"Int8": {
		{nil, int8(1), false},
		{int8(1), int8(2), true},
		{int8(-1), int8(1), true},
	},
	"Int16": {
		{nil, int16(1), false},
		{int16(1), int16(2), true},
		{int16(-1), int16(1), true},
	},
	"Int32": {
		{nil, int32(1), false},
		{int32(1), int32(2), true},
		{int32(-1), int32(1), true},
	},
	"Int64": {
		{nil, int64(1), false},
		{int64(1), int64(2), true},
		{int64(-1), int64(1), true},
	},
	"Map": {
		{nil, map[string]int(nil), false},
		{nil, map[string]int{}, false},
		{nil, map[string]int{"A": 1, "B": 2}, false},
		{map[string]int(nil), map[string]int{}, false},
		{map[string]int(nil), map[string]int{"A": 1, "B": 2}, false},
		{map[string]int{}, map[string]int{"A": 1, "B": 2}, false},
		{map[string]int{"A": 1, "B": 2}, map[string]int{"C": 1, "D": 2}, false},
		{map[string]int{"A": 1, "B": 2}, map[string]int{"A": 3, "B": 4}, false},
		{map[string]int{"A": 1, "B": 2}, map[string]int{"A": 1}, false},
	},
	"Slice": {
		{nil, []int(nil), false},
		{nil, []int{}, false},
		{nil, []int{1, 2}, false},
		{[]int(nil), []int{}, false},
		{[]int(nil), []int{1, 2}, false},
		{[]int{}, []int{1, 2}, false},
		{[]int{1, 2}, []int{3, 4}, false},
		{[]int{1, 2}, []int{2, 1}, false},
		{[]int{1}, []int{1, 2}, false},
		{[]es1{{1, "test1"}, {2, "test2"}}, []es1{{1, "test1"}, {3, "test3"}}, false},
		{[]es1{{1, "test1"}, {2, "test2"}}, []es2{{1, "test1"}, {2, "test2"}}, false},
		{[]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, []es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{5, "test5"}}}, false},
	},
	"String": {
		{nil, "test1", false},
		{nil, "", false},
		{"", "test1", true},
		{"test1", "test2", true},
	},
	"Struct": {
		{nil, es1{1, "A"}, false},
		{es1{}, es1{1, "A"}, false},
		{es1{1, "test1"}, es2{1, "test1"}, false},
		{es1{1, "test1"}, es1{1, "test2"}, false},
	},
	"Uint": {
		{nil, uint(1), false},
		{uint(1), uint(2), true},
	},
	"Uint8": {
		{nil, uint8(1), false},
		{uint8(1), uint8(2), true},
	},
	"Uint16": {
		{nil, uint16(1), false},
		{uint16(1), uint16(2), true},
	},
	"Uint32": {
		{nil, uint32(1), false},
		{uint32(1), uint32(2), true},
	},
	"Uint64": {
		{nil, uint64(1), false},
		{uint64(1), uint64(2), true},
	},
}

var cequal = map[string][]etc{
	"Array": {
		{[2]int{1, 2}, [2]int{1, 2}, false},
		{[2]es1{{1, "test1"}, {2, "test2"}}, [2]es1{{1, "test1"}, {2, "test2"}}, false},
		{[2]es1{{1, "test1"}, {2, "test2"}}, [2]es1{{1, "Test1"}, {2, "Test2"}}, false},
		{[2]es1{{1, "test1"}, {2, "test2"}}, [2]es1{{1, "TEST1"}, {2, "TEST2"}}, false},
		{[2]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, [2]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, false},
		{[2]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, [2]es3{{es1{1, "Test1"}, &es2{2, "Test2"}}, {es1{3, "Test3"}, &es2{4, "Test4"}}}, false},
		{[2]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, [2]es3{{es1{1, "TEST1"}, &es2{2, "TEST2"}}, {es1{3, "TEST3"}, &es2{4, "TEST4"}}}, false},
	},
	"Bool": {
		{true, true, false},
		{false, false, false},
	},
	"Float32": {
		{float32(0.2), float32(0.2), true},
		{float32(-3.2), float32(-3.2), true},
	},
	"Float64": {
		{float64(0.2), float64(0.2), true},
		{float64(-3.2), float64(-3.2), true},
	},
	"Int": {
		{int(1), int(1), true},
		{int(-1), int(-1), true},
	},
	"Int8": {
		{int8(1), int8(1), true},
		{int8(-1), int8(-1), true},
	},
	"Int16": {
		{int16(1), int16(1), true},
		{int16(-1), int16(-1), true},
	},
	"Int32": {
		{int32(1), int32(1), true},
		{int32(-1), int32(-1), true},
	},
	"Int64": {
		{int64(1), int64(1), true},
		{int64(-1), int64(-1), true},
	},
	"Map": {
		{map[string]int{"A": 1, "B": 2}, map[string]int{"A": 1, "B": 2}, false},
	},
	"Slice": {
		{[]int{1, 2}, []int{1, 2}, false},
		{[]es1{{1, "test1"}, {2, "test2"}}, []es1{{1, "test1"}, {2, "test2"}}, false},
		{[]es1{{1, "test1"}, {2, "test2"}}, []es1{{1, "Test1"}, {2, "Test2"}}, false},
		{[]es1{{1, "test1"}, {2, "test2"}}, []es1{{1, "TEST1"}, {2, "TEST2"}}, false},
		{[]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, []es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, false},
		{[]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, []es3{{es1{1, "Test1"}, &es2{2, "Test2"}}, {es1{3, "Test3"}, &es2{4, "Test4"}}}, false},
		{[]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, []es3{{es1{1, "TEST1"}, &es2{2, "TEST2"}}, {es1{3, "TEST3"}, &es2{4, "TEST4"}}}, false},
	},
	"String": {
		{"test1", "test1", true},
		{"Test1", "Test1", true},
		{"TEST1", "TEST1", true},
		{"Test1", "test1", true},
		{"TEST1", "test1", true},
	},
	"Struct": {
		{es1{1, "test1"}, es1{1, "test1"}, false},
		{es1{1, "Test1"}, es1{1, "Test1"}, false},
		{es1{1, "TEST1"}, es1{1, "TEST1"}, false},
		{es1{1, "Test1"}, es1{1, "test1"}, false},
		{es1{1, "TEST1"}, es1{1, "test1"}, false},
	},
	"Uint": {
		{uint(1), uint(1), true},
	},
	"Uint8": {
		{uint8(1), uint8(1), true},
	},
	"Uint16": {
		{uint16(1), uint16(1), true},
	},
	"Uint32": {
		{uint32(1), uint32(1), true},
	},
	"Uint64": {
		{uint64(1), uint64(1), true},
	},
}

var ddifferent = map[string][]dtc{
	"Array": {
		{nil, [0]int{}, false},
		{nil, [2]int{1, 2}, false},
		{[2]int{1, 2}, [2]int{3, 4}, false},
		{[2]int{1, 2}, [2]int{2, 1}, false},
		{[2]int{1}, [2]int{1, 2}, false},
		{[2]es1{{1, "test1"}, {2, "test2"}}, [2]es1{{1, "test1"}, {3, "test3"}}, false},
		{[2]es1{{1, "test1"}, {2, "test2"}}, [2]es1{{1, "Test1"}, {2, "Test2"}}, false},
		{[2]es1{{1, "test1"}, {2, "test2"}}, [2]es2{{1, "test1"}, {2, "test2"}}, false},
		{[2]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, [2]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{5, "test5"}}}, false},
		{[2]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, [2]es3{{es1{1, "Test1"}, &es2{2, "Test2"}}, {es1{3, "Test3"}, &es2{4, "Test4"}}}, false},
		{[2]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, [2]es3{{es1{1, "TEST1"}, &es2{2, "TEST2"}}, {es1{3, "TEST3"}, &es2{4, "TEST4"}}}, false},
	},
	"Bool": {
		{nil, true, false},
		{nil, false, false},
		{true, false, false},
	},
	"Float32": {
		{nil, float32(0.2), false},
		{float32(0.2), float32(1.7), true},
		{float32(-3.2), float32(3.2), true},
	},
	"Float64": {
		{nil, float64(0.2), false},
		{float64(0.2), float64(1.7), true},
		{float64(-3.2), float64(3.2), true},
	},
	"Int": {
		{nil, int(1), false},
		{int(1), int(2), true},
		{int(-1), int(1), true},
	},
	"Int8": {
		{nil, int8(1), false},
		{int8(1), int8(2), true},
		{int8(-1), int8(1), true},
	},
	"Int16": {
		{nil, int16(1), false},
		{int16(1), int16(2), true},
		{int16(-1), int16(1), true},
	},
	"Int32": {
		{nil, int32(1), false},
		{int32(1), int32(2), true},
		{int32(-1), int32(1), true},
	},
	"Int64": {
		{nil, int64(1), false},
		{int64(1), int64(2), true},
		{int64(-1), int64(1), true},
	},
	"Map": {
		{nil, map[string]int(nil), false},
		{nil, map[string]int{}, false},
		{nil, map[string]int{"A": 1, "B": 2}, false},
		{map[string]int(nil), map[string]int{}, false},
		{map[string]int(nil), map[string]int{"A": 1, "B": 2}, false},
		{map[string]int{}, map[string]int{"A": 1, "B": 2}, false},
		{map[string]int{"A": 1, "B": 2}, map[string]int{"C": 1, "D": 2}, false},
		{map[string]int{"A": 1, "B": 2}, map[string]int{"A": 3, "B": 4}, false},
		{map[string]int{"A": 1, "B": 2}, map[string]int{"A": 1}, false},
	},
	"Slice": {
		{nil, []int(nil), false},
		{nil, []int{}, false},
		{nil, []int{1, 2}, false},
		{[]int(nil), []int{}, false},
		{[]int(nil), []int{1, 2}, false},
		{[]int{}, []int{1, 2}, false},
		{[]int{1, 2}, []int{3, 4}, false},
		{[]int{1, 2}, []int{2, 1}, false},
		{[]int{1}, []int{1, 2}, false},
		{[]es1{{1, "test1"}, {2, "test2"}}, []es1{{1, "test1"}, {3, "test3"}}, false},
		{[]es1{{1, "test1"}, {2, "test2"}}, []es1{{1, "Test1"}, {2, "Test2"}}, false},
		{[]es1{{1, "test1"}, {2, "test2"}}, []es2{{1, "test1"}, {2, "test2"}}, false},
		{[]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, []es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{5, "test5"}}}, false},
		{[]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, []es3{{es1{1, "Test1"}, &es2{2, "Test2"}}, {es1{3, "Test3"}, &es2{4, "Test4"}}}, false},
		{[]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, []es3{{es1{1, "TEST1"}, &es2{2, "TEST2"}}, {es1{3, "TEST3"}, &es2{4, "TEST4"}}}, false},
	},
	"String": {
		{nil, "test1", false},
		{nil, "", false},
		{"", "test1", true},
		{"test1", "test2", true},
		{"Test1", "test1", true},
		{"TEST1", "test1", true},
	},
	"Struct": {
		{nil, es1{1, "test1"}, false},
		{es1{}, es1{1, "test1"}, false},
		{es1{1, "test1"}, es2{1, "test1"}, false},
		{es1{1, "test1"}, es1{1, "test2"}, false},
		{es3{es1{1, "test1"}, &es2{2, "test2"}}, es3{es1{1, "test3"}, &es2{2, "test2"}}, false},
		{es3{es1{1, "test1"}, &es2{2, "test2"}}, es3{es1{1, "test1"}, &es2{2, "test3"}}, false},
		{es3{es1{}, &es2{2, "test2"}}, es3{es1{1, "test1"}, &es2{2, "test2"}}, false},
		{es3{es1{1, "test1"}, nil}, es3{es1{1, "test1"}, &es2{2, "test2"}}, false},
	},
	"Uint": {
		{nil, uint(1), false},
		{uint(1), uint(2), true},
	},
	"Uint8": {
		{nil, uint8(1), false},
		{uint8(1), uint8(2), true},
	},
	"Uint16": {
		{nil, uint16(1), false},
		{uint16(1), uint16(2), true},
	},
	"Uint32": {
		{nil, uint32(1), false},
		{uint32(1), uint32(2), true},
	},
	"Uint64": {
		{nil, uint64(1), false},
		{uint64(1), uint64(2), true},
	},
}

var dequal = map[string][]etc{
	"Array": {
		{[2]int{1, 2}, [2]int{1, 2}, false},
		{[2]es1{{1, "test1"}, {2, "test2"}}, [2]es1{{1, "test1"}, {2, "test2"}}, false},
		{[2]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, [2]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, false},
	},
	"Bool": {
		{true, true, false},
		{false, false, false},
	},
	"Float32": {
		{float32(0.2), float32(0.2), true},
		{float32(-3.2), float32(-3.2), true},
	},
	"Float64": {
		{float64(0.2), float64(0.2), true},
		{float64(-3.2), float64(-3.2), true},
	},
	"Int": {
		{int(1), int(1), true},
		{int(-1), int(-1), true},
	},
	"Int8": {
		{int8(1), int8(1), true},
		{int8(-1), int8(-1), true},
	},
	"Int16": {
		{int16(1), int16(1), true},
		{int16(-1), int16(-1), true},
	},
	"Int32": {
		{int32(1), int32(1), true},
		{int32(-1), int32(-1), true},
	},
	"Int64": {
		{int64(1), int64(1), true},
		{int64(-1), int64(-1), true},
	},
	"Map": {
		{map[string]int{"A": 1, "B": 2}, map[string]int{"A": 1, "B": 2}, false},
	},
	"Slice": {
		{[]int{1, 2}, []int{1, 2}, false},
		{[]es1{{1, "test1"}, {2, "test2"}}, []es1{{1, "test1"}, {2, "test2"}}, false},
		{[]es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, []es3{{es1{1, "test1"}, &es2{2, "test2"}}, {es1{3, "test3"}, &es2{4, "test4"}}}, false},
	},
	"String": {
		{"test1", "test1", true},
		{"Test1", "Test1", true},
		{"TEST1", "TEST1", true},
	},
	"Struct": {
		{es1{}, es1{}, false},
		{es1{1, "test1"}, es1{1, "test1"}, false},
		{es3{es1{1, "test1"}, &es2{2, "test2"}}, es3{es1{1, "test1"}, &es2{2, "test2"}}, false},
		{es3{es1{}, &es2{2, "test2"}}, es3{es1{}, &es2{2, "test2"}}, false},
		{es3{es1{1, "test1"}, nil}, es3{es1{1, "test1"}, nil}, false},
	},
	"Uint": {
		{uint(1), uint(1), true},
	},
	"Uint8": {
		{uint8(1), uint8(1), true},
	},
	"Uint16": {
		{uint16(1), uint16(1), true},
	},
	"Uint32": {
		{uint32(1), uint32(1), true},
	},
	"Uint64": {
		{uint64(1), uint64(1), true},
	},
}

func TestCustomCompare(t *testing.T) {
	convert := func(v interface{}) (string, bool) {
		switch s := v.(type) {
		case string:
			return s, true
		default:
			return "", false
		}
	}
	comparator := func(_ string, a interface{}, b interface{}) (int, bool) {
		sa, ok := convert(a)
		if !ok {
			return 0, false
		}
		sa = strings.ToUpper(sa)

		sb, ok := convert(b)
		if !ok {
			return 0, false
		}
		sb = strings.ToUpper(sb)

		return strings.Compare(sa, sb), true
	}

	config := comparer.CustomComparator(comparator)
	c := comparer.New(config)

	run := func(t *testing.T, a interface{}, b interface{}, expected int, isComparable bool) {
		comparison, comparable := c.Compare(a, b)
		if comparable != isComparable {
			if isComparable {
				t.Errorf("The values should be comparable")
			} else {
				t.Errorf("The values should not be comparable")
			}
		} else if comparable && comparison != expected {
			if expected < 0 {
				t.Errorf("The value %+v shoud be greater than the value %+v", b, a)
			} else if expected > 0 {
				t.Errorf("The value %+v shoud be greater than the value %+v", a, b)
			} else {
				t.Errorf("The values should be equal")
			}
		}
	}

	for name, cases := range cdifferent {
		t.Run(name, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(fmt.Sprintf("comparer.Compare(%+v,%+v)", tc.min, tc.max), func(t *testing.T) {
					run(t, tc.min, tc.max, -1, tc.comparable)
				})
				t.Run(fmt.Sprintf("comparer.Compare(%+v,%+v)", tc.max, tc.min), func(t *testing.T) {
					run(t, tc.max, tc.min, 1, tc.comparable)
				})
				t.Run(fmt.Sprintf("comparer.Compare(&%+v,%+v)", &tc.min, tc.max), func(t *testing.T) {
					run(t, &tc.min, tc.max, 0, false)
				})
				t.Run(fmt.Sprintf("comparer.Compare(%+v,&%+v)", tc.min, &tc.max), func(t *testing.T) {
					run(t, tc.min, &tc.max, 0, false)
				})
				t.Run(fmt.Sprintf("comparer.Compare(&%+v,&%+v)", &tc.min, &tc.max), func(t *testing.T) {
					run(t, &tc.max, &tc.min, 0, false)
				})
			}
		})
	}

	for name, cases := range cequal {
		t.Run(name, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(fmt.Sprintf("comparer.Compare(%+v,%+v)", tc.a, tc.b), func(t *testing.T) {
					run(t, tc.a, tc.b, 0, tc.comparable)
				})
				t.Run(fmt.Sprintf("comparer.Compare(%+v,%+v)", tc.b, tc.a), func(t *testing.T) {
					run(t, tc.b, tc.a, 0, tc.comparable)
				})
			}
		})
	}
}

func TestCustomEqual(t *testing.T) {
	convert := func(v interface{}) (string, bool) {
		switch s := v.(type) {
		case string:
			return s, true
		default:
			return "", false
		}
	}
	comparator := func(_ string, a interface{}, b interface{}) (int, bool) {
		sa, ok := convert(a)
		if !ok {
			return 0, false
		}
		sa = strings.ToUpper(sa)

		sb, ok := convert(b)
		if !ok {
			return 0, false
		}
		sb = strings.ToUpper(sb)

		return strings.Compare(sa, sb), true
	}

	config := comparer.CustomComparator(comparator)
	c := comparer.New(config)

	run := func(t *testing.T, a interface{}, b interface{}, equal bool) {
		if c.Equal(a, b) != equal {
			if equal {
				t.Errorf("The values should be equal")
			} else {
				t.Errorf("The values should not be equal")
			}
		}
	}

	for name, cases := range cdifferent {
		t.Run(name, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(fmt.Sprintf("comparer.Equal(%+v,%+v)", tc.min, tc.max), func(t *testing.T) {
					run(t, tc.min, tc.max, false)
				})
				t.Run(fmt.Sprintf("comparer.Equal(%+v,%+v)", tc.max, tc.min), func(t *testing.T) {
					run(t, tc.max, tc.min, false)
				})
				t.Run(fmt.Sprintf("comparer.Equal(&%+v,%+v)", tc.min, tc.max), func(t *testing.T) {
					run(t, &tc.min, tc.max, false)
				})
				t.Run(fmt.Sprintf("comparer.Equal(%+v,&%+v)", tc.max, tc.min), func(t *testing.T) {
					run(t, tc.max, &tc.min, false)
				})
				t.Run(fmt.Sprintf("comparer.Equal(&%+v,&%+v)", tc.min, tc.max), func(t *testing.T) {
					run(t, &tc.min, &tc.max, false)
				})
				t.Run(fmt.Sprintf("comparer.Equal(&%+v,&%+v)", tc.max, tc.min), func(t *testing.T) {
					run(t, &tc.max, &tc.min, false)
				})
			}
		})
	}

	for name, cases := range cequal {
		t.Run(name, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(fmt.Sprintf("comparer.Equal(%+v,%+v)", tc.a, tc.b), func(t *testing.T) {
					run(t, tc.a, tc.b, true)
				})
				t.Run(fmt.Sprintf("comparer.Equal(%+v,%+v)", tc.b, tc.a), func(t *testing.T) {
					run(t, tc.b, tc.a, true)
				})
				t.Run(fmt.Sprintf("comparer.Equal(&%+v,%+v)", tc.a, tc.b), func(t *testing.T) {
					run(t, &tc.a, tc.b, false)
				})
				t.Run(fmt.Sprintf("comparer.Equal(%+v,&%+v)", tc.b, tc.a), func(t *testing.T) {
					run(t, tc.b, &tc.a, false)
				})
				t.Run(fmt.Sprintf("comparer.Equal(&%+v,&%+v)", tc.a, tc.b), func(t *testing.T) {
					run(t, &tc.a, &tc.b, true)
				})
				t.Run(fmt.Sprintf("comparer.Equal(&%+v,&%+v)", tc.b, tc.a), func(t *testing.T) {
					run(t, &tc.b, &tc.a, true)
				})
			}
		})
	}
}

func TestDefaultCompare(t *testing.T) {
	c := comparer.New()

	run := func(t *testing.T, a interface{}, b interface{}, expected int, isComparable bool) {
		comparison, comparable := c.Compare(a, b)
		if comparable != isComparable {
			if isComparable {
				t.Errorf("The values should be comparable")
			} else {
				t.Errorf("The values should not be comparable")
			}
		} else if comparable && comparison != expected {
			if expected < 0 {
				t.Errorf("The value %+v shoud be greater than the value %+v", b, a)
			} else if expected > 0 {
				t.Errorf("The value %+v shoud be greater than the value %+v", a, b)
			} else {
				t.Errorf("The values should be equal")
			}
		}
	}

	for name, cases := range ddifferent {
		t.Run(name, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(fmt.Sprintf("comparer.Compare(%+v,%+v)", tc.min, tc.max), func(t *testing.T) {
					run(t, tc.min, tc.max, -1, tc.comparable)
				})
				t.Run(fmt.Sprintf("comparer.Compare(%+v,%+v)", tc.max, tc.min), func(t *testing.T) {
					run(t, tc.max, tc.min, 1, tc.comparable)
				})
				t.Run(fmt.Sprintf("comparer.Compare(&%+v,%+v)", &tc.min, tc.max), func(t *testing.T) {
					run(t, &tc.min, tc.max, 0, false)
				})
				t.Run(fmt.Sprintf("comparer.Compare(%+v,&%+v)", tc.min, &tc.max), func(t *testing.T) {
					run(t, tc.min, &tc.max, 0, false)
				})
				t.Run(fmt.Sprintf("comparer.Compare(&%+v,&%+v)", &tc.min, &tc.max), func(t *testing.T) {
					run(t, &tc.max, &tc.min, 0, false)
				})
			}
		})
	}

	for name, cases := range dequal {
		t.Run(name, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(fmt.Sprintf("comparer.Compare(%+v,%+v)", tc.a, tc.b), func(t *testing.T) {
					run(t, tc.a, tc.b, 0, tc.comparable)
				})
				t.Run(fmt.Sprintf("comparer.Compare(%+v,%+v)", tc.b, tc.a), func(t *testing.T) {
					run(t, tc.b, tc.a, 0, tc.comparable)
				})
			}
		})
	}
}

func TestDefaultEqual(t *testing.T) {
	c := comparer.New()

	run := func(t *testing.T, a interface{}, b interface{}) {
		e := reflect.DeepEqual(a, b)
		if c.Equal(a, b) != e {
			if e {
				t.Errorf("The values should be equal")
			} else {
				t.Errorf("The values should not be equal")
			}
		}
	}

	for name, cases := range ddifferent {
		t.Run(name, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(fmt.Sprintf("comparer.Equal(%+v,%+v)", tc.min, tc.max), func(t *testing.T) {
					run(t, tc.min, tc.max)
				})
				t.Run(fmt.Sprintf("comparer.Equal(%+v,%+v)", tc.max, tc.min), func(t *testing.T) {
					run(t, tc.max, tc.min)
				})
				t.Run(fmt.Sprintf("comparer.Equal(&%+v,%+v)", tc.min, tc.max), func(t *testing.T) {
					run(t, &tc.min, tc.max)
				})
				t.Run(fmt.Sprintf("comparer.Equal(%+v,&%+v)", tc.max, tc.min), func(t *testing.T) {
					run(t, tc.max, &tc.min)
				})
				t.Run(fmt.Sprintf("comparer.Equal(&%+v,&%+v)", tc.min, tc.max), func(t *testing.T) {
					run(t, &tc.min, &tc.max)
				})
				t.Run(fmt.Sprintf("comparer.Equal(&%+v,&%+v)", tc.max, tc.min), func(t *testing.T) {
					run(t, &tc.max, &tc.min)
				})
			}
		})
	}

	for name, cases := range dequal {
		t.Run(name, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(fmt.Sprintf("comparer.Equal(%+v,%+v)", tc.a, tc.b), func(t *testing.T) {
					run(t, tc.a, tc.b)
				})
				t.Run(fmt.Sprintf("comparer.Equal(%+v,%+v)", tc.b, tc.a), func(t *testing.T) {
					run(t, tc.b, tc.a)
				})
				t.Run(fmt.Sprintf("comparer.Equal(&%+v,%+v)", tc.a, tc.b), func(t *testing.T) {
					run(t, &tc.a, tc.b)
				})
				t.Run(fmt.Sprintf("comparer.Equal(%+v,&%+v)", tc.b, tc.a), func(t *testing.T) {
					run(t, tc.b, &tc.a)
				})
				t.Run(fmt.Sprintf("comparer.Equal(&%+v,&%+v)", tc.a, tc.b), func(t *testing.T) {
					run(t, &tc.a, &tc.b)
				})
				t.Run(fmt.Sprintf("comparer.Equal(&%+v,&%+v)", tc.b, tc.a), func(t *testing.T) {
					run(t, &tc.b, &tc.a)
				})
			}
		})
	}
}
