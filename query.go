package roaringquery

import (
	"github.com/RoaringBitmap/roaring"
)

type Query interface {
	Execute() *roaring.Bitmap
	Iterator() roaring.IntIterable
	String() string
}
