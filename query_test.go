package roaringquery

import (
	"testing"

	"github.com/RoaringBitmap/roaring"
)

func postingsList(n int) *roaring.Bitmap {
	list := roaring.New()
	for i := 0; i < n; i++ {
		list.Add(uint32(i) * 3)
	}
	return list
}

func query(query Query) *roaring.Bitmap {
	return query.Execute()
}

func eq(t *testing.T, a, b *roaring.Bitmap) {
	if !a.Equals(b) {
		t.Logf("a != b: %s %s", a.String(), b.String())
		panic("no")
	}
}

func BenchmarkAnd1000(b *testing.B) {
	x := postingsList(1000000)
	y := postingsList(1000)

	for n := 0; n < b.N; n++ {
		sum := uint32(0)
		q := NewBoolAndQuery(
			NewTerm("x", x),
			NewTerm("y", y),
		)
		i := q.Iterator()
		for i.HasNext() {
			sum += i.Next()
		}
	}
}

func BenchmarkOr1000(b *testing.B) {
	x := postingsList(1000000)
	y := postingsList(1000)

	for n := 0; n < b.N; n++ {
		sum := uint32(0)
		q := NewBoolOrQuery(
			NewTerm("x", x),
			NewTerm("y", y),
		)
		i := q.Iterator()
		for i.HasNext() {
			sum += i.Next()
		}
	}
}

func TestModify(t *testing.T) {
	a := postingsList(100)
	b := postingsList(1000)
	c := postingsList(10000)
	d := postingsList(100000)
	e := postingsList(1000000)

	eq(t, a, query(NewTerm("x", a)))
	eq(t, b, query(NewTerm("x", b)))
	eq(t, c, query(NewTerm("x", c)))
	eq(t, d, query(NewTerm("x", d)))
	eq(t, e, query(NewTerm("x", e)))

	eq(t, b, query(NewBoolOrQuery(
		NewTerm("a", a),
		NewTerm("b", b),
	)))

	eq(t, c, query(NewBoolOrQuery(
		NewTerm("a", a),
		NewTerm("b", b),
		NewTerm("c", c),
	)))

	eq(t, e, query(NewBoolOrQuery(
		NewTerm("a", a),
		NewTerm("b", b),
		NewTerm("c", c),
		NewTerm("d", d),
		NewTerm("e", e),
	)))

	eq(t, a, query(NewBoolAndQuery(
		NewTerm("a", a),
		NewTerm("b", b),
		NewTerm("c", c),
		NewTerm("d", d),
		NewTerm("e", e),
	)))

	eq(t, roaring.BitmapOf(4, 6, 7, 8, 10), query(NewBoolAndNotQuery(
		NewTerm("x", roaring.BitmapOf(1, 2, 3, 9)),
		NewBoolOrQuery(
			NewTerm("x", roaring.BitmapOf(3, 4)),
			NewTerm("x", roaring.BitmapOf(1, 2, 3, 6, 7, 8, 9, 10)),
		),
	)))
	eq(t, roaring.BitmapOf(6, 7, 8, 10), query(NewBoolAndNotQuery(
		NewTerm("x", roaring.BitmapOf(1, 2, 3, 9)),
		NewBoolAndNotQuery(
			NewTerm("x", roaring.BitmapOf(4, 5)),
			NewTerm("x", roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)),
			NewTerm("x", roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)),
		),
	)))

	eq(t, roaring.BitmapOf(6, 7, 8, 10), query(NewBoolAndNotQuery(
		NewBoolOrQuery(
			NewTerm("x", roaring.BitmapOf(1, 2)),
			NewTerm("x", roaring.BitmapOf(3, 9))),
		NewBoolAndNotQuery(
			NewTerm("x", roaring.BitmapOf(4, 5)),
			NewTerm("x", roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)),
			NewTerm("x", roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)),
		),
	)))

	eq(t, roaring.BitmapOf(), query(NewBoolAndNotQuery(
		NewTerm("x", roaring.BitmapOf(1, 2, 3, 9)),
		NewTerm("x", roaring.BitmapOf(1, 2, 3, 9)),
	)))

	eq(t, roaring.BitmapOf(), query(NewBoolAndNotQuery(
		NewTerm("x", roaring.BitmapOf(1, 2, 3, 9)),
	)))

	eq(t, roaring.BitmapOf(1, 2, 3, 9), query(NewBoolAndNotQuery(
		NewTerm("x", roaring.BitmapOf()),
		NewTerm("x", roaring.BitmapOf(1, 2, 3, 9)),
	)))

	eq(t, b, query(NewBoolAndQuery(
		NewBoolOrQuery(
			NewTerm("x", a),
			NewTerm("x", b),
		),
		NewTerm("x", b),
		NewTerm("x", c),
		NewTerm("x", d),
		NewTerm("x", e),
	)))

	eq(t, c, query(NewBoolAndQuery(
		NewBoolOrQuery(
			NewTerm("x", a),
			NewTerm("x", b),
			NewBoolAndQuery(
				NewTerm("x", c),
				NewTerm("x", d),
			),
		),
		NewTerm("x", d),
		NewTerm("x", e),
	)))

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

	eq(t, roaring.BitmapOf(4, 5), q.Execute())
}
