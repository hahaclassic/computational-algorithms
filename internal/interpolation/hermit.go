package interpolation

import (
	"fmt"
	"math"
	"slices"
	"sort"

	"github.com/hahaclassic/computational-algorithms.git/internal/format"
)

type Hermit struct {
	points         [][]float64
	config         [][]float64
	differences    [][]float64
	numOfNodes     int
	numDerivatives int // number of derivatives at a point
}

// CreateHermitPolinomial() creates a Hermit structure that implements interpolation using the Hermit polynomial.
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
	sort.SliceStable(h.points, func(i, j int) bool {
		return h.points[i][0] < h.points[j][0]
	})

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
	sort.SliceStable(h.points, func(i, j int) bool {
		return h.points[i][0] < h.points[j][0]
	})
	return nil
}

// SetPoints() modifies the num of derivatives.
func (h *Hermit) SetNumDerivatives(numOfDerivates int) error {
	if len(h.points[0]) < 2+numOfDerivates {
		return ErrInvalidNumDerivates
	}
	h.numDerivatives = numOfDerivates
	return nil
}

// Calc() calculates the approximate value of y(x) for the degree of the polynomial n.
// x - the input value.
// n - the degree of the Hermit polynomial.
func (h *Hermit) Calc(x float64, n int) (float64, error) {
	if n < 0 {
		return UndefNum, ErrInvalidPolynomialDegree
	}
	err := h.configure(x, n)
	if err != nil {
		return UndefNum, err
	}
	h.buildDiff()

	return h.result(x), nil
}

// FindRoot() finds root of the function (y == 0)
// n - the degree of the Hermit polynomial.
func (h *Hermit) FindRoot(n int) (float64, error) {
	if n < 0 {
		return UndefNum, ErrInvalidPolynomialDegree
	}

	idx := -1
	for i := 0; i < len(h.points)-1; i++ {
		if math.Abs(h.points[i][1]) < delta {
			return h.points[i][0], nil
		}
		if h.points[i][1]*h.points[i+1][1] < 0 {
			idx = i + 1
			break
		}
	}
	if idx == -1 {
		return UndefNum, ErrNoRoot
	}

	source := h.points
	inverted, err := Inverse(h.points)
	if err != nil {
		return UndefNum, err
	}
	h.points = inverted

	h.fillConfig(0, n, idx)
	h.buildDiff()
	h.points = source

	return h.result(0), nil
}

// configure() creates a configuration of the values of the starting points. n + 1 points are selected, as close as possible to x.
// x - the input value.
// n - the degree of the Hermit polynomial.
func (h *Hermit) configure(x float64, n int) error {

	if len(h.points) <= n {
		return ErrNotEnoughInputData
	}

	index, _ := slices.BinarySearchFunc(h.points, x, func(point []float64, pointX float64) int {
		if point[0] >= x {
			return 1
		}
		return -1
	})

	h.fillConfig(x, n, index)

	return nil
}

// n - degree of the polynomial
func (h *Hermit) fillConfig(x float64, n int, idx int) {
	h.numOfNodes = (n + 1) * (h.numDerivatives + 1)
	left, right := idx-1, idx
	leftNodes := [][]float64{}
	rightNodes := [][]float64{}

	for count := 0; count < h.numOfNodes; {
		if left >= 0 && right < len(h.points) &&
			math.Abs(x-h.points[left][0]) < math.Abs(h.points[right][0]-x) {
			addCopies(&leftNodes, h.points[left][:2+h.numDerivatives], h.numDerivatives+1)
			left--
		} else if left >= 0 && right < len(h.points) {
			addCopies(&rightNodes, h.points[right][:2+h.numDerivatives], h.numDerivatives+1)
			right++
		} else if right < len(h.points) {
			addCopies(&rightNodes, h.points[right][:2+h.numDerivatives], h.numDerivatives+1)
			right++
		} else {
			addCopies(&leftNodes, h.points[left][:2+h.numDerivatives], h.numDerivatives+1)
			left--
		}
		count += h.numDerivatives + 1
	}
	slices.Reverse(leftNodes)
	h.config = leftNodes
	h.config = append(h.config, rightNodes...)
}

func addCopies(dst *[][]float64, src []float64, n int) {
	for i := 0; i < n; i++ {
		c := make([]float64, len(src))
		copy(c, src)
		(*dst) = append(*dst, c)
	}
}

// buildDiff() calculates the values of the divided differences.
// n - the degree of the Hermit polynomial.
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

func (h *Hermit) result(x float64) float64 {
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

	return result
}

// SepFiffTable() prints a table of the divided differences of the last operation.
func (h *Hermit) PrintDiffTable() {

	if len(h.differences) == 0 {
		return
	}

	k := len(h.differences[0])*18 + 1
	fmt.Println()
	format.PrintLine(k)

	fmt.Printf("|        x        |        y        ")
	for i := 2; i < len(h.differences[0]); i++ {
		fmt.Printf("|  y(x%-2d,..,x%-2d)  ", 0, i-1)
	}
	fmt.Println("|")

	format.PrintLine(k)

	for i := 0; i < len(h.differences); i++ {
		for j := 0; j < len(h.differences[i]); j++ {
			fmt.Printf("| ")
			if h.differences[i][j] >= 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("%-14f ", h.differences[i][j])
			if h.differences[i][j] < 0 {
				fmt.Printf(" ")
			}
		}
		for j := 0; j < len(h.differences[0])-len(h.differences[i]); j++ {
			fmt.Printf("|                 ")
		}
		fmt.Println("|")
	}

	format.PrintLine(k)
	fmt.Println()
}
