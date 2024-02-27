package reader

import (
	"encoding/csv"
	"os"

	"github.com/hahaclassic/computational-algorithms.git/internal/matrix"
)

func ReadCSV(fileName string, separator rune, fieldsPerRecord int) ([][]string, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = separator
	reader.FieldsPerRecord = fieldsPerRecord

	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, err
	}

	return data, nil
}

func ReadCSVFloatMatrix(fileName string, separator rune, fieldsPerRecord int) ([][]float64, error) {
	strData, err := ReadCSV(fileName, separator, fieldsPerRecord)
	if err != nil {
		return nil, err
	}
	data, err := matrix.MatrixAtof(strData[1:])
	if err != nil {
		return nil, err
	}
	return data, nil
}
