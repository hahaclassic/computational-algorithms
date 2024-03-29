package operations

type Operation int

const (
	Exit Operation = iota
	CalcNewton
	ShowNewtonTable
	CalcHermit
	ShowHermitTable
	FindRootValue
	ComparePolynomials
	SolveSystem
	ChangeDegree
)

func (op Operation) String() string {
	return []string{
		"Выход из программы.",
		"Вычислить значение функции y(x) c помощью полинома Ньютона.",
		"Вывести таблицу разделенных разностей для полинома Ньютона.",
		"Вычислить значение функции y(x) c помощью полинома Эрмитa.",
		"Вывести таблицу разделенных разностей для полинома Эрмита.",
		"Найти корень функции (y == 0).",
		"Сравнение работы полиномов при степенях n от 1 до 5.",
		"Решить систему уравнений.",
		"Изменить количество используемых степеней в полиноме Эрмита.",
	}[op]
}

const header string = `=========================================================================================
|                  Интерполяция с помощью полиномов Ньютона и Эрмита                    |
-----------------------------------------------------------------------------------------
| Операции                                                                              |
|---------------------------------------------------------------------------------------|`

const line string = "-----------------------------------------------------------------------------------------\n"
const emptyLine string = "|                                                                                       |\n"
