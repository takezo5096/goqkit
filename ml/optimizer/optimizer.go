package optimizer

type Optimizer interface {
	Gradient(grad [][]float64) [][]float64
}
