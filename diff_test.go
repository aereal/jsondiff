package jsondiff

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/itchyny/gojq"
)

func TestDiff(t *testing.T) {
	testCases := []struct {
		name         string
		opts         []Option
		wantDiffPath string
	}{
		{
			"nothing",
			[]Option{},
			"./testdata/nothing.diff",
		},
		{
			"ignore",
			[]Option{Ignore(parseQuery(t, ".b, .c"))},
			"./testdata/ignore.diff",
		},
		{
			"only",
			[]Option{Only(parseQuery(t, ".d"))},
			"./testdata/only.diff",
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
			want, err := os.ReadFile(tc.wantDiffPath)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(string(want), got); diff != "" {
				t.Errorf("-want,+got:\n%s", diff)
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
