package matrix

import "fmt"

type Matrix[T any] struct {
	mat    [][]T
	width  int
	height int
}

func NewMatrix[T any](width, height int) *Matrix[T] {
	mat := make([][]T, width)

	for i := 0; i < width; i++ {
		mat[i] = make([]T, height)
	}

	return &Matrix[T]{
		mat:    mat,
		width:  width,
		height: height,
	}
}

func zero[T any]() T {
	return *new(T)
}

func (m *Matrix[T]) Init(val T) {
	for i := 0; i < m.width; i++ {
		for j := 0; j < m.height; j++ {
			m.mat[i][j] = val
		}
	}
}

func (m *Matrix[T]) Width() int {
	return m.width
}

func (m *Matrix[T]) Height() int {
	return m.height
}

func (m *Matrix[T]) At(w, h int) (T, error) {
	if w < 0 || w > m.width-1 {
		return zero[T](), fmt.Errorf("width out of range")
	}

	if h < 0 || h >= m.height {
		return zero[T](), fmt.Errorf("height out of range")
	}

	return m.mat[w][h], nil
}

func (m *Matrix[T]) Set(w, h int, val T) (T, error) {
	if w < 0 || w > m.width-1 {
		return zero[T](), fmt.Errorf("width out of range")
	}

	if h < 0 || h > m.height-1 {
		return zero[T](), fmt.Errorf("height out of range")
	}

	prevVal := m.mat[w][h]
	m.mat[w][h] = val

	return prevVal, nil
}

func (m *Matrix[T]) RowAt(rowIdx int) ([]T, error) {
	if rowIdx < 0 || rowIdx > m.height-1 {
		return nil, fmt.Errorf("row index out of range")
	}

	return m.mat[rowIdx], nil
}

func (m *Matrix[T]) ColumnAt(colIdx int) ([]T, error) {
	if colIdx < 0 || colIdx > m.width-1 {
		return nil, fmt.Errorf("column index out of range")
	}

	row := make([]T, m.height)
	for i := 0; i < m.height; i++ {
		row[i] = m.mat[i][colIdx]
	}

	return row, nil
}

func (m *Matrix[T]) GetMatrix() [][]T {
	return m.mat
}

func (m *Matrix[T]) SetMatrix(mat [][]T) error {
	if m.width != len(mat) {
		return fmt.Errorf("matrices witdth does not match")
	}

	if m.height != len(mat[0]) {
		return fmt.Errorf("matrices height does not match")
	}

	for i := 0; i < m.width; i++ {
		for j := 0; j < m.height; j++ {
			m.mat[i][j] = mat[i][j]
		}
	}

	return nil
}

func (m *Matrix[T]) Expand(n int) error {
	if n < 0 {
		return fmt.Errorf("invalid expansion unit")
	}

	expandedMat := make([][]T, m.width+2*n)
	for i := range expandedMat {
		expandedMat[i] = make([]T, m.height+2*n)
	}

	for i := 0; i < m.width; i++ {
		for j := 0; j < m.height; j++ {
			expandedMat[i+n][j+n] = m.mat[i][j]
		}
	}

	m.mat = expandedMat
	return nil
}

func (m *Matrix[T]) PrintMatrix() {
	for _, row := range m.mat {
		for _, val := range row {
			fmt.Printf("%v ", val)
		}
		fmt.Println()
	}
}
