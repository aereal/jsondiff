package jsondiff

import (
	"github.com/itchyny/gojq"
)

func withUpdate(query *gojq.Query) *gojq.Query {
	var ret *gojq.Query
	qs := splitIntoTerms(query)
	for j := len(qs) - 1; j >= 0; j-- {
		if ret == nil { // most right leaf
			ret = &gojq.Query{
				Op:    gojq.OpAssign,
				Left:  qs[j],
				Right: nullRhs,
			}
			continue
		}

		ret = &gojq.Query{
			Op: gojq.OpPipe,
			Left: &gojq.Query{
				Op:    gojq.OpAssign,
				Left:  qs[j],
				Right: nullRhs,
			},
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
