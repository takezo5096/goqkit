package optimizer

import (
	"math"
)

type Adam struct {
	v [][]float64
	m [][]float64

	t float64

	LearningRate float64
}

func NewAdam(numberOfLayer int, numberOfParam int, learningRate float64) *Adam {
	a := new(Adam)

	a.v = make([][]float64, numberOfLayer)
	a.m = make([][]float64, numberOfLayer)
	for i := 0; i < numberOfLayer; i++ {
		a.v[i] = make([]float64, numberOfParam)
		a.m[i] = make([]float64, numberOfParam)
	}
	a.t = 1
	a.LearningRate = learningRate

	return a
}

func (a *Adam) Gradient(grad [][]float64) [][]float64 {
	p1 := 0.9
	p2 := 0.999
	e := 10e-8
	n := a.LearningRate

	delta := make([][]float64, len(grad))

	for j := 0; j < len(grad); j++ {
		delta[j] = make([]float64, len(grad[j]))
		for i := 0; i < len(grad[j]); i++ {
			if a.t == 1 {
				a.m[j][i] = grad[j][i]
				a.v[j][i] = grad[j][i] * grad[j][i]
			} else {
				a.m[j][i] = p1*a.m[j][i] + (1-p1)*grad[j][i]
				a.v[j][i] = p2*a.v[j][i] + (1-p2)*grad[j][i]*grad[j][i]
			}
			mh := a.m[j][i] / (1 - math.Pow(p1, a.t))
			vh := a.v[j][i] / (1 - math.Pow(p2, a.t))

			delta[j][i] = n / (math.Sqrt(vh + e)) * mh
		}
	}
	a.t += 1
	return delta
}
