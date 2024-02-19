package matrix

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
