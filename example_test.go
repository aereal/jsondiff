package jsondiff

import (
	"fmt"

	"github.com/itchyny/gojq"
)

var (
	lhs = map[string]interface{}{"a": 1, "b": 2, "c": 3, "d": 4}
	rhs = map[string]interface{}{"a": 1, "b": 1, "c": 2, "d": 3}
)

func ExampleDiff_only() {
	query, err := gojq.Parse(".d")
	if err != nil {
		panic(err)
	}
	diff, err := Diff(lhs, rhs, Only(query))
	if err != nil {
		panic(err)
	}
	fmt.Println(diff)
	// Output:
	// --- lhs
	// +++ rhs
	// @@ -1,2 +1,2 @@
	// -4
	// +3
}

func ExampleDiff_ignore() {
	query, err := gojq.Parse(".b, .c")
	if err != nil {
		panic(err)
	}
	diff, err := Diff(lhs, rhs, Ignore(query))
	if err != nil {
		panic(err)
	}
	fmt.Println(diff)
	// Output:
	// --- lhs
	// +++ rhs
	// @@ -2,6 +2,6 @@
	//    "a": 1,
	//    "b": null,
	//    "c": null,
	// -  "d": 4
	// +  "d": 3
	//  }
}
