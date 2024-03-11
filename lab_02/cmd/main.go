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
	mainFile        string
	separator       rune = ','
	FieldsPerRecord int  = 2
)

func init() {
	flag.StringVar(&mainFile, "data", "./data/src.csv", "the main data file")
	flag.Parse()

	if mainFile == "" {
		log.Fatal("File's name are not specified")
	}
}

func main() {
	data, err := reader.ReadCSVFloatMatrix(mainFile, separator, FieldsPerRecord)
	if err != nil {
		log.Fatal(err)
	}

	newton, err := interpolation.CreateNewtonPolinomial(data)
	if err != nil {
		log.Fatal(err)
	}
	spline, err := interpolation.CreateSpline(data)
	if err != nil {
		log.Fatal(err)
	}

	operation := op.ChooseOperation()
	for operation != op.Exit {
		switch operation {
		case op.CalcValue:
			err = op.CalcValues(newton, spline)
		case op.SetupNaturalCond:
			op.SetNaturalCond(spline)
		case op.SetupStart:
			err = op.SetStart(data[0][0], newton, spline)
		case op.SetupStartEnd:
			err = op.SetStartEnd(data[0][0], data[len(data)-1][0], newton, spline)
		case op.ComparePolynomials:
			err = op.Compare(data, newton, spline)
		}
		if err != nil {
			slog.Error(err.Error())
		}
		operation = op.ChooseOperation()
	}
	fmt.Println("Программа завершена.")
}
