package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func Example() {
	var err error
	cmd := exec.Command("go", "run", "github.com/aereal/jsondiff/cmd/jsondiff", "-only", ".d", "./testdata/from.json", "./testdata/to.json")
	cmd.Dir, err = filepath.Abs("../..")
	if err != nil {
		panic(err)
	}
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Print(string(out))
	// Output:
	// --- from.json
	// +++ to.json
	// @@ -1,2 +1,2 @@
	// -4
	// +3
}
