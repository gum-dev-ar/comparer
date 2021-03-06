# Comparer

This package allows you to make comparisons between two generic values and override the default behavior of go comparison operators.

## Installation

```bash
$ go get github.com/gum-dev-ar/comparer
```

## Usage

The `c.Equal(x, y) bool` method reports whether x and y are equal. The default behavior is the same as the [reflect.DeepEqual](https://golang.org/pkg/reflect/#DeepEqual) function.

The `c.Compare(x, y) (int, bool)` method returns an integer comparing the x and y, and a boolean indicating if the two values are comparable. The result will be 0 if x == y, -1 if x < y, and +1 if x > y.

### Default comparer example

In this example, we use a default comparer to illustrate the use of the provided interfaces.

```golang
package main

import (
	"fmt"
	"github.com/gum-dev-ar/comparer"
)

func main() {
	c := comparer.Default()

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
```

This example will produce the following output.

```
2 == 2
2 != 4
2 < 5
```

### Custom comparer example

In this example, we will define a comparator where all the even numbers are equal to each other and, in the same way, all the odd numbers are equal to each other. On the other hand, the even numbers are greater than the odd ones.

```golang
package main

import (
	"fmt"

	"github.com/gum-dev-ar/comparer"
)

func number(v interface{}) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	default:
		return 0, false
	}
}

func comparator(_ string, a interface{}, b interface{}) (int, bool) {
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

func main() {
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
```

This example will produce the following output.

```
2 == 2
2 == 4
2 > 5
```
