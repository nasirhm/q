package density

import (
	"math"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

type Matrix struct {
	internal matrix.Matrix
}

func New(v ...[]complex128) *Matrix {
	return &Matrix{matrix.New(v...)}
}

func (m *Matrix) Zero(dim int) {
	m.internal = make(matrix.Matrix, dim)
	for i := 0; i < dim; i++ {
		m.internal[i] = []complex128{}
		for j := 0; j < dim; j++ {
			m.internal[i] = append(m.internal[i], complex(0, 0))
		}
	}
}

func (m *Matrix) Add(p float64, q *qubit.Qubit) *Matrix {
	dim := q.Dimension()
	if len(m.internal) < 1 {
		m.Zero(dim)
	}

	op := q.OuterProduct(q).Mul(complex(p, 0))
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			m.internal[i][j] = m.internal[i][j] + op[i][j]
		}
	}

	return m
}

func (m *Matrix) Evlove(u matrix.Matrix) *Matrix {
	return &Matrix{u.Dagger().Apply(m.internal).Apply(u)}
}

func (m *Matrix) Measure(q *qubit.Qubit) complex128 {
	op := q.OuterProduct(q)
	return m.internal.Apply(op).Trace()
}

func (m *Matrix) ExpectedValue(u matrix.Matrix) complex128 {
	return m.internal.Apply(u).Trace()
}

func (m *Matrix) Trace() complex128 {
	return m.internal.Trace()
}

func (m *Matrix) PartialTrace() complex128 {
	// TODO
	return complex(0, 0)
}

func (m *Matrix) NumberOfBit() int {
	mm, _ := m.internal.Dimension()
	log := math.Log2(float64(mm))
	return int(log)
}

func (m *Matrix) Depolarizing(p float64) {
	n := m.NumberOfBit()
	i := gate.I(n).Mul(complex(p/2, 0))
	r := m.internal.Mul(complex(1-p, 0))

	m.internal = i.Add(r)
}

func Flip(p float64, m matrix.Matrix) (matrix.Matrix, matrix.Matrix) {
	e0 := gate.I().Mul(complex(math.Sqrt(p), 0))
	e1 := m.Mul(complex(math.Sqrt(1-p), 0))
	return e0, e1
}

func BitFlip(p float64) (matrix.Matrix, matrix.Matrix) {
	return Flip(p, gate.X())
}

func PhaseFlip(p float64) (matrix.Matrix, matrix.Matrix) {
	return Flip(p, gate.Z())
}

func BitPhaseFlip(p float64) (matrix.Matrix, matrix.Matrix) {
	return Flip(p, gate.Y())
}