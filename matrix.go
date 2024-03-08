package qr

import "fmt"

type matrix[T any] struct {
	mat    [][]T
	width  int
	height int
}

func newMatrix[T any](width, height int) *matrix[T] {
	mat := make([][]T, width)

	for i := 0; i < width; i++ {
		mat[i] = make([]T, height)
	}

	return &matrix[T]{
		mat:    mat,
		width:  width,
		height: height,
	}
}

func zero[T any]() T {
	return *new(T)
}

func (m *matrix[T]) Init(val T) {
	for i := 0; i < m.width; i++ {
		for j := 0; j < m.height; j++ {
			m.mat[i][j] = val
		}
	}
}

func (m *matrix[T]) Width() int {
	return m.width
}

func (m *matrix[T]) Height() int {
	return m.height
}

func (m *matrix[T]) At(w, h int) (T, error) {
	if w < 0 || w > m.width-1 {
		return zero[T](), fmt.Errorf("width out of range")
	}

	if h < 0 || h >= m.height {
		return zero[T](), fmt.Errorf("height out of range")
	}

	return m.mat[w][h], nil
}

func (m *matrix[T]) Set(w, h int, val T) (T, error) {
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

func (m *matrix[T]) RowAt(rowIdx int) ([]T, error) {
	if rowIdx < 0 || rowIdx > m.height-1 {
		return nil, fmt.Errorf("row index out of range")
	}

	return m.mat[rowIdx], nil
}

func (m *matrix[T]) ColumnAt(colIdx int) ([]T, error) {
	if colIdx < 0 || colIdx > m.width-1 {
		return nil, fmt.Errorf("column index out of range")
	}

	row := make([]T, m.height)
	for i := 0; i < m.height; i++ {
		row[i] = m.mat[i][colIdx]
	}

	return row, nil
}

func (m *matrix[T]) GetMatrix() [][]T {
	return m.mat
}

func (m *matrix[T]) SetMatrix(mat [][]T) error {
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

func (m *matrix[T]) Expand(n int) error {
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

func (m *matrix[T]) PrintMatrix() {
	for _, row := range m.mat {
		for _, val := range row {
			fmt.Printf("%v ", val)
		}
		fmt.Println()
	}
}
