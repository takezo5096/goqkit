package ml

import (
	"github.com/takezo5096/goqkit"
	"github.com/takezo5096/goqkit/ml/dataset"
	"github.com/takezo5096/goqkit/ml/optimizer"
	"math"
)

type TrainingStatusHandler func(int, float64, float64, int, int)

type Classifier struct {
	NumberOfQBits   int
	NumberOfLayers  int
	NumberOfClasses int

	trainXData [][]float64
	trainYData [][]float64

	theta [][]float64

	trainingStatusHandler TrainingStatusHandler
}

func (c *Classifier) SetTrainingData(X [][]float64, Y [][]float64) {
	c.trainXData = X
	c.trainYData = Y
}

func (c *Classifier) SetTrainingStatusHandler(handler TrainingStatusHandler) {
	c.trainingStatusHandler = handler
}

func (c *Classifier) Train(opti optimizer.Optimizer, epoch int) {

	c.theta = make([][]float64, c.NumberOfLayers)

	for i := 0; i < len(c.theta); i++ {
		c.theta[i] = make([]float64, c.NumberOfQBits)
		for j := 0; j < len(c.theta[i]); j++ {
			c.theta[i][j] = math.Pi
		}
	}

	lossList := []float64{}

	for ep := 1; ep <= epoch; ep++ {

		lossTmp := []float64{}

		for i := 0; i < len(c.trainXData); i++ {
			prediction := c.quantiumNN(c.trainXData[i], c.theta)

			lossVal := c.crossEntropyLoss(prediction, c.trainYData[i])
			lossTmp = append(lossTmp, lossVal)

			grad := c.gradient(c.trainXData[i], c.trainYData[i])
			delta := opti.Gradient(grad)
			for j := 0; j < len(c.theta); j++ {
				for k := 0; k < len(c.theta[j]); k++ {
					//theta[j][k] = theta[j][k] - eta * grad[j][k]
					c.theta[j][k] = c.theta[j][k] - delta[j][k]
				}
			}
		}
		lossList = append(lossList, dataset.Mean(lossTmp))

		acc, a, t := c.Accuracy(c.trainXData, c.trainYData)

		c.trainingStatusHandler(ep, lossList[len(lossList)-1], acc, a, t)
	}
}

func (c *Classifier) Predict(X [][]float64) [][]float64 {
	preds := make([][]float64, 0)

	for i := 0; i < len(X); i++ {
		pred := c.quantiumNN(X[i], c.theta)
		preds = append(preds, pred)
	}
	return preds
}

func (c *Classifier) featureMap(X []float64) (*goqkit.QBitsCircuit, *goqkit.Register) {
	circuit := goqkit.MakeQBitsCircuit(c.NumberOfQBits)
	register := circuit.AssignQBits(c.NumberOfQBits, "register")

	for i := 0; i < c.NumberOfQBits; i++ {
		circuit.Had(1<<i, 0)
	}

	for i, v := range X {
		deg := v * 180 / math.Pi
		circuit.RotY(1<<i, 0, deg)
	}

	return &circuit, register
}

func (c *Classifier) valiationalCircut(circuit *goqkit.QBitsCircuit, theta [][]float64) {

	N := int(circuit.QBitNumber)

	for j := 0; j < N-1; j++ {
		circuit.Not(1<<(j+1), 1<<j)
	}
	circuit.Not(0x01, 1<<(N-1))

	for i := 0; i < len(theta); i++ {
		for j := 0; j < len(theta[i]); j++ {
			degTheta := theta[i][j] * 180 / math.Pi
			circuit.RotY(1<<j, 0, degTheta)
		}
	}
}

func (c *Classifier) quantiumNN(X []float64, theta [][]float64) []float64 {

	circuit, _ := c.featureMap(X)
	c.valiationalCircut(circuit, theta)

	probs := make([]float64, c.NumberOfClasses)
	for i := 0; i < c.NumberOfClasses; i++ {
		_, probs[i] = circuit.Probability(1 << i)
	}
	return probs
}

func (c *Classifier) softmax(pred []float64, target float64) float64 {
	s := dataset.Sum(pred)
	return target / s
}

func (c *Classifier) crossEntropyLoss(prediction []float64, target []float64) float64 {
	cel := make([]float64, len(prediction))
	for i := 0; i < len(prediction); i++ {
		cel[i] = c.softmax(prediction, prediction[i])
	}
	e := 0.0
	for i := 0; i < len(cel); i++ {
		e += target[i] * math.Log(cel[i])
	}
	return -e
}

func (c *Classifier) loss(prediction []float64, target []float64) float64 {
	s := 0.0
	for i := 0; i < len(prediction); i++ {
		s += math.Pow(prediction[i]-target[i], 2)
	}
	return s
}

func (c *Classifier) gradient(X []float64, Y []float64) [][]float64 {

	delta := 10e-8

	grad := make([][]float64, len(c.theta))
	dtheta := make([][]float64, len(c.theta))
	dtheta2 := make([][]float64, len(c.theta))

	for i := 0; i < len(c.theta); i++ {
		dtheta[i] = dataset.CopyArray(c.theta[i])
		dtheta2[i] = dataset.CopyArray(c.theta[i])
	}

	for i := 0; i < len(c.theta); i++ {
		for j := 0; j < len(c.theta[i]); j++ {
			dtheta[i][j] += delta
			dtheta2[i][j] -= delta

			pred := c.quantiumNN(X, dtheta)
			pred2 := c.quantiumNN(X, dtheta2)

			grad[i] = append(grad[i], (c.crossEntropyLoss(pred, Y)-c.crossEntropyLoss(pred2, Y))/(delta*2))
		}
	}
	return grad
}

func (c *Classifier) Accuracy(X [][]float64, Y [][]float64) (float64, int, int) {
	cnt := 0
	for i := 0; i < len(X); i++ {
		pred := c.quantiumNN(X[i], c.theta)

		_, pIdx := dataset.Max(pred)
		_, yIdx := dataset.Max(Y[i])

		if pIdx == yIdx {
			cnt++
		}
	}
	return float64(cnt) / float64(len(X)), cnt, len(X)
}

func (c *Classifier) GetThetaParameters() [][]float64 {
	return c.theta
}
