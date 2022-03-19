package jsondiff

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/itchyny/gojq"
)

func TestDiff(t *testing.T) {
	testCases := []struct {
		name     string
		opts     []Option
		wantDiff string
	}{
		{
			"nothing",
			[]Option{},
			"--- lhs\n+++ rhs\n@@ -1,7 +1,7 @@\n {\n   \"a\": 1,\n-  \"b\": 2,\n-  \"c\": 3,\n-  \"d\": 4\n+  \"b\": 1,\n+  \"c\": 2,\n+  \"d\": 3\n }\n \n",
		},
		{
			"ignore",
			[]Option{Ignore(parseQuery(t, ".b, .c"))},
			"--- lhs\n+++ rhs\n@@ -2,6 +2,6 @@\n   \"a\": 1,\n   \"b\": null,\n   \"c\": null,\n-  \"d\": 4\n+  \"d\": 3\n }\n \n",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lhs := map[string]interface{}{"a": 1, "b": 2, "c": 3, "d": 4}
			rhs := map[string]interface{}{"a": 1, "b": 1, "c": 2, "d": 3}
			got, err := Diff(lhs, rhs, tc.opts...)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.wantDiff {
				t.Log(got)
				t.Errorf("got:\n%q\nwant:\n%q", got, tc.wantDiff)
			}
		})
	}
}

func Test_withUpdate(t *testing.T) {
	queries := []struct {
		query string
		want  string
	}{
		{".age", ".age = null"},
		{".age, .name", ".age = null | .name = null"},
		{".age, .name, .meta", ".age = null | .name = null | .meta = null"},
		{".meta[]", ".meta[] = null"},
		{".meta[0:-1]", ".meta[0:-1] = null"},
	}
	for _, c := range queries {
		t.Run(c.query, func(t *testing.T) {
			want := parseQuery(t, c.want)
			got := withUpdate(parseQuery(t, c.query))
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("-want, +got:\n%s", diff)
			}
		})
	}
}

func parseQuery(t *testing.T, q string) *gojq.Query {
	t.Helper()
	parsed, err := gojq.Parse(q)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return parsed
}
