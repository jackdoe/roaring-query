package roaringquery

import (
	"strings"

	"github.com/RoaringBitmap/roaring"
)

type BoolOrQuery struct {
	queries []Query
}

func (q *BoolOrQuery) AddSubQuery(sub Query) {
	q.queries = append(q.queries, sub)
}

func NewBoolOrQuery(queries ...Query) *BoolOrQuery {
	return &BoolOrQuery{
		queries: queries,
	}
}

func (q *BoolOrQuery) String() string {
	out := []string{}
	for _, v := range q.queries {
		out = append(out, v.String())
	}
	return "{" + strings.Join(out, " OR ") + "}"
}

func (q *BoolOrQuery) Execute() *roaring.Bitmap {
	m := []*roaring.Bitmap{}
	for _, s := range q.queries {
		m = append(m, s.Execute())
	}
	out := roaring.FastOr(m...)
	return out
}

func (q *BoolOrQuery) Iterator() roaring.IntIterable {
	return q.Execute().Iterator()
}
