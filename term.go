package roaringquery

import (
	"github.com/RoaringBitmap/roaring"
)

type termQuery struct {
	r    *roaring.Bitmap
	term string
}

func Term(t string, r *roaring.Bitmap) *termQuery {
	return &termQuery{
		term: t,
		r:    r,
	}
}

func (t *termQuery) String() string {
	return t.term
}

func (t *termQuery) Execute() *roaring.Bitmap {
	return t.r
}

func (t *termQuery) Iterator() roaring.IntIterable {
	return t.r.Iterator()
}
