package interpolation

import (
	"errors"
	"math"
)

var (
	ErrNotEnoughInputData      = errors.New("not enough input data") // decrease the degree of the polynomial or increase the number of input points
	ErrCantInverseFunc         = errors.New("one of the derivatives is equal to zero")
	ErrNoRoot                  = errors.New("at this interval, the function has no valid roots")
	ErrInvalidPolynomialDegree = errors.New("invalid polynomial degree")
	ErrInvalidNumDerivates     = errors.New("invalid num of derivates")
)

const (
	delta    float64 = 1e-7
	UndefNum float64 = -1
)

// Inverse() Changes the data dependency y(x) to x(y).
// If derivatives are present in the data, they are replaced with the reverse ones.
// Before:
// points[i][0] - x coordinate.
// points[i][1] - y coordinate.
// points[i][2] - the first derivative y' (optional)
// points[i][3] - the second derivative y” (optional)
// After:
// points[i][0] - y coordinate.
// points[i][1] - x coordinate.
// points[i][2] - the first derivative x' (optional)
// points[i][3] - the second derivative x” (optional)
func Inverse(points [][]float64) ([][]float64, error) {
	result := make([][]float64, len(points))
	for i := 0; i < len(result); i++ {
		if len(points[i]) < 2 {
			return nil, ErrNotEnoughInputData
		}
		result[i] = make([]float64, len(points[i]))
	}

	for i := 0; i < len(points); i++ {
		result[i][0] = points[i][1]
		result[i][1] = points[i][0]
		if len(result[i]) > 2 {
			if math.Abs(points[i][2]) < delta {
				return nil, ErrCantInverseFunc
			} else {
				result[i][2] = 1 / points[i][2]
			}
		}
		if len(result[i]) > 3 {
			result[i][3] = -(points[i][3] / math.Pow(points[i][2], 3))
		}
	}

	return result, nil
}
