package interpolation

import (
	"slices"
	"sort"
)

type Spline struct {
	points       [][]float64
	config       [][4]float64
	isConfigured bool
}

// CreateSpline()
// points[i][0] - x coordinate.
// points[i][1] - y coordinate.
func CreateSpline(points [][]float64) (*Spline, error) {
	spline := &Spline{
		points: make([][]float64, len(points)),
	}
	for i := 0; i < len(points); i++ {
		if len(points[i]) < 2 {
			return nil, ErrNotEnoughInputData
		}
		spline.points[i] = make([]float64, 2)
		copy(spline.points[i], points[i][:2])
	}
	sort.SliceStable(spline.points, func(i, j int) bool {
		return spline.points[i][0] < spline.points[j][0]
	})
	// Config starts from 1 index
	spline.config = make([][4]float64, len(points))

	return spline, nil
}

// startC, endC - C1 and Ci+1
func (s *Spline) SetBoundaryCond(startC, endC float64) {

	// Метод прогонки, вычисляет коэффициенты Ci
	s.shuttle(startC, endC)

	// Вычисляет значение коэффициентов a, b, d
	s.configure()

	s.isConfigured = true
}

func (s *Spline) Calc(x float64) float64 {
	if !s.isConfigured {
		s.SetBoundaryCond(0, 0)
		s.isConfigured = true
	}

	index, _ := slices.BinarySearchFunc(s.points, x, func(point []float64, pointX float64) int {
		if point[0] >= x {
			return 1
		}
		return -1
	})
	if index == 0 {
		index++
	} else if index == len(s.points) {
		index--
	}

	var result float64
	var diff float64 = 1
	for i := 0; i < 4; i++ {
		result += s.config[index][i] * diff
		diff *= (x - s.points[index-1][0])
	}

	return result
}

// Hi = Xi - Xi-1
func (s *Spline) paramH(idx int) float64 {
	return s.points[idx][0] - s.points[idx-1][0]
}

func (s *Spline) paramF(idx int, h1, h2 float64) float64 {

	a := (s.points[idx][1] - s.points[idx-1][1]) / h2
	b := (s.points[idx-1][1] - s.points[idx-2][1]) / h1

	return 3 * (a - b)
}

func (s *Spline) shuttle(start, end float64) {

	ksiValues := make([]float64, len(s.points))
	tetaValues := make([]float64, len(s.points))
	ksiValues[1], tetaValues[1] = start, end

	for i := 2; i < len(s.points); i++ {

		h1 := s.paramH(i - 1)
		h2 := s.paramH(i)
		fi := s.paramF(i, h1, h2)

		ksiValues[i] = -h1 / (h2*ksiValues[i-1] + 2*(h2+h1))
		tetaValues[i] = (fi - h1*tetaValues[i-1]) / (h1*ksiValues[i-1] + 2*(h2+h1))
	}

	// set c parameter
	s.config[1][2] = start
	s.config[len(s.config)-1][2] = end //tetaValues[len(s.config)-1]
	for i := len(s.config) - 1; i > 1; i-- {
		s.config[i-1][2] = ksiValues[i-1]*s.config[i][2] + tetaValues[i-1]
	}
}

func (s *Spline) calcB(idx int, h float64) float64 {
	return (s.points[idx][1]-s.points[idx-1][1])/h -
		h*(s.config[idx+1][2]+2*s.config[idx][2])/3
}

func (s *Spline) calcD(idx int, h float64) float64 {
	return (s.config[idx+1][2] - s.config[idx][2]) / (3.0 * h)
}

func (s *Spline) configure() {
	// idx := 2
	for i := 1; i < len(s.config)-1; i++ {
		h := s.paramH(i)
		// set Ai
		s.config[i][0] = s.points[i-1][1]
		// set Bi
		s.config[i][1] = s.calcB(i, h)
		// set Di
		s.config[i][3] = s.calcD(i, h)
	}

	l := len(s.config) - 1
	h := s.paramH(l)
	// set An
	s.config[l][0] = s.points[l-1][1]
	// set Bn
	s.config[l][1] = (s.points[l][1]-s.points[l-1][1])/h - (2.0 / 3.0 * h * s.config[l][2])
	// set Dn
	s.config[l][3] = -s.config[l][2] / (3.0 * h)
}
