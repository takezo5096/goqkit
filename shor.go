package goqkit

import (
	"fmt"
	"math"
)

// Here are some values of N to try:
// 15, 21, 35, 39, 51, 55, 69, 77, 85, 87, 91, 93, 95, 111, 115, 117,
// 119, 123, 133, 155, 187, 203, 221, 247, 259, 287, 341, 451

// Larger numbers require more bits of precision.
// N = 15    precision_bits >= 4
// N = 21    precision_bits >= 5
// N = 35    precision_bits >= 6
// N = 123   precision_bits >= 7
// N = 341   precision_bits >= 8  time: about 6 seconds
// N = 451   precision_bits >= 9  time: about 23 seconds

/*
Implement of Shor algorithm

N: The number which you want to do prime factorization

precisionBits: The number of QBits used as Precision

coprime: Fixed value in this Circuit. must be 2
 */
func ShorCircuit(N, precisionBits, coprime int) (*QBitsCircuit, int, []int, []int, error) {

	var qc, readResult, repeatPeriod, factor, err = Shor(N, precisionBits, coprime)

	return qc, readResult, repeatPeriod, factor, err
}

func Shor(N, precisionBits, coprime int) (*QBitsCircuit, int, []int, []int, error) {
	var qc, readResult, repeatPeriod = ShorQPU(N, precisionBits, coprime) // quantum part
	//var repeat_period = ShorNoQPU(N, precision_bits, coprime) // quantum part(but no QPU)
	var factors = ShorLogic(N, repeatPeriod, coprime) // classical part
	fact, err := checkResult(N, factors)
	return qc, readResult, repeatPeriod, fact, err
}

func gcd(a, b int) int {
	// return the greatest common divisor of a,b
	for b != 0 {
		var m = a % b
		a = b
		b = m
	}
	return a
}

func checkResult(N int, factorCandidates [][]int) ([]int, error) {
	for i := 0; i < len(factorCandidates); i++ {
		var factors = factorCandidates[i]
		if factors[0] * factors[1] == N {
			if factors[0] != 1 && factors[1] != 1 {
				// Success!
				return factors, nil
			}
		}
	}
	// Failure
	return nil, fmt.Errorf("failure: no non-trivial factors were found")
}

func ShorLogic(N int, repeatPeriodCandidates []int, coprime int) [][]int {
	factorCandidates := make([][]int, 0)
	for i:= 0; i < len(repeatPeriodCandidates); i++ {
		var repeatPeriod = repeatPeriodCandidates[i]
		// Given the repeat period, find the actual factors
		var ar2 = math.Pow(float64(coprime), float64(repeatPeriod) / 2.0)
		var factor1 = gcd(N, int(ar2) - 1)
		var factor2 = gcd(N, int(ar2) + 1)
		factorCandidates = append(factorCandidates, []int{factor1, factor2})
	}
	return factorCandidates
}

func ShorNoQPU(N, precisionBits, coprime int) []int {
	// Classical replacement for the quantum part of Shor
	var work = 1
	var maxLoops = int(math.Pow(2, float64(precisionBits)))
	for iter := 0; iter < maxLoops; iter++ {
		work = (work * coprime) % N
		if work == 1 { // found the repeat
			return []int{iter + 1}
		}
	}
	return nil
}

func ShorQPU(N, precisionBits, coprime int) (*QBitsCircuit, int, []int) {

	// Quantum part of Shor's algorithm
	// For this implementation, the coprime must be 2.
	coprime = 2

	// For some numbers (like 15 and 21) the "mod" in a^xmod(N)
	// is not needed, because a^x wraps neatly around. This makes the
	// code simpler, and much easier to follow.
	//if N == 15 || N == 21 {
	return ShorQPUWithoutModulo(N, precisionBits)
	//}
	//else{
	//	return ShorQPU_WithModulo(N, precision_bits, coprime)
	//}
	//return nil
}

func readUnsigned(qreg *Register) int {
	var value = qreg.ReadAll()
	return value & ((1 << qreg.NumberOfQBits()) - 1)
}

func ShorQPUWithoutModulo(N, precisionBits int) (*QBitsCircuit, int, []int) {
	var NBits = 1
	for (1 << NBits) < N {
		NBits++
	}
	if N != 15 { // For this implementation, numbers other than 15 need an extra bit
		NBits++
	}
	var totalBits = NBits + precisionBits
	//fmt.Println("total_bits:", totalBits, NBits, precisionBits)

	// Set up the QPU and the working registers
	qc := MakeQBitsCircuit(totalBits)
	var num = qc.AssignQBits(NBits)
	var precision = qc.AssignQBits(precisionBits)

	num.Write(1)
	precision.Write(0)
	precision.HadAll()

	// Perform 2^x for all possible values of x in superposition
	for iter := 0; iter < precisionBits; iter++ {
		var numShifts = 1 << iter
		var condition = precision.ToGlobalQBits(numShifts)
		numShifts %= num.NumberOfQBits()
		if numShifts == 0 {
			continue
		}
		num.ShiftLeft(condition, numShifts)
	}
	// Perform the QFT
	precision.QFT()

	//qc.PrintQBits(PrintTypePolar, 0, 1000)

	var readResult = readUnsigned(precision)
	var repeatPeriodCandidates = estimateNumSpikes(readResult, 1 << precisionBits)

	return &qc, readResult, repeatPeriodCandidates
}

func estimateNumSpikes(spike, range1 int) []int {
	if float64(spike) < float64(range1) / 2.0 {
		spike = range1 - spike
	}
	var bestError = 1.0
	var e0 = 0.0
	var e1 = 0.0
	var e2 = 0.0
	var actual = float64(spike) / float64(range1)
	var candidates = make([]int, 0)
	for denom := 1; denom < spike; denom++{
		var numerator = math.Round(float64(denom) * actual)
		var estimated = numerator / float64(denom)
		var err = math.Abs(estimated - actual)
		e0 = e1
		e1 = e2
		e2 = err
		// Look for a local minimum which beats our current best error
		if e1 <= bestError && e1 < e0 && e1 < e2 {
			var repeatPeriod = denom - 1
			candidates = append(candidates, repeatPeriod)
			bestError = e1
		}
	}
	return candidates
}