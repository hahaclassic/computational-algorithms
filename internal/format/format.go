package format

import "strconv"

func MatrixAtof(data [][]string) ([][]float64, error) {

	result := make([][]float64, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = make([]float64, len(data[i]))
	}

	var err error

	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			result[i][j], err = strconv.ParseFloat(data[i][j], 64)
			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func MatrixTranspose(data [][]float64) [][]float64 {

	matrix := make([][]float64, len(data[0]))
	for i := 0; i < len(matrix); i++ {
		matrix[i] = make([]float64, len(data))
	}

	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			matrix[j][j] = data[i][j]
		}
	}

	return matrix
}
