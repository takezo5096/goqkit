package dataset

import (
	"encoding/csv"
	"os"
	"strconv"
)

func IrisDataset(filePath string, limit int) ([][]float64, []float64, error) {

	nameMap := map[string]float64{"Iris-setosa": 0, "Iris-versicolor": 1, "Iris-virginica": 2}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	record, err := r.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	n := len(record)
	if limit > 0 {
		n = limit
	}

	//double array for iris
	xData := MakeDoubleArray(n, 4)
	yData := make([]float64, n)

	for i := 0; i < n; i++ {
		xData[i] = make([]float64, 4)
		for j := 0; j < 4; j++ {
			xData[i][j], _ = strconv.ParseFloat(record[i][j], 64)
		}
		yData[i] = nameMap[record[i][4]]
	}

	return xData, yData, nil
}
