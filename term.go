package roaringquery

import (
	"github.com/RoaringBitmap/roaring"
)

type Term struct {
	r    *roaring.Bitmap
	term string
}

func NewTerm(t string, r *roaring.Bitmap) *Term {
	return &Term{
		term: t,
		r:    r,
	}
}

func (t *Term) String() string {
	return t.term
}

func (t *Term) Execute() *roaring.Bitmap {
	return t.r
}

func (t *Term) Iterator() roaring.IntIterable {
	return t.r.Iterator()
}
