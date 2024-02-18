package format

import "fmt"

func ReadPolynomialDegree() (int, error) {
	var n int
	fmt.Printf("Введите степень полинома (целое): ")
	_, err := fmt.Scan(&n)
	return n, err
}

func ReadValue() (float64, error) {
	var x float64
	fmt.Print("Введите значение x (вещественное): ")
	_, err := fmt.Scan(&x)
	return x, err
}

func ReadNumOfDerivates() (int, error) {
	var n int
	fmt.Print("Введите количество используемых производных: ")
	_, err := fmt.Scan(&n)
	return n, err
}

func PrintHermitResult(res float64) {
	fmt.Println("\nРезультат, полученный с помощью полинома Эрмита H(x):", res)
}

func PrintNewtonResult(res float64) {
	fmt.Println("\nРезультат, полученный с помощью полинома Ньютона P(x):", res)
}

func PrintLine(k int) {
	for i := 0; i < k; i++ {
		fmt.Print("-")
	}
	fmt.Println()
}
