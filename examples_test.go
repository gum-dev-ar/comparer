package comparer_test

import (
	"fmt"
	"github.com/gum-dev-ar/comparer"
)

func Example_default() {
	c := comparer.New()

	x, y, z := 2, 4, 5

	if c.Equal(x, x) {
		fmt.Printf("%v == %v\n", x, x)
	} else {
		fmt.Printf("%v != %v\n", x, x)
	}

	if c.Equal(x, y) {
		fmt.Printf("%v == %v\n", x, y)
	} else {
		fmt.Printf("%v != %v\n", x, y)
	}

	comparison, comparable := c.Compare(x, z)
	if !comparable {
		fmt.Printf("%v and %v are not comparable\n", x, z)
	} else if comparison < 0 {
		fmt.Printf("%v < %v\n", x, z)
	} else if comparison > 0 {
		fmt.Printf("%v > %v\n", x, z)
	} else {
		fmt.Printf("%v == %v\n", x, z)
	}
}

func Example_custom() {
	number := func(v interface{}) (int, bool) {
		switch n := v.(type) {
		case int:
			return n, true
		default:
			return 0, false
		}
	}

	comparator := func(_ string, a interface{}, b interface{}) (int, bool) {
		na, ok := number(a)
		if !ok {
			return 0, false
		}

		nb, ok := number(b)
		if !ok {
			return 0, false
		}

		if (na%2) == 0 && (nb%2) == 0 {
			return 0, true
		} else if (na%2) == 0 && (nb%2) != 0 {
			return 1, true
		} else if (na%2) != 0 && (nb%2) == 0 {
			return -1, true
		} else {
			return 0, true
		}
	}

	config := comparer.CustomComparator(comparator)
	c := comparer.New(config)

	x, y, z := 2, 4, 5

	if c.Equal(x, x) {
		fmt.Printf("%v == %v\n", x, x)
	} else {
		fmt.Printf("%v != %v\n", x, x)
	}

	if c.Equal(x, y) {
		fmt.Printf("%v == %v\n", x, y)
	} else {
		fmt.Printf("%v != %v\n", x, y)
	}

	comparison, comparable := c.Compare(x, z)
	if !comparable {
		fmt.Printf("%v and %v are not comparable\n", x, z)
	} else if comparison < 0 {
		fmt.Printf("%v < %v\n", x, z)
	} else if comparison > 0 {
		fmt.Printf("%v > %v\n", x, z)
	} else {
		fmt.Printf("%v == %v\n", x, z)
	}
}
