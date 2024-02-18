package interpolation

import (
	"fmt"
	"math"
	"slices"
	"sort"

	"github.com/hahaclassic/computational-algorithms.git/internal/format"
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
	sort.SliceStable(newton.points, func(i, j int) bool {
		return newton.points[i][0] < newton.points[j][0]
	})

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
	sort.SliceStable(newton.points, func(i, j int) bool {
		return newton.points[i][0] < newton.points[j][0]
	})

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

	return newton.result(x), nil
}

// FindRoot() finds root of the function (y == 0)
// n - the degree of the Newton polynomial.
func (newton *Newton) FindRoot(n int) (float64, error) {

	idx := -1
	for i := 0; i < len(newton.points)-1; i++ {
		if math.Abs(newton.points[i][1]) < delta {
			return newton.points[i][0], nil
		}
		if newton.points[i][1]*newton.points[i+1][1] < 0 {
			idx = i + 1
			break
		}
	}
	if idx == -1 {
		return UndefNum, ErrNoRoot
	}

	source := newton.points
	inverted, _ := Inverse(newton.points)
	newton.points = inverted

	newton.fillConfig(0, n, idx)
	newton.buildDiff()
	newton.points = source

	return newton.result(0), nil
}

// configure() creates a configuration of the values of the starting points. n + 1 points are selected, as close as possible to x.
// x - the input value.
// n - the degree of the Newton polynomial.
func (newton *Newton) configure(x float64, n int) error {

	if len(newton.points) <= n {
		return ErrNotEnoughInputData
	}

	index, _ := slices.BinarySearchFunc(newton.points, x, func(point []float64, pointX float64) int {
		if point[0] >= x {
			return 1
		}
		return -1
	})

	newton.fillConfig(x, n, index)

	return nil
}

// n - degree of the polynomial
func (newton *Newton) fillConfig(x float64, n int, idx int) {
	leftNodes := [][]float64{}
	rightNodes := [][]float64{}
	left, right := idx-1, idx

	for count := 0; count < n+1; count++ {
		if left >= 0 && right < len(newton.points) &&
			math.Abs(x-newton.points[left][0]) < math.Abs(newton.points[right][0]-x) {
			leftNodes = append(leftNodes, newton.points[left][:2])
			left--
		} else if left >= 0 && right < len(newton.points) {
			rightNodes = append(rightNodes, newton.points[right][:2])
			right++
		} else if right < len(newton.points) {
			rightNodes = append(rightNodes, newton.points[right][:2])
			right++
		} else {
			leftNodes = append(leftNodes, newton.points[left][:2])
			left--
		}
	}
	slices.Reverse(leftNodes)
	newton.config = leftNodes
	newton.config = append(newton.config, rightNodes...)
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

func (newton *Newton) result(x float64) float64 {
	args := []float64{}
	for i := 1; i < len(newton.config[0]); i++ {
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

	return result
}

// SepFiffTable() prints a table of the divided differences of the last operation.
func (newton *Newton) PrintDiffTable() {
	if len(newton.config) == 0 {
		return
	}

	k := len(newton.config[0])*18 + 1
	fmt.Println()
	format.PrintLine(k)

	fmt.Printf("|        x        |        y        ")
	for i := 2; i < len(newton.config[0]); i++ {
		fmt.Printf("|  y(x%-2d,..,x%-2d)  ", 0, i-1)
	}
	fmt.Println("|")

	format.PrintLine(k)

	for i := 0; i < len(newton.config); i++ {
		for j := 0; j < len(newton.config[i]); j++ {
			fmt.Printf("| ")
			if newton.config[i][j] >= 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("%-14f ", newton.config[i][j])
			if newton.config[i][j] < 0 {
				fmt.Printf(" ")
			}
		}
		for j := 0; j < len(newton.config[0])-len(newton.config[i]); j++ {
			fmt.Printf("|                 ")
		}
		fmt.Println("|")
	}

	format.PrintLine(k)
	fmt.Println()
}
