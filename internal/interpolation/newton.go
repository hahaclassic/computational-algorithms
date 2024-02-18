package interpolation

import (
	"errors"
	"fmt"
	"slices"
	"sort"
)

var (
	ErrNotEnoughInputData = errors.New("not enough input data") // decrease the degree of the polynomial or increase the number of input points")
)

const (
	delta    float64 = 1e-7
	UndefNum float64 = -1
)

type Newton struct {
	points [][]float64
	config [][]float64
}

// CreateNewtonPolinomial() creates a Newton structure that implements interpolation using the Newton polynomial.
// points[i][0] - x coordinate.
// points[i][1] - y coordinate.
func CreateNewtonPolinomial(points [][]float64) (*Newton, error) {
	newton := &Newton{
		points: make([][]float64, len(points)),
	}
	for i := 0; i < len(points); i++ {
		if len(points[i]) < 2 {
			return nil, ErrNotEnoughInputData
		}
		newton.points[i] = make([]float64, 2)
		copy(newton.points[i], points[i][:2])
	}

	return newton, nil
}

// SetPoints() modifies the set of points from which the approximate value is calculated.
// points[i][0] - x coordinate.
// points[i][1] - y coordinate.
func (newton *Newton) SetPoints(points [][]float64) error {
	newton.points = make([][]float64, len(points))
	for i := 0; i < len(points); i++ {
		if len(points[i]) < 2 {
			return ErrNotEnoughInputData
		}
		newton.points[i] = make([]float64, 2)
		copy(newton.points[i], points[i][:2])
	}

	return nil
}

// Calc() calculates the approximate value of y(x) for the degree of the polynomial n.
// x - the input value.
// n - the degree of the Newton polynomial.
func (newton *Newton) Calc(x float64, n int) (float64, error) {
	err := newton.configure(x, n)
	if err != nil {
		return UndefNum, err
	}

	newton.buildDiff()

	args := []float64{}
	for i := 1; i <= n+1; i++ {
		args = append(args, newton.config[0][i])
	}

	var result, temp float64
	for i := 0; i < len(args); i++ {
		temp = args[i]
		for j := 0; j < i; j++ {
			temp *= (x - newton.config[j][0])
		}
		result += temp
	}

	return result, nil
}

// configure() creates a configuration of the values of the starting points. n + 1 points are selected, as close as possible to x.
// x - the input value.
// n - the degree of the Newton polynomial.
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

// buildDiff() calculates the values of the divided differences.
func (newton *Newton) buildDiff() {

	numOfNodes := len(newton.config)
	n := numOfNodes - 1 // n - the degree of the Newton polynomial.

	for k := 1; k <= n; k++ {
		idx := len(newton.config[0]) - 1
		for i := 0; i < numOfNodes-1; i++ {
			diff := (newton.config[i][idx] - newton.config[i+1][idx]) / (newton.config[i][0] - newton.config[i+k][0])
			newton.config[i] = append(newton.config[i], diff)
		}
		numOfNodes--
	}
}

// SepFiffTable() prints a table of the divided differences of the last operation.
func (newton *Newton) PrintDiffTable() {
	k := len(newton.config[0])*16 + 1
	printLine(k)

	fmt.Printf("|       x       |       y       ")
	for i := 2; i < len(newton.config[0]); i++ {
		fmt.Printf("| y(x%-2d,..,x%-2d) ", 0, i-1)
	}
	fmt.Println("|")

	printLine(k)

	for i := 0; i < len(newton.config); i++ {
		for j := 0; j < len(newton.config[i]); j++ {
			fmt.Printf("| ")
			if newton.config[i][j] >= 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("%-12f ", newton.config[i][j])
			if newton.config[i][j] < 0 {
				fmt.Printf(" ")
			}
		}
		for j := 0; j < len(newton.config[0])-len(newton.config[i]); j++ {
			fmt.Printf("|               ")
		}
		fmt.Println("|")
	}

	printLine(k)
}

func printLine(k int) {
	for i := 0; i < k; i++ {
		fmt.Print("-")
	}
	fmt.Println()
}
