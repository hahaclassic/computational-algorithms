package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"

	"github.com/hahaclassic/computational-algorithms.git/internal/interpolation"
	op "github.com/hahaclassic/computational-algorithms.git/internal/operations"
	"github.com/hahaclassic/computational-algorithms.git/internal/reader"
)

var (
	mainFile            string
	XYFile              string
	YXFile              string
	separator           rune = ','
	MainFieldsPerRecord int  = 4
	FieldsPerRecord     int  = 2
)

func init() {
	flag.StringVar(&mainFile, "main", "./data/source_data.csv", "the main data file")
	flag.StringVar(&XYFile, "xy", "./data/xy.csv", "the xy data file")
	flag.StringVar(&YXFile, "yx", "./data/yx.csv", "the yx data file")
	flag.Parse()

	if mainFile == "" || XYFile == "" || YXFile == "" {
		log.Fatal("Files' names are not specified")
	}
}

func main() {
	data, err := reader.ReadCSVFloatMatrix(mainFile, separator, MainFieldsPerRecord)
	if err != nil {
		log.Fatal(err)
	}
	dataXY, err := reader.ReadCSVFloatMatrix(XYFile, separator, FieldsPerRecord)
	if err != nil {
		log.Fatal(err)
	}
	dataYX, err := reader.ReadCSVFloatMatrix(YXFile, separator, FieldsPerRecord)
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
			err = op.SolveSystemOfEquations(dataXY, dataYX)
		case op.ChangeDegree:
			err = op.SetNumDerivatives(hermit)
		}
		if err != nil {
			slog.Error(err.Error())
		}
		operation = op.ChooseOperation()
	}
	fmt.Println("Программа завершена.")
}
