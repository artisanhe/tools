package benchmark

import (
	"testing"
)

type S struct{ a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, q, r, s, t, u, v, w, x, y, z int64 }

func Value() int64 {
	return 0
}

func (S) Value() int64 {
	return 0
}

func (t S) ValueWithVar() int64 {
	return 0
}

func (t S) ValueAccess() int64 {
	return t.a
}

func (*S) Pointer() int64 {
	return 0
}

func (t *S) PointerAccess() int64 {
	return t.a
}

func Benchmark_Access(b *testing.B) {
	s := S{}
	v := int64(0)
	for i := 0; i < b.N; i++ {
		v = s.a
	}
	b.Log(v)
}

func Benchmark_Value(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Value()
	}
}

func Benchmark_ValueReceiver(b *testing.B) {
	for i := 0; i < b.N; i++ {
		(S{}).Value()
	}
}

func Benchmark_ValueReceiverWithVar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		(S{}).ValueWithVar()
	}
}

func Benchmark_ValueReceiverWithAccess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		(S{}).ValueAccess()
	}
}

func Benchmark_ValueReceiverFromSomeId(b *testing.B) {
	v := S{}
	for i := 0; i < b.N; i++ {
		v.Value()
	}
}

func Benchmark_ValueReceiverWithAccessFromSomeId(b *testing.B) {
	v := S{}
	for i := 0; i < b.N; i++ {
		v.ValueAccess()
	}
}

func Benchmark_ValueReceiverForPointer(b *testing.B) {
	p := &S{}
	for i := 0; i < b.N; i++ {
		p.Value()
	}
}

func Benchmark_PointerReceiver(b *testing.B) {
	p := &S{}
	for i := 0; i < b.N; i++ {
		p.Pointer()
	}
}
