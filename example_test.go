package jsondiff

import (
	"fmt"

	"github.com/itchyny/gojq"
)

var (
	lhs = map[string]interface{}{"a": 1, "b": 2, "c": 3, "d": 4}
	rhs = map[string]interface{}{"a": 1, "b": 1, "c": 2, "d": 3}
)

func ExampleDiffFromObjects_only() {
	query, err := gojq.Parse(".d")
	if err != nil {
		panic(err)
	}
	diff, err := DiffFromObjects(lhs, rhs, Only(query))
	if err != nil {
		panic(err)
	}
	fmt.Println(diff)
	// Output:
	// --- from
	// +++ to
	// @@ -1,2 +1,2 @@
	// -4
	// +3
}

func ExampleDiffFromObjects_ignore() {
	query, err := gojq.Parse(".b, .c")
	if err != nil {
		panic(err)
	}
	diff, err := DiffFromObjects(lhs, rhs, Ignore(query))
	if err != nil {
		panic(err)
	}
	fmt.Println(diff)
	// Output:
	// --- from
	// +++ to
	// @@ -2,6 +2,6 @@
	//    "a": 1,
	//    "b": null,
	//    "c": null,
	// -  "d": 4
	// +  "d": 3
	//  }
}
