example:

	q := NewBoolAndNotQuery(
		NewBoolOrQuery(
			NewTerm("a", roaring.BitmapOf(1, 2)),
			NewTerm("b", roaring.BitmapOf(3, 9))),
		NewBoolAndQuery(
			NewTerm("c", roaring.BitmapOf(4, 5)),
			NewTerm("d", roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)),
		),
                NewTerm("e", roaring.BitmapOf(4, 5, 6)),
	)
        // ({c AND d} AND e) AND NOT {a OR b}

        iter := q.Iterator()
        for iter.HasNext() {
            match := iter.Next()
            log.Printf("document id: %d",match)
        }

