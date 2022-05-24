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
		wantErr      bool
	}{
		{
			"nothing",
			[]Option{},
			"./testdata/nothing.diff",
			false,
		},
		{
			"both only and ignore",
			[]Option{Ignore(parseQuery(t, ".b, .c")), Only(parseQuery(t, ".d"))},
			"",
			true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lhs := map[string]interface{}{"a": 1, "b": 2, "c": 3, "d": 4}
			rhs := map[string]interface{}{"a": 1, "b": 1, "c": 2, "d": 3}
			got, err := DiffFromObjects(lhs, rhs, tc.opts...)
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v got=%v (%#v)", tc.wantErr, err, err)
			}
			if tc.wantErr {
				return
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

func Test_removing(t *testing.T) {
	queries := []struct {
		query string
		want  string
	}{
		{".age", "del(.age)"},
		{".age, .name", "del(.age, .name)"},
		{".age, .name, .meta", "del(.age, .name, .meta)"},
		{".meta[]", "del(.meta[])"},
		{".meta[0:-1]", "del(.meta[0:-1])"},
	}
	for _, c := range queries {
		t.Run(c.query, func(t *testing.T) {
			want := parseQuery(t, c.want)
			got := removing(parseQuery(t, c.query))
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
