package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test(t *testing.T) {
	type testCase struct {
		name        string
		argv        []string
		wantStatus  int
		wantOutPath string
		wantErr     string
	}
	testCases := []testCase{
		{"ok; no queries", []string{"jsondiff", "../../testdata/from.json", "../../testdata/to.json"}, codeOK, "./testdata/ok_no_queries.stdout.txt", ""},
		{"ok; with -exit-code", []string{"jsondiff", "-exit-code", "../../testdata/from.json", "../../testdata/to.json"}, codeHaveDifferences, "./testdata/ok_with_exit_code.stdout.txt", ""},
		{"ok; with -only option", []string{"jsondiff", "-only", ".d", "../../testdata/from.json", "../../testdata/to.json"}, codeOK, "./testdata/ok_with_only_option.stdout.txt", ""},
		{"ok; with -ignore option", []string{"jsondiff", "-ignore", ".d", "../../testdata/from.json", "../../testdata/to.json"}, codeOK, "./testdata/ok_with_ignore_option.stdout.txt", ""},
		{"no paths given", []string{"jsondiff"}, codeError, "", "2 file paths must be passed\n"},
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
			var wantOut string
			if tc.wantOutPath != "" {
				out, err := os.ReadFile(tc.wantOutPath)
				if err != nil {
					t.Fatal(err)
				}
				wantOut = string(out)
			}
			gotOut := outStream.String()
			if diff := cmp.Diff(wantOut, gotOut); diff != "" {
				t.Errorf("output (-want, +got):\n%s", diff)
			}
			gotErr := errStream.String()
			if diff := cmp.Diff(tc.wantErr, gotErr); diff != "" {
				t.Errorf("error (-want, +got):\n%s", diff)
			}
		})
	}
}
