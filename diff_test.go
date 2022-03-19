package jsondiff

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/itchyny/gojq"
)

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
