package helper

import (
	"errors"
)

func MatrixAdditionWithNumber(matrix [][]float64, number float64) [][]float64 {
	for row := range matrix {
		for col := range matrix[row] {
			matrix[row][col] += number
		}
	}
	return matrix
}

func TransponseMatrix(x [][]float32) [][]float32 {
	out := make([][]float32, len(x[0]))
	for i := 0; i < len(x); i += 1 {
		for j := 0; j < len(x[0]); j += 1 {
			out[j] = append(out[j], x[i][j])
		}
	}
	return out
}

func MatrixMultiplication(matrix1, matrix2 [][]float64) ([][]float64, error) {
	if len(matrix1) == 0 || len(matrix2) == 0 {
		return nil, errors.New("Nil Matrix Length")
	}

	matrix1ColumnLength := len(matrix1[0])
	matrix2RowLength := len(matrix2)

	if matrix1ColumnLength != matrix2RowLength {
		return nil, errors.New("Invalid Matrix Length")
	}

	resultMatrix := make([][]float64, len(matrix1))
	for i := 0; i < len(matrix1); i++ {
		resultMatrix[i] = make([]float64, len(matrix2[0]))
		for j := 0; j < len(matrix2[0]); j++ {
			for k := 0; k < len(matrix2); k++ {
				resultMatrix[i][j] += matrix1[i][k] * matrix2[k][j]
			}
		}
	}

	return resultMatrix, nil
}
