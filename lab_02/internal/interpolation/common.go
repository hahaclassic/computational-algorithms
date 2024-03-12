package interpolation

import (
	"errors"
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
	dx       float64 = 1e-5
	UndefNum float64 = -1
)
