package jsondiff

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/itchyny/gojq"
)

type opt struct {
	ignore *gojq.Query
	only   *gojq.Query
}

// Option is a function to modify Diff's behavior.
type Option func(*opt)

// Ignore returns Option function that indicates the function to ignore structures pointed by the query.
func Ignore(query *gojq.Query) Option {
	return func(o *opt) {
		o.ignore = query
	}
}

// Only returns Option function that indicates the function to calculate differences based on the structure pointed by the query.
func Only(query *gojq.Query) Option {
	return func(o *opt) {
		o.only = query
	}
}

// Diff calculates differences with lhs and rhs.
func Diff(lhs, rhs interface{}, opts ...Option) (string, error) {
	o := &opt{}
	for _, f := range opts {
		f(o)
	}
	var (
		left  = lhs
		right = rhs
	)
	switch {
	case o.ignore != nil:
		var err error
		q := withUpdate(o.ignore)
		left, err = modifyValue(q, lhs)
		if err != nil {
			return "", fmt.Errorf("modify(lhs): %v", err)
		}
		right, err = modifyValue(q, rhs)
		if err != nil {
			return "", fmt.Errorf("modify(rhs): %v", err)
		}
	case o.only != nil:
		var err error
		left, err = modifyValue(o.only, lhs)
		if err != nil {
			return "", fmt.Errorf("modify(lhs): %v", err)
		}
		right, err = modifyValue(o.only, right)
		if err != nil {
			return "", fmt.Errorf("modify(lhs): %v", err)
		}
	}
	l, err := toJSON(left)
	if err != nil {
		return "", fmt.Errorf("toJSON(lhs): %v", err)
	}
	r, err := toJSON(right)
	if err != nil {
		return "", fmt.Errorf("toJSON(rhs): %v", err)
	}
	edits := myers.ComputeEdits(span.URIFromPath(""), l, r)
	d := gotextdiff.ToUnified("lhs", "rhs", l, edits)
	return fmt.Sprint(d), nil
}

func modifyValue(query *gojq.Query, x interface{}) (interface{}, error) {
	iter := query.Run(x)
	var ret interface{}
	for {
		got, hasNext := iter.Next()
		if !hasNext {
			break
		}
		if err, ok := got.(error); ok {
			return nil, err
		}
		ret = got
	}
	return ret, nil
}

func toJSON(x interface{}) (string, error) {
	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.SetIndent("", "  ")
	if err := enc.Encode(x); err != nil {
		return "", err
	}
	b.WriteRune('\n')
	return b.String(), nil
}

func withUpdate(query *gojq.Query) *gojq.Query {
	var ret *gojq.Query
	qs := splitIntoTerms(query)
	for j := len(qs) - 1; j >= 0; j-- {
		left := &gojq.Query{
			Op:    gojq.OpAssign,
			Left:  qs[j],
			Right: nullRhs,
		}
		if ret == nil { // most right leaf
			ret = left
			continue
		}

		ret = &gojq.Query{
			Op:    gojq.OpPipe,
			Left:  left,
			Right: ret,
		}
	}
	return ret
}

var nullRhs *gojq.Query

func init() {
	var err error
	nullRhs, err = gojq.Parse("null")
	if err != nil {
		panic(err)
	}
}

func splitIntoTerms(q *gojq.Query) []*gojq.Query {
	ret := []*gojq.Query{}
	if q.Term != nil {
		ret = append(ret, q)
		return ret
	}
	if q.Left != nil {
		ret = append(ret, splitIntoTerms(q.Left)...)
	}
	if q.Right != nil {
		ret = append(ret, splitIntoTerms(q.Right)...)
	}
	return ret
}
