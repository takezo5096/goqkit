package dataset

import (
	"math"
	"math/rand"
	"time"
)

func MakeDoubleArray(rows, cols int) [][]float64 {
	data := make([][]float64, rows)

	for i := 0; i < rows; i++ {
		data[i] = make([]float64, cols)
	}
	return data
}

func ShuffleArray(data []float64) []float64 {
	rows := len(data)
	newData := make([]float64, rows)

	rand.Seed(time.Now().UnixNano())
	for i := rows - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		newData[i], newData[j] = data[j], data[i]
	}
	return newData
}

func ShuffleDoubleArray(data [][]float64) [][]float64 {
	rows := len(data)
	cols := len(data[0])
	newData := MakeDoubleArray(rows, cols)

	rand.Seed(time.Now().UnixNano())
	for i := rows - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		newData[i], newData[j] = data[j], data[i]
	}
	return newData
}

func CopyArray(data []float64) []float64 {
	rows := len(data)

	newData := make([]float64, rows)

	for i := 0; i < rows; i++ {
		newData[i] = data[i]
	}
	return newData
}

func CopyDoubleArray(data [][]float64) [][]float64 {
	rows := len(data)
	cols := len(data[0])

	newData := make([][]float64, rows)

	for i := 0; i < rows; i++ {
		newData[i] = make([]float64, cols)
		newData[i] = data[i]
	}
	return newData
}

func ShuffleData(xData [][]float64, yData []float64) ([][]float64, []float64) {
	rows := len(xData)

	newXData := CopyDoubleArray(xData)
	newYData := CopyArray(yData)

	rand.Seed(time.Now().UnixNano())
	for i := rows - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		newXData[i], newXData[j] = newXData[j], newXData[i]
		newYData[i], newYData[j] = newYData[j], newYData[i]
	}
	return newXData, newYData
}

/*
func MakeTrainTestData(xData [][]float64, yData []float64, testDataRatio float64) ([][]float64, []float64, [][]float64, []float64) {

	sXData, sYData := ShuffleData(xData, yData)

	n := len(xData)

	testDataLimit := int(float64(n) * (1.0-testDataRatio))

	trainXData := sXData[:testDataLimit]
	testXData := sXData[testDataLimit:]
	trainYData := sYData[:testDataLimit]
	testYData := sYData[testDataLimit:]

	return trainXData, trainYData, testXData, testYData
}
*/
func SliceByClass(xData [][]float64, yData []float64) ([][][]float64, [][]float64) {
	classesMap := make(map[int]int)
	cnt := 0
	for i := 0; i < len(yData); i++ {
		if _, ok := classesMap[int(yData[i])]; !ok {
			classesMap[int(yData[i])] = cnt
			cnt++
		}
	}

	xDataClass := make([][][]float64, len(classesMap))
	yDataClass := make([][]float64, len(classesMap))
	cnt = 0
	for classKey, idx := range classesMap {
		xDataClass[idx] = make([][]float64, 0)
		yDataClass[idx] = make([]float64, 0)
		for i := 0; i < len(yData); i++ {
			if int(yData[i]) == classKey {
				xDataClass[idx] = append(xDataClass[idx], xData[i])
				yDataClass[idx] = append(yDataClass[idx], yData[i])
			}
		}
		cnt++
	}
	return xDataClass, yDataClass
}

func MakeTrainTestDataOneHot(xData [][]float64, yData []float64, testDataRatio float64) ([][]float64, [][]float64, [][]float64, [][]float64) {
	trainXData, trainYData, testXData, testYData := MakeTrainTestData(xData, yData, testDataRatio)
	trainYOnehotData := OneHot(trainYData)
	testYOnehotData := OneHot(testYData)
	return trainXData, trainYOnehotData, testXData, testYOnehotData
}
func MakeTrainTestData(xData [][]float64, yData []float64, testDataRatio float64) ([][]float64, []float64, [][]float64, []float64) {

	xDataC, yDataC := SliceByClass(xData, yData)

	trainXData := make([][]float64, 0)
	trainYData := make([]float64, 0)
	testXData := make([][]float64, 0)
	testYData := make([]float64, 0)

	for i := 0; i < len(xDataC); i++ {

		testDataLimit := int(float64(len(xDataC[i])) * (1.0 - testDataRatio))

		trainXD := xDataC[i][:testDataLimit]
		trainYD := yDataC[i][:testDataLimit]
		testXD := xDataC[i][testDataLimit:]
		testYD := yDataC[i][testDataLimit:]

		trainXData = append(trainXData, trainXD...)
		trainYData = append(trainYData, trainYD...)

		testXData = append(testXData, testXD...)
		testYData = append(testYData, testYD...)

	}
	trainXData, trainYData = ShuffleData(trainXData, trainYData)
	testXData, testYData = ShuffleData(testXData, testYData)

	return trainXData, trainYData, testXData, testYData
}

func Mean(data []float64) float64 {
	s := 0.0
	for i := 0; i < len(data); i++ {
		s += data[i]
	}
	return s / float64(len(data))
}

func Max(data []float64) (float64, int) {
	max := 0.0
	idx := -1
	for i := 0; i < len(data); i++ {
		if data[i] > max {
			max = data[i]
			idx = i
		}
	}
	return max, idx
}
func Min(data []float64) float64 {
	min := 999999999.0
	for i := 0; i < len(data); i++ {
		if data[i] < min {
			min = data[i]
		}
	}
	return min
}

func Sigma(data []float64) float64 {
	m := Mean(data)
	s := 0.0
	for i := 0; i < len(data); i++ {
		s += math.Pow(data[i]-m, 2)
	}
	return s / float64(len(data))
}
func NormalizeMinMax(data [][]float64, rangeMin float64, rangeMax float64) [][]float64 {
	norm := make([][]float64, len(data))

	min := 99999999.0
	max := 0.0
	for i := 0; i < len(data); i++ {
		minTmp := Min(data[i])
		maxTmp, _ := Max(data[i])
		if minTmp < min {
			min = minTmp
		}
		if maxTmp > max {
			max = maxTmp
		}
	}

	for i := 0; i < len(data); i++ {
		norm[i] = make([]float64, len(data[i]))
		for j := 0; j < len(data[i]); j++ {
			norm[i][j] = (data[i][j]-min)/(max-min)*(rangeMax-rangeMin) + rangeMin
		}
	}
	return norm
}

func Normalize(data []float64) []float64 {
	m := Mean(data)
	s := math.Sqrt(Sigma(data))
	norm := make([]float64, len(data))
	for i := 0; i < len(data); i++ {
		norm[i] = (data[i] - m) / s
	}
	return norm
}

func Sum(data []float64) float64 {
	s := 0.0
	for i := 0; i < len(data); i++ {
		s += data[i]
	}
	return s
}

func MeanMat(data [][]float64) float64 {
	s := 0.0
	cnt := 0
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			s += data[i][j]
			cnt++
		}
	}
	return s / float64(cnt)
}

func SigmaMat(data [][]float64) float64 {
	m := MeanMat(data)
	s := 0.0
	cnt := 0
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			s += math.Pow(data[i][j]-m, 2)
			cnt++
		}
	}
	return s / float64(cnt)
}

func NormalizeMat(data [][]float64) [][]float64 {
	m := MeanMat(data)
	s := math.Sqrt(SigmaMat(data))
	norm := make([][]float64, len(data))
	for i := 0; i < len(data); i++ {
		norm[i] = make([]float64, len(data[i]))
		for j := 0; j < len(data[i]); j++ {
			norm[i][j] = (data[i][j] - m) / s
		}
	}
	return norm
}

func OneHot(data []float64) [][]float64 {

	classesMap := map[int]int{}
	for i := 0; i < len(data); i++ {
		classesMap[int(data[i])] = 1
	}

	classes := make([]int, len(classesMap))
	for k, _ := range classesMap {
		classes[k] = k
	}

	classesMap2 := map[int]int{}
	for i := 0; i < len(classes); i++ {
		classesMap2[classes[i]] = i
	}

	onehots := make([][]float64, len(data))
	for i := 0; i < len(data); i++ {
		onehots[i] = make([]float64, len(classes))
		for j := 0; j < len(classes); j++ {
			if classesMap2[int(data[i])] == j {
				onehots[i][j] = 1.0
			} else {
				onehots[i][j] = 0.0
			}
		}
	}
	return onehots
}
