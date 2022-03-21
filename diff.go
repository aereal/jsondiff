package jsondiff

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/itchyny/gojq"
)

var ErrEitherOnlyOneOption = errors.New("either of only one of Ignore() or Only() must be specified")

type opt struct {
	ignore *gojq.Query
	only   *gojq.Query
}

func (o *opt) validate() error {
	if o.ignore != nil && o.only != nil {
		return ErrEitherOnlyOneOption
	}
	return nil
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

// DiffFromFiles calculates differences with flies' contents.
func DiffFromFiles(from, to fs.File, opts ...Option) (string, error) {
	l, err := newInputFromFile(from)
	if err != nil {
		return "", fmt.Errorf("left: %w", err)
	}
	r, err := newInputFromFile(to)
	if err != nil {
		return "", fmt.Errorf("right: %w", err)
	}
	return takeDiff(l, r, opts...)
}

// DiffFromObjects calculates differences with from and to.
func DiffFromObjects(from, to interface{}, opts ...Option) (string, error) {
	return takeDiff(&input{x: from, name: "from"}, &input{x: to, name: "to"}, opts...)
}

// Diff calculates differences with from and to.
//
// Deprecated: use DiffObjects()
func Diff(from, to interface{}, opts ...Option) (string, error) {
	return DiffFromObjects(from, to, opts...)
}

func newInputFromFile(f fs.File) (*input, error) {
	var i input
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	i.name = st.Name()
	if err := json.NewDecoder(f).Decode(&i.x); err != nil {
		return nil, err
	}
	return &i, nil
}

type input struct {
	name string
	x    interface{}
}

func takeDiff(from, to *input, opts ...Option) (string, error) {
	o := &opt{}
	for _, f := range opts {
		f(o)
	}
	if err := o.validate(); err != nil {
		return "", err
	}
	var (
		fromObj = from.x
		toObj   = to.x
	)
	switch {
	case o.ignore != nil:
		var err error
		q := withUpdate(o.ignore)
		fromObj, err = modifyValue(q, fromObj)
		if err != nil {
			return "", fmt.Errorf("modify(lhs): %v", err)
		}
		toObj, err = modifyValue(q, toObj)
		if err != nil {
			return "", fmt.Errorf("modify(rhs): %v", err)
		}
	case o.only != nil:
		var err error
		fromObj, err = modifyValue(o.only, fromObj)
		if err != nil {
			return "", fmt.Errorf("modify(lhs): %v", err)
		}
		toObj, err = modifyValue(o.only, toObj)
		if err != nil {
			return "", fmt.Errorf("modify(lhs): %v", err)
		}
	}
	l, err := toJSON(fromObj)
	if err != nil {
		return "", fmt.Errorf("toJSON(lhs): %v", err)
	}
	r, err := toJSON(toObj)
	if err != nil {
		return "", fmt.Errorf("toJSON(rhs): %v", err)
	}
	edits := myers.ComputeEdits(span.URIFromPath(""), l, r)
	d := gotextdiff.ToUnified(from.name, to.name, l, edits)
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
