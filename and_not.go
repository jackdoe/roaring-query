package roaringquery

import (
	"fmt"
	"strings"

	"github.com/RoaringBitmap/roaring"
)

type boolAndQuery struct {
	queries []Query
	not     Query
}

func (q *boolAndQuery) AddSubQuery(sub Query) {
	q.queries = append(q.queries, sub)
}

func AndNot(not Query, queries ...Query) *boolAndQuery {
	return And(queries...).SetNot(not)
}

func And(queries ...Query) *boolAndQuery {
	return &boolAndQuery{
		queries: queries,
	}
}

func (q *boolAndQuery) SetNot(not Query) *boolAndQuery {
	q.not = not
	return q
}

func (q *boolAndQuery) String() string {
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

func (q *boolAndQuery) Execute() *roaring.Bitmap {
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

func (q *boolAndQuery) Iterator() roaring.IntIterable {
	return q.Execute().Iterator()
}
