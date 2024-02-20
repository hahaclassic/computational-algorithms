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
	for i := CalcNewton; i <= ChangeDegree; i++ {
		fmt.Printf("| %d. %-82s |\n", int(i), i)
		if i == ShowNewtonTable || i == ShowHermitTable ||
			i == SolveSystem || i == ChangeDegree {
			fmt.Print(emptyLine)
		}
	}
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
		if num >= int(Exit) && num <= int(ChangeDegree) {
			break
		}
		fmt.Println("[ERR]: Неверный номер операции. Введите номер повторно.")
	}

	return Operation(num)
}

func CalcValueByNewton(newton *interpolation.Newton) error {

	x, err := format.ReadValue()
	if err != nil {
		return err
	}
	n, err := format.ReadPolynomialDegree()
	if err != nil {
		return err
	}

	result, err := newton.Calc(x, n)
	if err != nil {
		return err
	}

	format.PrintNewtonResult(result)
	return nil
}

func CalcValueByHermit(hermit *interpolation.Hermit) error {
	x, err := format.ReadValue()
	if err != nil {
		return err
	}
	n, err := format.ReadPolynomialDegree()
	if err != nil {
		return err
	}
	result, err := hermit.Calc(x, n)
	if err != nil {
		return err
	}

	format.PrintHermitResult(result)
	return nil
}

func FindRoot(newton *interpolation.Newton, hermit *interpolation.Hermit) error {
	n, err := format.ReadPolynomialDegree()
	if err != nil {
		return err
	}
	result, err := newton.FindRoot(n)
	if err != nil {
		return err
	}
	format.PrintNewtonResult(result)

	result, err = hermit.FindRoot(n)
	if err != nil {
		return err
	}
	format.PrintHermitResult(result)
	return nil
}

func SetNumDerivatives(hermit *interpolation.Hermit) error {
	n, err := format.ReadNumOfDerivates()
	if err != nil {
		return err
	}
	hermit.SetNumDerivatives(n)
	return nil
}

func Compare(newton *interpolation.Newton, hermit *interpolation.Hermit) error {
	x, err := format.ReadValue()
	if err != nil {
		return err
	}
	degree := 5
	newtonResults, hermitResults := make([]float64, degree), make([]float64, degree)
	for i := 0; i < degree; i++ {
		res, err := newton.Calc(x, i+1)
		if err != nil {
			return err
		}
		newtonResults[i] = res

		res, err = hermit.Calc(x, i+1)
		if err != nil {
			return err
		}
		hermitResults[i] = res
	}

	format.PrintLine(20*5 - 1)
	fmt.Print("| Polynomial ")
	for i := 0; i < degree; i++ {
		fmt.Printf("|        %d       ", i+1)
	}
	fmt.Println("|\n|-------------------------------------------------------------------------------------------------|")
	fmt.Print("| Newton     ")
	for i := 0; i < degree; i++ {
		fmt.Printf("| %-14f ", newtonResults[i])
	}
	fmt.Println("|\n|-------------------------------------------------------------------------------------------------|")
	fmt.Print("| Hermit     ")
	for i := 0; i < degree; i++ {
		fmt.Printf("| %-14f ", hermitResults[i])
	}
	fmt.Println("|")
	format.PrintLine(20*5 - 1)
	fmt.Println()
	return nil
}

func SolveSystemOfEquations(pointsXY, pointsYX [][]float64) error {
	n, err := format.ReadPolynomialDegree()
	if err != nil {
		return err
	}
	inverted, _ := interpolation.Inverse(pointsYX)

	newtonXY, err := interpolation.CreateNewtonPolinomial(pointsXY)
	if err != nil {
		return err
	}

	newtonYX, err := interpolation.CreateNewtonPolinomial(inverted)
	if err != nil {
		return err
	}

	diff := make([][]float64, len(pointsXY))
	for i := 0; i < len(pointsXY); i++ {
		y, err := newtonYX.Calc(pointsXY[i][0], n)
		if err != nil {
			return err
		}
		diff[i] = append(diff[i], pointsXY[i][0], pointsXY[i][1]-y)
	}

	newtonXY.SetPoints(diff)
	root, err := newtonXY.FindRoot(n)
	if err != nil {
		return err
	}
	y, err := newtonYX.Calc(root, n)
	format.PrintSystemResult(root, y)
	return nil
}
