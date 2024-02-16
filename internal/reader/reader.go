package reader

import (
	"encoding/csv"
	"log"
	"os"
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
		log.Println("[ERR]: Can't read the file")
		return nil, err
	}
	if len(data) == 0 {
		log.Println("[ERR]: File is empty")
		return nil, err
	}

	return data, nil
}
