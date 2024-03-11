package operations

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/hahaclassic/computational-algorithms.git/internal/format"
	"github.com/hahaclassic/computational-algorithms.git/internal/interpolation"
)

func menu() {
	fmt.Println(header)
	for i := CalcValue; i <= ComparePolynomials; i++ {
		fmt.Printf("| %d. %-82s |\n", int(i), i)
	}
	fmt.Print(emptyLine)
	fmt.Printf("| 0. %-82s |\n", Exit)
	fmt.Print(line)
}

func ChooseOperation() Operation {
	menu()
	var (
		num int
		err error
	)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Введите номер операции: ")
		scanner.Scan()
		num, err = strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("[ERR]: Неверный номер операции. Введите номер повторно.")
			continue
		}
		if num >= int(Exit) && num <= int(ComparePolynomials) {
			break
		}
		fmt.Println("[ERR]: Неверный номер операции. Введите номер повторно.")
	}

	return Operation(num)
}

func CalcValues(newton *interpolation.Newton, spline *interpolation.Spline) error {

	x, err := format.ReadValue()
	if err != nil {
		return err
	}

	newtonResult, err := newton.Calc(x, 3)
	if err != nil {
		return err
	}

	splineResult := spline.Calc(x)

	format.PrintNewtonResult(newtonResult)
	format.PrintSplineResult(splineResult)
	return nil
}

func SetNaturalCond(spline *interpolation.Spline) {
	spline.SetBoundaryCond(0, 0)
}

func SetStart(startX float64, newton *interpolation.Newton, spline *interpolation.Spline) error {
	p, err := newton.Derivative2(startX, 3)
	if err != nil {
		return err
	}

	spline.SetBoundaryCond(p/2, 0)

	return nil
}

func SetStartEnd(startX, endX float64, newton *interpolation.Newton, spline *interpolation.Spline) error {
	p1, err := newton.Derivative2(startX, 3)
	if err != nil {
		return err
	}
	p2, err := newton.Derivative2(endX, 3)
	if err != nil {
		return err
	}

	spline.SetBoundaryCond(p1/2, p2/2)
	return nil
}

func addPoints(x []float64, start, end float64) []float64 {
	diff := (end - start) / 4
	for i := 0; i < 3; i++ {
		start += diff
		x = append(x, start)
	}
	return x
}

func Compare(points [][]float64, newton *interpolation.Newton, spline *interpolation.Spline) error {

	x := []float64{}
	x = addPoints(x, points[0][0], points[1][0])
	x = addPoints(x, points[len(points)/2][0], points[len(points)/2+1][0])
	x = addPoints(x, points[len(points)-2][0], points[len(points)-1][0])

	newtonResults, splineResults := make([]float64, len(x)), make([]float64, len(x))
	for i := 0; i < len(x); i++ {
		res, err := newton.Calc(x[i], 3)
		if err != nil {
			return err
		}
		newtonResults[i] = res

		res = spline.Calc(x[i])
		splineResults[i] = res
	}

	fmt.Println()
	format.PrintLine(34*5 - 3)
	fmt.Print("\n|     x      ")
	for i := 0; i < 9; i++ {
		fmt.Printf("| %-14f ", x[i])
	}
	fmt.Print("|\n|")
	format.PrintLine(34*5 - 5)
	//fmt.Println("|\n|---------------------------------------------------------------------------------------------------------------------------------------------------------------------|")
	fmt.Print("|\n| Newton     ")
	for i := 0; i < 9; i++ {
		fmt.Printf("| %-14f ", newtonResults[i])
	}
	fmt.Print("|\n|")
	format.PrintLine(34*5 - 5)
	//fmt.Println("|\n|---------------------------------------------------------------------------------------------------------------------------------------------------------------------|")
	fmt.Print("|\n| Spline     ")
	for i := 0; i < 9; i++ {
		fmt.Printf("| %-14f ", splineResults[i])
	}
	fmt.Println("|")
	format.PrintLine(34*5 - 3)
	fmt.Print("\n\n")
	return nil
}
