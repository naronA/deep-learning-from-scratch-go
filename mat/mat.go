package mat

import (
	"errors"
	"fmt"
	"sync"

	"github.com/naronA/zero_deeplearning/scalar"
	"github.com/naronA/zero_deeplearning/vec"
)

type Mat interface {
	Element(int, int) interface{}
	Shape() ( int, int )
  SliceRow(r int) vec.Vector
}
type Matrix struct {
	Array   vec.Vector
	Rows    int
	Columns int
}

func (m *Matrix) Shape() (int, int) {
	if m.Rows == 1 {
		return m.Columns, 0
	}
	return m.Rows, m.Columns
}

func (m *Matrix) Element(r int, c int) float64 {
	return m.Array[r*m.Columns+c]
}

func (m *Matrix) SliceRow(r int) vec.Vector {
	slice := make(vec.Vector, m.Columns)
	for i := 0; i < len(slice); i++ {
		slice[i] = m.Array[i+r*m.Columns]
	}
	return slice
}

func (m *Matrix) String() string {
	str := "[\n"
	for i := 0; i < m.Rows; i++ {
		str += fmt.Sprintf("  %v,\n", m.SliceRow(i))
	}
	str += "]"
	return str
}

func Zeros(rows int, cols int) *Matrix {
	zeros := make(vec.Vector, rows*cols)
	for i := range zeros {
		zeros[i] = 0
	}
	mat, err := NewMat64(rows, cols, zeros)
	if err != nil {
		panic(err)
	}
	return mat
}

func ZerosLike(x *Matrix) *Matrix {
	zeros := make(vec.Vector, x.Rows*x.Columns)
	for i := range zeros {
		zeros[i] = 0
	}
	mat, err := NewMat64(x.Rows, x.Columns, zeros)
	if err != nil {
		panic(err)
	}
	return mat
}

func NewMat64(row int, column int, vec vec.Vector) (*Matrix, error) {
	if row == 0 || column == 0 {
		return nil, errors.New("row/columns is zero.")
	}
	return &Matrix{
		Array:   vec,
		Rows:    row,
		Columns: column,
	}, nil
}

func NewRandnMat64(row int, column int) (*Matrix, error) {
	if row == 0 || column == 0 {
		return nil, errors.New("row/columns is zero.")
	}
	vec := vec.Randn(row * column)
	return &Matrix{
		Array:   vec,
		Rows:    row,
		Columns: column,
	}, nil
}

func (m1 *Matrix) NotEqual(m2 *Matrix) bool {
	if m1.Rows == m2.Rows &&
		m1.Columns == m2.Columns &&
		m1.Array.Equal(m2.Array) {
		return false
	}
	return true
}

func (m1 *Matrix) Equal(m2 *Matrix) bool {
	if m1.Rows == m2.Rows &&
		m1.Columns == m2.Columns &&
		m1.Array.Equal(m2.Array) {
		return true
	}
	return false
}

func (m1 *Matrix) DotGo(m2 *Matrix) *Matrix {
	if m1.Columns != m2.Rows {
		return nil
	}
	arys := make([]vec.Vector, m1.Columns)
	wg := &sync.WaitGroup{}
	ch := make(chan int)
	for i := 0; i < m1.Columns; i++ {
		wg.Add(1)
		arys[i] = make(vec.Vector, m1.Rows*m2.Columns)
		go func(ch chan int) {
			defer wg.Done()
			i := <-ch
			for c := 0; c < m2.Columns; c++ {
				for r := 0; r < m1.Rows; r++ {
					arys[i][r*m2.Columns+c] += m1.Element(r, i) * m2.Element(i, c)
				}
			}
		}(ch)
		ch <- i
	}
	wg.Wait()
	sum := vec.Zeros(m1.Rows * m2.Columns)
	for _, ary := range arys {
		sum = sum.Add(ary)
	}
	close(ch)
	return &Matrix{
		Array:   sum,
		Rows:    m1.Rows,
		Columns: m2.Columns,
	}
}

func (m1 *Matrix) Dot(m2 *Matrix) *Matrix {
	if m1.Columns != m2.Rows {
		return nil
	}
	arys := make([]vec.Vector, m1.Columns)
	for i := 0; i < m1.Columns; i++ {
		arys[i] = make(vec.Vector, m1.Rows*m2.Columns)
		for c := 0; c < m2.Columns; c++ {
			for r := 0; r < m1.Rows; r++ {
				arys[i][r*m2.Columns+c] += m1.Element(r, i) * m2.Element(i, c)
			}
		}
	}
	sum := vec.Zeros(m1.Rows * m2.Columns)
	for _, ary := range arys {
		sum = sum.Add(ary)
	}
	return &Matrix{
		Array:   sum,
		Rows:    m1.Rows,
		Columns: m2.Columns,
	}
}

func (m1 *Matrix) Mul(m2 *Matrix) *Matrix {
	mul := m1.Array.Multi(m2.Array)
	return &Matrix{
		Array:   mul,
		Rows:    m1.Rows,
		Columns: m1.Columns,
	}
}

func (m1 *Matrix) Sub(m2 *Matrix) *Matrix {
	// 左辺の行数と、右辺の列数があっていないの掛け算できない
	if m1.Columns != m2.Columns && m1.Rows != m2.Rows {
		return nil
	}

	mat := make([]float64, m1.Rows*m1.Columns)
	for r := 0; r < m1.Rows; r++ {
		for c := 0; c < m2.Columns; c++ {
			index := r*m1.Columns + c
			mat[index] = m1.Element(r, c) - m2.Element(r, c)
		}
	}
	return &Matrix{
		Array:   mat,
		Rows:    m1.Rows,
		Columns: m1.Columns,
	}
}

func (m1 *Matrix) Add(m2 *Matrix) *Matrix {
	// 左辺の行数と、右辺の列数があっていないの掛け算できない
	if m1.Columns != m2.Columns && m1.Rows != m2.Rows {
		return nil
	}

	mat := make([]float64, m1.Rows*m1.Columns)
	for r := 0; r < m1.Rows; r++ {
		for c := 0; c < m2.Columns; c++ {
			index := r*m1.Columns + c
			mat[index] = m1.Element(r, c) + m2.Element(r, c)
		}
	}
	return &Matrix{
		Array:   mat,
		Rows:    m1.Rows,
		Columns: m1.Columns,
	}
}

func (m1 *Matrix) AddBroadCast(m2 *Matrix) *Matrix {
	// 左辺の行数と、右辺の列数があっていないの掛け算できない
	if m1.Columns != m2.Columns {
		return nil
	}

	mat := make([]float64, m1.Rows*m1.Columns)
	for r := 0; r < m1.Rows; r++ {
		for c := 0; c < m1.Columns; c++ {
			index := r*m1.Columns + c
			mat[index] = m1.Element(r, c) + m2.Element(0, c)
		}
	}
	return &Matrix{
		Array:   mat,
		Rows:    m1.Rows,
		Columns: m1.Columns,
	}
}
func (m *Matrix) AddAll(a float64) *Matrix {

	mat := make([]float64, m.Rows*m.Columns)
	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Columns; c++ {
			index := r*m.Columns + c
			mat[index] = m.Element(r, c) + a
		}
	}
	return &Matrix{
		Array:   mat,
		Rows:    m.Rows,
		Columns: m.Columns,
	}
}

func (m *Matrix) MulAll(a float64) *Matrix {
	// 左辺の行数と、右辺の列数があっていないの掛け算できない
	mat := make([]float64, m.Rows*m.Columns)
	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Columns; c++ {
			index := r*m.Columns + c
			mat[index] = a * m.Element(r, c)
		}
	}
	return &Matrix{
		Array:   mat,
		Rows:    m.Rows,
		Columns: m.Columns,
	}
}

func Sigmoid(m *Matrix) *Matrix {
	mat := make([]float64, m.Rows*m.Columns)
	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Columns; c++ {
			index := r*m.Columns + c
			mat[index] = scalar.Sigmoid(m.Element(r, c))
		}
	}
	return &Matrix{
		Array:   mat,
		Rows:    m.Rows,
		Columns: m.Columns,
	}
}

func Relu(m *Matrix) *Matrix {
	mat := make([]float64, m.Rows*m.Columns)
	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Columns; c++ {
			index := r*m.Columns + c
			mat[index] = scalar.Relu(m.Element(r, c))
		}
	}
	return &Matrix{
		Array:   mat,
		Rows:    m.Rows,
		Columns: m.Columns,
	}
}

func Log(m *Matrix) *Matrix {
	log := vec.Log(m.Array)
	return &Matrix{
		Array:   log,
		Rows:    m.Rows,
		Columns: m.Columns,
	}
}

func Sum(m *Matrix) float64 {
	return vec.Sum(m.Array)
}

func ArgMax(x *Matrix) []int {
	r := make([]int, x.Rows)
	for i := 0; i < x.Rows; i++ {
		row := x.SliceRow(i)
		r[i] = vec.ArgMax(row)
	}
	return r
}

func Softmax(x *Matrix) *Matrix {
	m := vec.Vector{}
	for i := 0; i < x.Rows; i++ {
		xRow := x.SliceRow(i)
		m = append(m, vec.Softmax(xRow)...)
	}
	r, err := NewMat64(x.Rows, x.Columns, m)
	if err != nil {
		panic(err)
	}
	return r
}

func CrossEntropyError(y, t *Matrix) float64 {
	r := make(vec.Vector, y.Rows)
	for i := 0; i < y.Rows; i++ {
		yRow := y.SliceRow(i)
		tRow := t.SliceRow(i)
		r[i] = vec.CrossEntropyError(yRow, tRow)
	}
	return vec.Sum(r) / float64(y.Rows)
}

func NumericalGradient(f func(vec.Vector) float64, x *Matrix) *Matrix {
	grad := vec.NumericalGradient(f, x.Array)
	mat, _ := NewMat64(x.Rows, x.Columns, grad)
	return mat
}
