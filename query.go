// Simple query dsl on top of roaring bitmaps
//
// Example:
//
//  package main
//
//  import (
//  	"log"
//  	rq "github.com/jackdoe/roaring-query"
//  	"github.com/RoaringBitmap/roaring"
//  )
//
//  func main() {
//  	q := rq.AndNot(
//  		rq.Or(
//  			rq.Term("a", roaring.BitmapOf(1, 2)),
//  			rq.Term("b", roaring.BitmapOf(3, 9))),
//  		rq.And(
//  			rq.Term("c", roaring.BitmapOf(4, 5)),
//  			rq.Term("d", roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)),
//  		),
//  		rq.Term("e", roaring.BitmapOf(4, 5, 6)),
//  	)
//  	// ({c AND d} AND e) AND NOT {a OR b}
//
//  	iter := q.Iterator()
//  	for iter.HasNext() {
//  		match := iter.Next()
//  		log.Printf("document id: %d", match)
//  	}
//  }
package roaringquery

import (
	"github.com/RoaringBitmap/roaring"
)

type Query interface {
	Execute() *roaring.Bitmap
	Iterator() roaring.IntIterable
	String() string
}
