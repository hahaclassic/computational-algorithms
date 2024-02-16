package interpolation

import (
	"errors"
	"slices"
	"sort"
)

var (
	ErrNotEnoughInputData = errors.New("not enough input data") // decrease the degree of the polynomial or increase the number of input points")
)

const (
	UndefNum float64 = -1
)

type Newton struct {
	points [][]float64
	config [][]float64
}

func CreateNewtonPolinomial(points [][]float64) *Newton {
	return &Newton{points: points}
}

func (newton *Newton) Calc(x float64, n int) (float64, error) {
	err := newton.configure(x, n)
	if err != nil {
		return UndefNum, err
	}

	newton.buildDiff(n)

	args := []float64{}
	for i := 1; i < n+1; i++ {
		args = append(args, newton.config[0][i])
	}

	var result, temp float64
	for i := 0; i < n; i++ {
		temp = args[i]
		for j := 0; j < i; j++ {
			temp *= (x - newton.config[j][0])
		}
		result += temp
	}

	return result, nil
}

func (newton *Newton) configure(x float64, n int) error {

	if len(newton.points) <= n {
		return ErrNotEnoughInputData
	}

	sort.SliceStable(newton.points, func(i, j int) bool {
		return newton.points[i][0] < newton.points[j][0]
	})
	index, _ := slices.BinarySearchFunc(newton.points, x, func(point []float64, pointX float64) int {
		if point[0] >= x {
			return 1
		}
		return -1
	})

	newton.config = make([][]float64, n+1)
	count := 0
	left, right := index-1, index

	for count < n+1 {
		if left >= 0 && right < len(newton.points) && x-newton.points[left][0] < newton.points[right][0]-x {
			newton.config[count] = newton.points[left][:2]
			left--
		} else if left >= 0 && right < len(newton.points) {
			newton.config[count] = newton.points[right][:2]
			right++
		} else if right < len(newton.points) {
			newton.config[count] = newton.points[right][:2]
			right++
		} else {
			newton.config[count] = newton.points[left][:2]
			left--
		}
		count++
	}

	sort.SliceStable(newton.config, func(i, j int) bool {
		return newton.config[i][0] < newton.config[j][0]
	})

	return nil
}

func (newton *Newton) buildDiff(n int) {

	length := n + 1
	k := 1

	for k <= n {
		idx := len(newton.config[0]) - 1
		for i := 0; i < length-1; i++ {
			diff := (newton.config[i][idx] - newton.config[i+1][idx]) / (newton.config[i][0] - newton.config[i+k][0])
			newton.config[i] = append(newton.config[i], diff)
		}
		length--
		k++
	}
}
