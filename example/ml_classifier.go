package example

import (
	"fmt"
	"github.com/takezo5096/goqkit/ml"
	"github.com/takezo5096/goqkit/ml/dataset"
	"github.com/takezo5096/goqkit/ml/optimizer"
	"log"
	"math"
)

func MLClassifier() {

	//should change the path for iris.data
	X, Y, err := dataset.IrisDataset("./test/iris.data", -1)
	if err != nil {
		log.Fatal(err)
	}
	X = dataset.NormalizeMinMax(X, 0, 2*math.Pi)
	trainXData, trainYData, testXData, testYData := dataset.MakeTrainTestDataOneHot(X, Y, 0.20)

	nClasses := 3
	nLayers := 5
	nQBits := 4
	clf := ml.Classifier{NumberOfClasses: nClasses, NumberOfLayers: nLayers, NumberOfQBits: nQBits}
	clf.SetTrainingData(trainXData, trainYData)
	clf.SetTrainingStatusHandler(func(epoch int, loss float64, acc float64, a int, t int) {
		fmt.Printf("%d loss:%f train acc:%f (%d/%d)\n", epoch, loss, acc, a, t)
	})
	var opti optimizer.Optimizer = optimizer.NewAdam(nLayers, nQBits, 0.001)
	clf.Train(opti, 50)
	acc, a, t := clf.Accuracy(testXData, testYData)
	fmt.Printf("test acc:%f (%d/%d)\n", acc, a, t)

	//to predict
	//pred := clf.Predict(testXData)

	//to get parameters of classifier
	//theta := clf.GetThetaParameters()
}
