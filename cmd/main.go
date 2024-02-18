package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hahaclassic/computational-algorithms.git/internal/interpolation"
	"github.com/hahaclassic/computational-algorithms.git/internal/matrix"
	op "github.com/hahaclassic/computational-algorithms.git/internal/operations"
	"github.com/hahaclassic/computational-algorithms.git/internal/reader"
)

var (
	fileName        string
	separator       rune = ','
	fieldsPerRecord int  = 4
	//x               float64
	//n               int
)

func init() {
	flag.StringVar(&fileName, "filename", "", "the data file")
	// flag.Float64Var(&x, "x", 0, "x")
	// flag.IntVar(&n, "n", 0, "n")
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

	data, err := matrix.MatrixAtof(strData[1:])
	if err != nil {
		log.Fatal(err)
	}

	newton, err := interpolation.CreateNewtonPolinomial(data)
	if err != nil {
		log.Fatal(err)
	}
	hermit, err := interpolation.CreateHermitPolinomial(data, 2)
	if err != nil {
		log.Fatal(err)
	}

	operation := op.ChooseOperation()
	for operation != op.Exit {
		switch operation {
		case op.CalcNewton:
			err = op.CalcValueByNewton(newton)
		case op.ShowNewtonTable:
			newton.PrintDiffTable()
		case op.CalcHermit:
			err = op.CalcValueByHermit(hermit)
		case op.ShowHermitTable:
			hermit.PrintDiffTable()
		case op.FindRootValue:
			err = op.FindRoot(newton, hermit)
		case op.ComparePolynomials:
			err = op.Compare(newton, hermit)
		case op.SolveSystem:
			//
		case op.ChangeDegree:
			err = op.SetNumDerivatives(hermit)
		}
		if err != nil {
			fmt.Println(err)
		}
		operation = op.ChooseOperation()
	}
}
