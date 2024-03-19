package qr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatrixWidth(t *testing.T) {
	assert := assert.New(t)

	m := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	mat := newMatrix[int](3, 3)
	mat.SetMatrix(m)

	assert.Equal(3, mat.Width(), "width should match")
}

func TestMatrixHeight(t *testing.T) {
	assert := assert.New(t)

	m := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	mat := newMatrix[int](3, 3)
	mat.SetMatrix(m)

	assert.Equal(3, mat.Height(), "height should match")
}

func TestMatrixAt(t *testing.T) {
	assert := assert.New(t)

	m := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	mat := newMatrix[int](3, 3)
	mat.SetMatrix(m)

	tests := []struct {
		name     string
		w, h     int
		expected int
		err      string
	}{
		{
			name:     "WithinRange",
			w:        1,
			h:        2,
			expected: 6,
			err:      "",
		},
		{
			name:     "WidthOutOfRange",
			w:        3,
			h:        1,
			expected: 0,
			err:      "width out of range",
		},
		{
			name:     "HeightOutOfRange",
			w:        1,
			h:        3,
			expected: 0,
			err:      "height out of range",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := mat.At(test.w, test.h)

			if err != nil {
				assert.Equal(test.err, err.Error(), "Error messages should match")
			} else {
				assert.Equal(test.expected, actual, "Cell values should match")
			}
		})
	}
}

func TestMatrixSet(t *testing.T) {
	assert := assert.New(t)

	m := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	mat := newMatrix[int](3, 3)
	mat.SetMatrix(m)

	tests := []struct {
		name            string
		w, h, val       int
		expectedPrevVal int
		err             string
	}{
		{
			name:            "WithinRange",
			w:               1,
			h:               2,
			val:             10,
			expectedPrevVal: 6,
			err:             "",
		},
		{
			name:            "WidthOutOfRange",
			w:               3,
			h:               1,
			val:             10,
			expectedPrevVal: 0,
			err:             "width out of range",
		},
		{
			name:            "HeightOutOfRange",
			w:               1,
			h:               3,
			val:             10,
			expectedPrevVal: 0,
			err:             "height out of range",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prevVal, err := mat.Set(test.w, test.h, test.val)

			if err != nil {
				assert.Equal(test.err, err.Error(), "Error messages should match")
			} else {
				assert.Equal(test.expectedPrevVal, prevVal, "Previous value returned should match")
				newVal, _ := mat.At(test.w, test.h)
				assert.Equal(test.val, newVal, "New value set should match")
			}
		})
	}
}

func TestMatrixRowAt(t *testing.T) {
	assert := assert.New(t)

	m := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	mat := newMatrix[int](3, 3)
	mat.SetMatrix(m)

	tests := []struct {
		name     string
		rowIdx   int
		expected []int
		err      string
	}{
		{
			name:     "WithinRange",
			rowIdx:   1,
			expected: []int{4, 5, 6},
			err:      "",
		},
		{
			name:     "RowOutOfRange",
			rowIdx:   3,
			expected: nil,
			err:      "row index out of range",
		},
		{
			name:     "NegativeRowOutOfRange",
			rowIdx:   -1,
			expected: nil,
			err:      "row index out of range",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := mat.RowAt(test.rowIdx)

			if err != nil {
				assert.Equal(test.err, err.Error(), "error messages should match")
			} else {
				assert.Equal(test.expected, actual, "rows should match")
			}
		})
	}
}

func TestMatrixColAt(t *testing.T) {
	assert := assert.New(t)

	m := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	mat := newMatrix[int](3, 3)
	mat.SetMatrix(m)

	tests := []struct {
		name     string
		colIdx   int
		expected []int
		err      string
	}{
		{
			name:     "WithinRange",
			colIdx:   1,
			expected: []int{2, 5, 8},
			err:      "",
		},
		{
			name:     "ColumnOutOgRange",
			colIdx:   3,
			expected: nil,
			err:      "column index out of range",
		},
		{
			name:     "NegativeColumn",
			colIdx:   -1,
			expected: nil,
			err:      "column index out of range",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := mat.ColumnAt(test.colIdx)

			if err != nil {
				assert.Equal(test.err, err.Error(), "error messages should match")
			} else {
				assert.Equal(test.expected, actual, "columns should match")
			}
		})
	}
}
