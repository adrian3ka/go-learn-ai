package helper

import "math"

func Softmax(matrix []float64) []float64 {
	divisor := float64(1)

	for idx := range matrix {
		divisor += math.Pow(math.E, matrix[idx])
	}
	for idx := range matrix {
		matrix[idx] = matrix[idx] / divisor
	}
	return matrix
}
