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
		return zero[T](), fmt.Errorf("Width out of range")
	}

	if h < 0 || h >= m.height {
		return zero[T](), fmt.Errorf("Height out of range")
	}

	return m.mat[w][h], nil
}

func (m *Matrix[T]) Set(w, h int, val T) (T, error) {
	if w < 0 || w > m.width-1 {
		return zero[T](), fmt.Errorf("Width out of range")
	}

	if h < 0 || h > m.height-1 {
		return zero[T](), fmt.Errorf("Height out of range")
	}

	prevVal := m.mat[w][h]
	m.mat[w][h] = val

	return prevVal, nil
}

func (m *Matrix[T]) RowAt(rowIdx int) ([]T, error) {
	if rowIdx < 0 || rowIdx > m.height-1 {
		return nil, fmt.Errorf("Row index out of range")
	}

	return m.mat[rowIdx], nil
}

func (m *Matrix[T]) ColumnAt(colIdx int) ([]T, error) {
	if colIdx < 0 || colIdx > m.width-1 {
		return nil, fmt.Errorf("Column index out of range")
	}

	row := make([]T, m.height)
	for i := 0; i < m.height; i++ {
		row[i] = m.mat[i][colIdx]
	}

	return row, nil
}

func (m *Matrix[T]) setMatrix(mat [][]T) error {
	if m.width != len(mat) {
		return fmt.Errorf("Matrices witdth does not match")
	}

	if m.height != len(mat[0]) {
		return fmt.Errorf("Matrices height does not match")
	}

	for i := 0; i < m.width; i++ {
		for j := 0; j < m.height; j++ {
			m.mat[i][j] = mat[i][j]
		}
	}

	return nil
}
