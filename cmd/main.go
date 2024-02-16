package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hahaclassic/computational-algorithms.git/internal/interpolation"
	"github.com/hahaclassic/computational-algorithms.git/internal/matrix"
	"github.com/hahaclassic/computational-algorithms.git/internal/reader"
)

var (
	fileName        string
	separator       rune = ','
	fieldsPerRecord int  = 4
	x               float64
	n               int
)

func init() {
	flag.StringVar(&fileName, "filename", "", "the data file")
	flag.Float64Var(&x, "x", 0, "x")
	flag.IntVar(&n, "n", 0, "n")
	flag.Parse()

	if fileName == "" {
		fmt.Println(fileName)
		log.Fatal("File name is not specified")
	}
}

func main() {

	strData, err := reader.ReadCSV(fileName, separator, fieldsPerRecord)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(strData)

	data, err := matrix.MatrixAtof(strData[1:])
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(data)

	newton := interpolation.CreateNewtonPolinomial(data)

	res, err := newton.Calc(x, n)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	fmt.Println(newton.SepDiffTable())
}
