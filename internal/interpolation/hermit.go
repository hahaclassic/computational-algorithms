package interpolation

import (
	"fmt"
	"math"
	"slices"
	"sort"
)

type Hermit struct {
	points         [][]float64
	config         [][]float64
	differences    [][]float64
	numOfNodes     int
	numDerivatives int // number of derivatives at a point
}

// CreateNewtonPolinomial() creates a Newton structure that implements interpolation using the Newton polynomial.
// points[i][0] - x coordinate.
// points[i][1] - y coordinate.
// points[i][2] - the first derivative
// points[i][3] - the second derivative
func CreateHermitPolinomial(points [][]float64, numDerivatives int) (*Hermit, error) {

	h := &Hermit{
		points:         make([][]float64, len(points)),
		numDerivatives: numDerivatives,
	}
	for i := 0; i < len(points); i++ {
		if len(points[i]) < 2+numDerivatives {
			return nil, ErrNotEnoughInputData
		}
		h.points[i] = make([]float64, 2+numDerivatives)
		copy(h.points[i], points[i][:2+numDerivatives])
	}

	return h, nil
}

// SetPoints() modifies the set of points from which the approximate value is calculated.
// points[i][0] - x coordinate.
// points[i][1] - y coordinate.
// points[i][2] - the first derivative
// points[i][3] - the second derivative
func (h *Hermit) SetPoints(points [][]float64, numDerivatives int) error {
	for i := 0; i < len(points); i++ {
		if len(points[i]) < 2+numDerivatives {
			return ErrNotEnoughInputData
		}
		h.points[i] = make([]float64, 2+numDerivatives)
		copy(h.points[i], points[i][:2+numDerivatives])
	}
	h.numDerivatives = numDerivatives

	return nil
}

// Calc() calculates the approximate value of y(x) for the degree of the polynomial n.
// x - the input value.
// n - the degree of the Newton polynomial.
func (h *Hermit) Calc(x float64, n int) (float64, error) {
	err := h.configure(x, n)
	if err != nil {
		return UndefNum, err
	}

	h.buildDiff()

	args := []float64{}
	for i := 1; i <= h.numOfNodes; i++ {
		args = append(args, h.differences[0][i])
	}

	var result, temp float64
	for i := 0; i < len(args); i++ {
		temp = args[i]
		for j := 0; j < i; j++ {
			temp *= (x - h.differences[j][0])
		}
		result += temp
	}

	return result + temp, nil
}

// configure() creates a configuration of the values of the starting points. n + 1 points are selected, as close as possible to x.
// x - the input value.
// n - the degree of the Newton polynomial.
func (h *Hermit) configure(x float64, n int) error {

	if len(h.points) <= n {
		return ErrNotEnoughInputData
	}

	h.numOfNodes = (n + 1) * (h.numDerivatives + 1)
	h.config = make([][]float64, h.numOfNodes)

	sort.SliceStable(h.points, func(i, j int) bool {
		return h.points[i][0] < h.points[j][0]
	})
	index, _ := slices.BinarySearchFunc(h.points, x, func(point []float64, pointX float64) int {
		if point[0] >= x {
			return 1
		}
		return -1
	})

	left, right := index-1, index
	var point []float64

	for count := 0; count < h.numOfNodes; {
		if left >= 0 && right < len(h.points) && x-h.points[left][0] < h.points[right][0]-x {
			point = h.points[left][:2+h.numDerivatives]
			left--
		} else if left >= 0 && right < len(h.points) {
			point = h.points[right][:2+h.numDerivatives]
			right++
		} else if right < len(h.points) {
			point = h.points[right][:2+h.numDerivatives]
			right++
		} else {
			point = h.points[left][:2+h.numDerivatives]
			left--
		}
		for i := 0; i < h.numDerivatives+1; i++ {
			h.config[count] = make([]float64, len(point))
			copy(h.config[count], point)
			count++
		}
	}

	sort.SliceStable(h.config, func(i, j int) bool {
		return h.config[i][0] < h.config[j][0]
	})

	return nil
}

// buildDiff() calculates the values of the divided differences.
// n - the degree of the Newton polynomial.
func (h *Hermit) buildDiff() {

	h.differences = make([][]float64, h.numOfNodes)
	for i := 0; i < h.numOfNodes; i++ {
		h.differences[i] = h.config[i][:2]
	}
	numOfNodes := h.numOfNodes
	n := h.numOfNodes - 1

	for k := 1; k <= n; k++ {
		idx := len(h.differences[0]) - 1

		for i := 0; i < numOfNodes-1; i++ {
			var diff float64
			if k <= h.numDerivatives && math.Abs(h.differences[i][0]-h.differences[i+k][0]) < delta {
				diff = h.config[i][k+1]
			} else {
				diff = (h.differences[i][idx] - h.differences[i+1][idx]) /
					(h.differences[i][0] - h.differences[i+k][0])
			}
			h.differences[i] = append(h.differences[i], diff)
		}
		numOfNodes--
	}
}

// SepFiffTable() prints a table of the divided differences of the last operation.
func (h *Hermit) PrintDiffTable() {
	k := len(h.differences[0])*16 + 1
	printLine(k)

	fmt.Printf("|       x       |       y       ")
	for i := 2; i < len(h.differences[0]); i++ {
		fmt.Printf("| y(x%-2d,..,x%-2d) ", 0, i-1)
	}
	fmt.Println("|")

	printLine(k)

	for i := 0; i < len(h.differences); i++ {
		for j := 0; j < len(h.differences[i]); j++ {
			fmt.Printf("| ")
			if h.differences[i][j] >= 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("%-12f ", h.differences[i][j])
			if h.differences[i][j] < 0 {
				fmt.Printf(" ")
			}
		}
		for j := 0; j < len(h.differences[0])-len(h.differences[i]); j++ {
			fmt.Printf("|               ")
		}
		fmt.Println("|")
	}

	printLine(k)
}
