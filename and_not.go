package roaringquery

import (
	"fmt"
	"strings"

	"github.com/RoaringBitmap/roaring"
)

type BoolAndQuery struct {
	queries []Query
	not     Query
}

func (q *BoolAndQuery) AddSubQuery(sub Query) {
	q.queries = append(q.queries, sub)
}

func NewBoolAndNotQuery(not Query, queries ...Query) *BoolAndQuery {
	return NewBoolAndQuery(queries...).SetNot(not)
}

func NewBoolAndQuery(queries ...Query) *BoolAndQuery {
	return &BoolAndQuery{
		queries: queries,
	}
}

func (q *BoolAndQuery) SetNot(not Query) *BoolAndQuery {
	q.not = not
	return q
}

func (q *BoolAndQuery) String() string {
	out := []string{}
	for _, v := range q.queries {
		out = append(out, v.String())
	}
	s := strings.Join(out, " AND ")
	if q.not != nil {
		s = fmt.Sprintf("%s[-(%s)]", s, q.not.String())
	}
	return "{" + s + "}"
}

func (q *BoolAndQuery) Execute() *roaring.Bitmap {
	m := []*roaring.Bitmap{}
	for _, s := range q.queries {
		m = append(m, s.Execute())
	}

	out := roaring.FastAnd(m...)

	if q.not != nil {
		out.AndNot(q.not.Execute())
	}
	return out
}

func (q *BoolAndQuery) Iterator() roaring.IntIterable {
	return q.Execute().Iterator()
}
