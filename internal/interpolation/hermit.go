package interpolation

import (
	"fmt"
	"slices"
	"sort"
)

type Hermit struct {
	points [][]float64
	config [][]float64
}

// CreateNewtonPolinomial() creates a Newton structure that implements interpolation using the Newton polynomial.
// points[i][0] - x coordinate.
// points[i][1] - y coordinate.
// points[i][2] - the first derivative
// points[i][3] - the second derivative
func CreateHermitPolinomial(points [][]float64) *Hermit {
	return &Hermit{points: points}
}

// SetPoints() modifies the set of points from which the approximate value is calculated.
// points[i][0] - x coordinate.
// points[i][1] - y coordinate.
func (h *Hermit) SetPoints(points [][]float64) {
	h.points = points
}

// Calc() calculates the approximate value of y(x) for the degree of the polynomial n.
// x - the input value.
// n - the degree of the Newton polynomial.
func (h *Hermit) Calc(x float64, n int) (float64, error) {
	err := h.configure(x, n)
	if err != nil {
		return UndefNum, err
	}

	h.buildDiff(n)

	args := []float64{}
	for i := 1; i < n+1; i++ {
		args = append(args, h.config[0][i])
	}

	var result, temp float64
	for i := 0; i < n; i++ {
		temp = args[i]
		for j := 0; j < i; j++ {
			temp *= (x - h.config[j][0])
		}
		result += temp
	}

	return result, nil
}

// configure() creates a configuration of the values of the starting points. n + 1 points are selected, as close as possible to x.
// x - the input value.
// n - the degree of the Newton polynomial.
func (h *Hermit) configure(x float64, n int) error {

	if len(h.points) <= n {
		return ErrNotEnoughInputData
	}

	sort.SliceStable(h.points, func(i, j int) bool {
		return h.points[i][0] < h.points[j][0]
	})
	index, _ := slices.BinarySearchFunc(h.points, x, func(point []float64, pointX float64) int {
		if point[0] >= x {
			return 1
		}
		return -1
	})

	h.config = make([][]float64, n+1)
	count := 0
	left, right := index-1, index

	for count < n+1 {
		if left >= 0 && right < len(h.points) && x-h.points[left][0] < h.points[right][0]-x {
			h.config[count] = h.points[left][:2]
			left--
		} else if left >= 0 && right < len(h.points) {
			h.config[count] = h.points[right][:2]
			right++
		} else if right < len(h.points) {
			h.config[count] = h.points[right][:2]
			right++
		} else {
			h.config[count] = h.points[left][:2]
			left--
		}
		count++
	}

	sort.SliceStable(h.config, func(i, j int) bool {
		return h.config[i][0] < h.config[j][0]
	})

	return nil
}

// buildDiff() calculates the values of the separated differences.
// n - the degree of the Newton polynomial.
func (h *Hermit) buildDiff(n int) {

	length := n + 1
	k := 1

	for k <= n {
		idx := len(h.config[0]) - 1
		for i := 0; i < length-1; i++ {
			diff := (h.config[i][idx] - h.config[i+1][idx]) / (h.config[i][0] - h.config[i+k][0])
			h.config[i] = append(h.config[i], diff)
		}
		length--
		k++
	}
}

// SepFiffTable() prints a table of the separated differences of the last operation.
func (h *Hermit) PrintDiffTable() {
	k := len(h.config[0])*16 + 1
	printLine(k)

	fmt.Printf("|       x       |       y       ")
	for i := 2; i < len(h.config[0]); i++ {
		fmt.Printf("| y(x%-2d,..,x%-2d) ", 0, i-1)
	}
	fmt.Println("|")

	printLine(k)

	for i := 0; i < len(h.config); i++ {
		for j := 0; j < len(h.config[i]); j++ {
			fmt.Printf("| ")
			if h.config[i][j] >= 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("%-12f ", h.config[i][j])
			if h.config[i][j] < 0 {
				fmt.Printf(" ")
			}
		}
		for j := 0; j < len(h.config[0])-len(h.config[i]); j++ {
			fmt.Printf("|               ")
		}
		fmt.Println("|")
	}

	printLine(k)
}
