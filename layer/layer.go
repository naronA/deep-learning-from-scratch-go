package layer

import (
	"github.com/naronA/zero_deeplearning/mat"
	"github.com/naronA/zero_deeplearning/vec"
)

type Affine struct {
	W  *mat.Matrix
	B  *mat.Matrix
	X  *mat.Matrix
	DW *mat.Matrix
	DB *mat.Matrix
}

func NewAffine(w, b *mat.Matrix) *Affine {
	return &Affine{
		W: w,
		B: b,
	}
}

func (a *Affine) Forward(x *mat.Matrix) *mat.Matrix {
	a.X = x
	out := mat.Dot(x, a.W).Add(a.B)
	return out
}

func (a *Affine) Backward(dout *mat.Matrix) *mat.Matrix {
	dx := mat.Dot(dout, a.W.T())
	a.DW = mat.Dot(a.X.T(), dout)
	a.DB = mat.Sum(dout, 0)
	return dx
}

type Sigmoid struct {
	Out *mat.Matrix
}

func NewSigmoid() *Sigmoid {
	return &Sigmoid{}
}

func (s *Sigmoid) Forward(x *mat.Matrix) *mat.Matrix {
	minusX := x.Mul(-1.0)
	exp := mat.Exp(minusX)
	plusX := exp.Add(1.0)
	out := mat.Pow(plusX, -1)
	s.Out = out
	return out
}

func (s *Sigmoid) Backward(dout *mat.Matrix) *mat.Matrix {
	minus := s.Out.Mul(-1.0)
	sub := minus.Add(1.0)
	mul := dout.Mul(sub)
	dx := mul.Mul(s.Out)
	return dx
}

type Relu struct {
	mask []bool
}

func NewRelu() *Relu {
	return &Relu{}
}

func (r *Relu) Forward(x *mat.Matrix) *mat.Matrix {
	v := x.Vector
	r.mask = make([]bool, len(v))
	out := vec.ZerosLike(v)
	for i, e := range v {
		if e <= 0 {
			r.mask[i] = true
			out[i] = 0
		} else {
			out[i] = e
		}
	}

	return &mat.Matrix{
		Vector:  out,
		Rows:    x.Rows,
		Columns: x.Columns,
	}
}

func (r *Relu) Backward(dout *mat.Matrix) *mat.Matrix {
	v := dout.Vector
	dv := vec.ZerosLike(v)
	for i, e := range v {
		if r.mask[i] {
			dv[i] = 0
		} else {
			dv[i] = e
		}
	}
	dx := &mat.Matrix{
		Vector:  dv,
		Rows:    dout.Rows,
		Columns: dout.Columns,
	}
	return dx
}

type SoftmaxWithLoss struct {
	loss float64
	y    *mat.Matrix
	t    *mat.Matrix
}

func NewSfotmaxWithLoss() *SoftmaxWithLoss {
	return &SoftmaxWithLoss{}
}

func (s *SoftmaxWithLoss) Forward(x, t *mat.Matrix) float64 {
	s.t = t
	s.y = mat.Softmax(x)
	s.loss = mat.CrossEntropyError(s.y, s.t)
	return s.loss
}

func (s *SoftmaxWithLoss) Backward(_ float64) *mat.Matrix {
	batchSize, _ := s.t.Shape()
	sub := s.y.Sub(s.t)
	dx := sub.Div(float64(batchSize))
	return dx
}