Simple query dsl on top of roaring bitmaps

Example:

    package main

    import (
    	"log"
  	rq "github.com/jackdoe/roaring-query"
    	"github.com/RoaringBitmap/roaring"
    )

    func main() {
    	q := rq.AndNot(
    		rq.Or(
    			rq.Term("a", roaring.BitmapOf(1, 2)),
    			rq.Term("b", roaring.BitmapOf(3, 9))),
    		rq.And(
    			rq.Term("c", roaring.BitmapOf(4, 5)),
    			rq.Term("d", roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)),
    		),
    		rq.Term("e", roaring.BitmapOf(4, 5, 6)),
    	)
    	// ({c AND d} AND e) AND NOT {a OR b}

    	iter := q.Iterator()
    	for iter.HasNext() {
    		match := iter.Next()
    		log.Printf("document id: %d", match)
    	}
    }

func And(queries ...Query) *boolAndQuery
func AndNot(not Query, queries ...Query) *boolAndQuery
func Or(queries ...Query) *boolOrQuery
func Term(t string, r *roaring.Bitmap) *termQuery
type Query interface{ ... }
