package roaringquery

import (
	"strings"

	"github.com/RoaringBitmap/roaring"
)

type boolOrQuery struct {
	queries []Query
}

func Or(queries ...Query) *boolOrQuery {
	return &boolOrQuery{
		queries: queries,
	}
}

func (q *boolOrQuery) AddSubQuery(sub Query) {
	q.queries = append(q.queries, sub)
}

func (q *boolOrQuery) String() string {
	out := []string{}
	for _, v := range q.queries {
		out = append(out, v.String())
	}
	return "{" + strings.Join(out, " OR ") + "}"
}

func (q *boolOrQuery) Execute() *roaring.Bitmap {
	m := []*roaring.Bitmap{}
	for _, s := range q.queries {
		m = append(m, s.Execute())
	}
	out := roaring.FastOr(m...)
	return out
}

func (q *boolOrQuery) Iterator() roaring.IntIterable {
	return q.Execute().Iterator()
}
