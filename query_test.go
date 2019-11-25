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
		q := And(
			Term("x", x),
			Term("y", y),
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
		q := Or(
			Term("x", x),
			Term("y", y),
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

	eq(t, a, query(Term("x", a)))
	eq(t, b, query(Term("x", b)))
	eq(t, c, query(Term("x", c)))
	eq(t, d, query(Term("x", d)))
	eq(t, e, query(Term("x", e)))

	eq(t, b, query(Or(
		Term("a", a),
		Term("b", b),
	)))

	eq(t, c, query(Or(
		Term("a", a),
		Term("b", b),
		Term("c", c),
	)))

	eq(t, e, query(Or(
		Term("a", a),
		Term("b", b),
		Term("c", c),
		Term("d", d),
		Term("e", e),
	)))

	eq(t, a, query(And(
		Term("a", a),
		Term("b", b),
		Term("c", c),
		Term("d", d),
		Term("e", e),
	)))

	eq(t, roaring.BitmapOf(4, 6, 7, 8, 10), query(AndNot(
		Term("x", roaring.BitmapOf(1, 2, 3, 9)),
		Or(
			Term("x", roaring.BitmapOf(3, 4)),
			Term("x", roaring.BitmapOf(1, 2, 3, 6, 7, 8, 9, 10)),
		),
	)))
	eq(t, roaring.BitmapOf(6, 7, 8, 10), query(AndNot(
		Term("x", roaring.BitmapOf(1, 2, 3, 9)),
		AndNot(
			Term("x", roaring.BitmapOf(4, 5)),
			Term("x", roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)),
			Term("x", roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)),
		),
	)))

	eq(t, roaring.BitmapOf(6, 7, 8, 10), query(AndNot(
		Or(
			Term("x", roaring.BitmapOf(1, 2)),
			Term("x", roaring.BitmapOf(3, 9))),
		AndNot(
			Term("x", roaring.BitmapOf(4, 5)),
			Term("x", roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)),
			Term("x", roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)),
		),
	)))

	eq(t, roaring.BitmapOf(), query(AndNot(
		Term("x", roaring.BitmapOf(1, 2, 3, 9)),
		Term("x", roaring.BitmapOf(1, 2, 3, 9)),
	)))

	eq(t, roaring.BitmapOf(), query(AndNot(
		Term("x", roaring.BitmapOf(1, 2, 3, 9)),
	)))

	eq(t, roaring.BitmapOf(1, 2, 3, 9), query(AndNot(
		Term("x", roaring.BitmapOf()),
		Term("x", roaring.BitmapOf(1, 2, 3, 9)),
	)))

	eq(t, b, query(And(
		Or(
			Term("x", a),
			Term("x", b),
		),
		Term("x", b),
		Term("x", c),
		Term("x", d),
		Term("x", e),
	)))

	eq(t, c, query(And(
		Or(
			Term("x", a),
			Term("x", b),
			And(
				Term("x", c),
				Term("x", d),
			),
		),
		Term("x", d),
		Term("x", e),
	)))

	q := AndNot(
		Or(
			Term("a", roaring.BitmapOf(1, 2)),
			Term("b", roaring.BitmapOf(3, 9))),
		And(
			Term("c", roaring.BitmapOf(4, 5)),
			Term("d", roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)),
		),
		Term("e", roaring.BitmapOf(4, 5, 6)),
	)

	eq(t, roaring.BitmapOf(4, 5), q.Execute())
}
