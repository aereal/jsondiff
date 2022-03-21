package main

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test(t *testing.T) {
	type testCase struct {
		name       string
		argv       []string
		wantStatus int
		wantOut    string
		wantErr    string
	}
	testCases := []testCase{
		{"ok; no queries", []string{"jsondiff", "../../testdata/from.json", "../../testdata/to.json"}, 0, "--- from.json\n+++ to.json\n@@ -1,7 +1,7 @@\n {\n   \"a\": 1,\n-  \"b\": 2,\n-  \"c\": 3,\n-  \"d\": 4\n+  \"b\": 1,\n+  \"c\": 2,\n+  \"d\": 3\n }\n \n", ""},
		{"ok; with -only option", []string{"jsondiff", "-only", ".d", "../../testdata/from.json", "../../testdata/to.json"}, 0, "--- from.json\n+++ to.json\n@@ -1,2 +1,2 @@\n-4\n+3\n \n", ""},
		{"ok; with -ignore option", []string{"jsondiff", "-ignore", ".d", "../../testdata/from.json", "../../testdata/to.json"}, 0, "--- from.json\n+++ to.json\n@@ -1,7 +1,7 @@\n {\n   \"a\": 1,\n-  \"b\": 2,\n-  \"c\": 3,\n+  \"b\": 1,\n+  \"c\": 2,\n   \"d\": null\n }\n \n", ""},
		{"no paths given", []string{"jsondiff"}, 1, "", "2 file paths must be passed\n"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				outStream bytes.Buffer
				errStream bytes.Buffer
			)
			a := &app{outstream: &outStream, errStream: &errStream}
			gotStatus := a.run(tc.argv)
			if gotStatus != tc.wantStatus {
				t.Errorf("status code: got=%d want=%d", tc.wantStatus, gotStatus)
			}
			if diff := cmp.Diff(tc.wantOut, outStream.String()); diff != "" {
				t.Errorf("output:\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantErr, errStream.String()); diff != "" {
				t.Errorf("error:\n%s", diff)
			}
		})
	}
}
