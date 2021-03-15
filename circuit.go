/*
qkit core package

to make the quantum circuit, register, and manipulate basic quantum gates
*/
package goqkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/takezo5096/goqkit/mat"
	"github.com/takezo5096/goqkit/util/queue"
	"log"
	"math"
	"math/cmplx"
	"math/rand"
	"os"
	"time"
)

/*
QBits Circuit which has all of qbits and some of base quantum gates
*/
type QBitsCircuit struct {
	//The number of all qbits in this circuit.
	QBitNumber uint
	//The vector of all qbits which require 2^QBitNumber number.
	RawQBits mat.Vector
	//Temporary queue to assign qbits
	qBitsQueue queue.Queue

	qBitRegisters []*Register

	operations []Operation

	printBuffer string
}

const (
	OperationTypeSpace  = "Sp"
	OperationTypeRead   = "R"
	OperationTypeWrite  = "W"
	OperationTypeHad    = "H"
	OperationTypePhase  = "P"
	OperationTypeRotate = "Ro"
	OperationTypeNot    = "N"
	OperationTypeSwap   = "S"
	OperationTypeX      = "X"
	OperationTypeY      = "Y"
	OperationTypeZ      = "Z"
)

type Operation struct {
	OpName             string    `json:"op_name"`
	RegisterName       int       `json:"register_name"`
	RegisterNameString string    `json:"register_name_string"`
	TargetQBit         uint      `json:"target_qbit"`
	ControlQBits       []uint    `json:"control_qbits"`
	SwapQBit           uint      `json:"swap_qbit"`
	Options            []float64 `json:"options"`
}

type DumpFormat struct {
	Message    string               `json:"message"`
	Operations []Operation          `json:"operations"`
	Registers  []DumpFormatRegister `json:"registers"`
	QBits      [][]float64          `json:"qbits"`
}
type DumpFormatRegister struct {
	NumberOfQBits int    `json:"number_of_qbits"`
	QBits         []uint `json:"qbits"`
	Shift         int    `json:"shift"`
	Name          string `json:"reg_name"`
}

const (
	PrintTypePolar   = iota
	PrintTypeComplex = iota
)

/*
Make a instance of a qbits circuit.
*/
func MakeQBitsCircuit(qBitNumber int) QBitsCircuit {
	var i uint
	// qbits need 2^qBitNumber
	i = 1 << qBitNumber
	v := mat.NewVector(i)
	v.Set(0, 1)

	qBitRegisters := make([]*Register, 0)

	qBitsQueue := queue.Queue{}
	for j := 0; j < qBitNumber; j++ {
		qbit := int(math.Pow(2, float64(j)))
		qBitsQueue.Enqueue(qbit)
	}
	return QBitsCircuit{RawQBits: v, QBitNumber: uint(qBitNumber), qBitsQueue: qBitsQueue, qBitRegisters: qBitRegisters}
}

/*
Return 0 or 1 value in terms of val of the bit position which is specified by i
*/
func refBit(val int, i uint) int {
	return (val >> i) & 1
}

/*
Assign qbits for the register

num: The number of qbits which you want to assign
*/
func (q *QBitsCircuit) AssignQBits(num int, name string) *Register {
	var qbits uint = 0
	var shift = int(q.QBitNumber) - q.qBitsQueue.Size()
	cnt := 0
	for i := 0; i < num; i++ {
		qbitIndex := q.qBitsQueue.Dequeue()
		qbits = qbits | uint(qbitIndex)
		cnt++
	}
	reg := Register{numberOfQBits: cnt, qBits: qbits, shift: shift, circuit: q, Name: name}
	q.qBitRegisters = append(q.qBitRegisters, &reg)
	return &reg
}

/*
Print all qbits in a register with Polar mode.
*/
func (q *QBitsCircuit) PrintQBits() {
	q.printQBitsImpl(PrintTypePolar, -1, -1)
}

/*
Print all qbits in a register with Complex mode.
*/
func (q *QBitsCircuit) PrintQBitsComplex() {
	q.printQBitsImpl(PrintTypeComplex, -1, -1)
}

func (q *QBitsCircuit) printQBitsImpl(printType int, start, end int) {
	if printType == PrintTypePolar {
		fmt.Printf("printing in polar mode.\n")
	} else {
		fmt.Printf("printing in complex mode.\n")
	}
	if start >= 0 && end >= 0 && start <= end {
		fmt.Printf("%d to %d vector elements\n", start, end)
	}
	if start > end {
		start = -1
		end = -1
	}

	v := q.RawQBits
	var i int

	fmt.Printf("%d: ", 0)
	for i = 0; i < int(v.N); i++ {
		if start >= 0 && i < start {
			continue
		}
		if end >= 0 && i >= end {
			fmt.Println("")
			if int(v.N)-1 != end {
				fmt.Printf("printed %d to %d vector elements. still more elements remains to display...", start, end)
			}
			break
		}
		if i != 0 && i%16 == 0 {
			fmt.Println("")
			fmt.Printf("%d: ", i)
		}
		e := v.Data[i]

		if printType == PrintTypePolar {
			r, theta := cmplx.Polar(e)
			if r > 0 {
				fmt.Printf("(\x1b[32m%.2f\x1b[0m,", r)
			} else if r < 0 {
				fmt.Printf("(\x1b[31m%.2f\x1b[0m,", r)
			} else {
				fmt.Printf("(%.2f,", r)
			}
			if theta > 0 {
				fmt.Printf("\x1b[32m%.2f\x1b[0m)", theta/math.Pi*180)
			} else if theta < 0 {
				fmt.Printf("\x1b[31m%.2f\x1b[0m)", theta/math.Pi*180)
			} else {
				fmt.Printf("%.2f)", theta/math.Pi*180)
			}
		} else {
			fmt.Printf("%.2f", e)
		}
	}
	fmt.Println("")
}

/*
Get target qbits(Array) specified by val

Find bit which is set to 1 in this val
*/
func (q *QBitsCircuit) GetQBits(val int) []uint {
	bits := make([]uint, 0)
	var i uint
	for i = 0; i < q.QBitNumber; i++ {
		r := refBit(val, i)
		//get an only qbit is turned on
		if r == 1 {
			qbit := uint(math.Pow(2, float64(i)))
			bits = append(bits, qbit)
		}
	}
	return bits
}

/*
Find pair of qbits will be applied to a gate.
*/
func (q *QBitsCircuit) GetQBitPairs(targetQBit uint) [][]uint {

	var pairs [][]uint

	n := q.RawQBits.N

	limit := targetQBit
	bk := n / 2
	step := limit + 1
	if step < 0 {
		return pairs
	}

	if targetQBit > bk {
		return pairs
	}

	minusCnt := limit

	var i uint
	for i = 0; i < n; i++ {

		if minusCnt == 0 {
			i += limit
			minusCnt = limit
		}
		if i+limit >= n {
			break
		}
		pair := []uint{i, i + limit}
		pairs = append(pairs, pair)

		minusCnt--

	}
	return pairs
}

/*
Read all qbits and return value in this circuit
*/
func (q *QBitsCircuit) Read() int {
	ret := 0
	var i uint
	for i = 0; i < q.QBitNumber; i++ {
		idx := 1 << i
		r := int(q.ReadQBit(uint(idx)))
		if r == 1 {
			ret = ret | idx
		}
		q.addOperation(OperationTypeRead, q.GetRegister(idx), idx, 0, 0, []float64{float64(r)})
	}
	return ret
}

/*
Read qbits specified by val and return val
*/
func (q *QBitsCircuit) ReadQBits(val int) int {
	ret := 0
	qbits := q.GetQBits(val)
	for _, qbit := range qbits {
		r := q.ReadQBit(qbit)
		if r == 1 {
			ret = ret | int(qbit)
		}
		q.addOperation(OperationTypeRead, q.GetRegister(int(qbit)), int(qbit), 0, 0, []float64{float64(r)})
	}
	return ret
}

func (q *QBitsCircuit) Probability(targetIndex uint) (float64, float64) {
	pairs := q.GetQBitPairs(targetIndex)

	var v0 float64
	var v1 float64

	v0 = 0
	v1 = 0
	v0q := make(map[uint]int, 0)
	v1q := make(map[uint]int, 0)
	for _, pair := range pairs {
		tmp0 := math.Pow(cmplx.Abs(q.RawQBits.At(pair[0])), 2)
		tmp1 := math.Pow(cmplx.Abs(q.RawQBits.At(pair[1])), 2)
		if tmp0 > 0 {
			v0q[pair[0]]++
		}
		if tmp1 > 0 {
			v1q[pair[1]]++
		}
		v0 += tmp0
		v1 += tmp1
	}

	prob0 := v0 / (v0 + v1)

	return prob0, 1.0 - prob0
}

/*
Read qbits specified by val and return val
*/
func (q *QBitsCircuit) ReadQBit(targetIndex uint) uint {

	pairs := q.GetQBitPairs(targetIndex)

	var v0 float64
	var v1 float64

	v0 = 0
	v1 = 0
	v0q := make(map[uint]int, 0)
	v1q := make(map[uint]int, 0)
	for _, pair := range pairs {
		tmp0 := math.Pow(cmplx.Abs(q.RawQBits.At(pair[0])), 2)
		tmp1 := math.Pow(cmplx.Abs(q.RawQBits.At(pair[1])), 2)
		if tmp0 > 0 {
			v0q[pair[0]]++
		}
		if tmp1 > 0 {
			v1q[pair[1]]++
		}
		v0 += tmp0
		v1 += tmp1
	}

	prob0 := v0 / (v0 + v1)

	rand.Seed(time.Now().UnixNano())

	var returnVal uint
	if prob0 > rand.Float64() {

		for _, pair := range pairs {
			if _, ok := v0q[pair[0]]; ok {
				q.RawQBits.Set(pair[0], cmplx.Sqrt(complex(1.0/float64(len(v0q)), 0)))
			}
			q.RawQBits.Set(pair[1], complex(0, 0))
		}
		returnVal = 0
	} else {
		for _, pair := range pairs {
			if _, ok := v1q[pair[1]]; ok {
				q.RawQBits.Set(pair[1], cmplx.Sqrt(complex(1.0/float64(len(v1q)), 0)))
			}
			q.RawQBits.Set(pair[0], complex(0, 0))
		}
		returnVal = 1
	}
	return returnVal
}

/*
Write the val to qbits in this circuit
*/
func (q *QBitsCircuit) Write(val int) {

	readResult := 0
	for _, qbit := range q.GetQBits(val) {
		r := q.ReadQBit(qbit)
		readResult = readResult | int(r)
	}

	if readResult != val {
		q.NotWithoutOp(val, 0)

		q.addOperation(OperationTypeWrite, q.GetRegister(val), val, 0, 0, nil)
	}
}

/*
Apply the unitary matrix to the vector of qbits
*/
func (q *QBitsCircuit) Unitary(val int, controlValue int, m *mat.Matrix) {
	targetQBits := q.GetQBits(val)

	for _, targetQBit := range targetQBits {
		pairs := q.GetQBitPairs(targetQBit)
		for _, pair := range pairs {
			if controlValue == 0 || int(pair[0])&controlValue == controlValue {
				qbit := mat.NewVector(2)
				qbit.Set(0, q.RawQBits.At(pair[0]))
				qbit.Set(1, q.RawQBits.At(pair[1]))
				newQBit := m.Dot(qbit)
				q.RawQBits.Set(pair[0], newQBit.At(0))
				q.RawQBits.Set(pair[1], newQBit.At(1))
			}
		}
	}
}

func (q QBitsCircuit) GetRegister(val int) *Register {
	var targetReg *Register
	for _, reg := range q.qBitRegisters {
		if reg.qBits&uint(val) != 0 {
			targetReg = reg
			break
		}
	}
	return targetReg
}

/*
Hadamard gate
*/
func (q *QBitsCircuit) Had(val int, controlValue int) {
	sqrt2 := 1.0 / complex(math.Sqrt(2), 0)
	m := mat.NewMatrix(2, 2)
	m.Set(0, 0, sqrt2)
	m.Set(0, 1, sqrt2)
	m.Set(1, 0, sqrt2)
	m.Set(1, 1, -sqrt2)

	q.Unitary(val, controlValue, &m)

	q.addOperation(OperationTypeHad, q.GetRegister(val), val, controlValue, 0, nil)
}

/*
Not gate
*/
func (q *QBitsCircuit) Not(val int, controlValue int) {
	m := mat.NewMatrix(2, 2)
	m.Set(0, 0, 0)
	m.Set(0, 1, 1)
	m.Set(1, 0, 1)
	m.Set(1, 1, 0)

	q.Unitary(val, controlValue, &m)

	q.addOperation(OperationTypeNot, q.GetRegister(val), val, controlValue, 0, nil)
}
func (q *QBitsCircuit) NotWithoutOp(val int, controlValue int) {
	m := mat.NewMatrix(2, 2)
	m.Set(0, 0, 0)
	m.Set(0, 1, 1)
	m.Set(1, 0, 1)
	m.Set(1, 1, 0)

	q.Unitary(val, controlValue, &m)
}

/*
Rotate X gate
*/
func (q *QBitsCircuit) RotX(val int, controlValue int, deg float64) {
	q.rotImpl(val, controlValue, deg, 0, 0)
}

/*
Rotate Y gate
*/
func (q *QBitsCircuit) RotY(val int, controlValue int, deg float64) {
	q.rotImpl(val, controlValue, 0, deg, 0)
}

/*
Rotate X gate
*/
func (q *QBitsCircuit) RotZ(val int, controlValue int, deg float64) {
	q.rotImpl(val, controlValue, 0, 0, deg)
}

/*
Rotate gate
*/
func (q *QBitsCircuit) rotImpl(val int, controlValue int, degX, degY, degZ float64) {

	thetaX := degX * (math.Pi / 180.0)
	thetaY := degY * (math.Pi / 180.0)
	thetaZ := degZ * (math.Pi / 180.0)

	var v00, v01, v10, v11 complex128

	if thetaX > 0 {
		v00 = complex(math.Cos(thetaX/2.0), 0)
		v01 = complex(0, -math.Sin(thetaX/2.0))
		v10 = complex(0, -math.Sin(thetaX/2.0))
		v11 = complex(math.Cos(thetaX/2.0), 0)
	}
	if thetaY > 0 {
		v00 = complex(math.Cos(thetaY/2.0), 0)
		v01 = complex(-math.Sin(thetaY/2.0), 0)
		v10 = complex(math.Sin(thetaY/2.0), 0)
		v11 = complex(math.Cos(thetaY/2.0), 0)
	}
	if thetaZ > 0 {
		v00 = cmplx.Exp(complex(0, -thetaZ/2.0))
		v01 = complex(0, 0)
		v10 = complex(0, 0)
		v11 = cmplx.Exp(complex(0, thetaZ/2.0))
	}

	m := mat.NewMatrix(2, 2)
	m.Set(0, 0, v00)
	m.Set(0, 1, v01)
	m.Set(1, 0, v10)
	m.Set(1, 1, v11)

	q.Unitary(val, controlValue, &m)

	q.addOperation(OperationTypeRotate, q.GetRegister(val), val, controlValue, 0, []float64{degX, degY, degZ})
}

/*
Phase Gate
*/
func (q *QBitsCircuit) Phase(val int, controlValue int, deg float64) {
	thetaZ := deg * (math.Pi / 180.0)

	m := mat.NewMatrix(2, 2)
	m.Set(0, 0, 1)
	m.Set(0, 1, 0)
	m.Set(1, 0, 0)
	m.Set(1, 1, cmplx.Exp(complex(0, thetaZ)))

	q.Unitary(val, controlValue, &m)

	q.addOperation(OperationTypePhase, q.GetRegister(val), val, controlValue, 0, []float64{deg})
}

/*
X Gate
*/
func (q *QBitsCircuit) X(val int, controlValue int) {
	m := mat.NewMatrix(2, 2)
	m.Set(0, 0, complex(0, 0))
	m.Set(0, 1, complex(1, 0))
	m.Set(1, 0, complex(1, 0))
	m.Set(1, 1, complex(0, 0))

	q.Unitary(val, controlValue, &m)

	q.addOperation(OperationTypeX, q.GetRegister(val), val, controlValue, 0, nil)
}

/*
Y Gate
*/
func (q *QBitsCircuit) Y(val int, controlValue int) {
	m := mat.NewMatrix(2, 2)
	m.Set(0, 0, complex(0, 0))
	m.Set(0, 1, complex(0, -1))
	m.Set(1, 0, complex(0, 1))
	m.Set(1, 1, complex(0, 0))

	q.Unitary(val, controlValue, &m)

	q.addOperation(OperationTypeY, q.GetRegister(val), val, controlValue, 0, nil)
}

/*
Z Gate
*/
func (q *QBitsCircuit) Z(val int, controlValue int) {
	m := mat.NewMatrix(2, 2)
	m.Set(0, 0, complex(1, 0))
	m.Set(0, 1, complex(0, 0))
	m.Set(1, 0, complex(0, 0))
	m.Set(1, 1, complex(-1, 0))

	q.Unitary(val, controlValue, &m)

	q.addOperation(OperationTypeZ, q.GetRegister(val), val, controlValue, 0, nil)
}

/*
Swap gate
*/
func (q *QBitsCircuit) Swap(targetVal int, swapVal int, controlValue int) {

	newControlValue := controlValue | swapVal
	q.NotWithoutOp(targetVal, newControlValue)

	newControlValue = controlValue | targetVal
	q.NotWithoutOp(swapVal, newControlValue)

	newControlValue = controlValue | swapVal
	q.NotWithoutOp(targetVal, newControlValue)

	q.addOperation(OperationTypeSwap, q.GetRegister(targetVal), targetVal, controlValue, swapVal, nil)
}

/*
Shift left
*/
func (q *QBitsCircuit) ShiftLeft(targetVal, controlVal int, shift int) {
	targetQBits := q.GetQBits(targetVal)

	tl := len(targetQBits)

	for i := 0; i < tl-1; i++ {
		tb := tl - i - 1
		pair1 := targetQBits[tb]
		pair2 := targetQBits[tb-shift]
		q.Swap(int(pair1), int(pair2), controlVal)
		if tb-shift == 0 {
			shift--
		}
	}
}

/*
QFT

Quantum version of Discrete Fourier transform(DFT)
*/
func (q *QBitsCircuit) QFT(val int) {
	idxs := q.GetQBits(val)

	for j := len(idxs) - 1; j >= 0; j-- {
		highestQbit := idxs[j]
		q.Had(int(highestQbit), 0)
		deg := -90.0
		for i := j - 1; i >= 0; i-- {
			q.Phase(int(highestQbit), int(idxs[i]), deg)
			deg = deg / 2.0
		}
	}
	for j := len(idxs) - 1; j >= len(idxs)/2; j-- {
		lowestQbit := idxs[len(idxs)-1-j]
		highestQbit := idxs[j]
		q.Swap(int(highestQbit), int(lowestQbit), 0)
	}
}

/*
Inversed QFT
*/
func (q *QBitsCircuit) InversedQFT(val int) {
	idxs := q.GetQBits(val)

	for j := len(idxs) - 1; j >= len(idxs)/2; j-- {
		lowestQbit := idxs[len(idxs)-1-j]
		highestQbit := idxs[j]
		q.Swap(int(highestQbit), int(lowestQbit), 0)

	}
	for j := 0; j < len(idxs); j++ {
		lowestQbit := idxs[j]
		q.Had(int(lowestQbit), 0)
		deg := 90.0
		for i := j + 1; i < len(idxs); i++ {
			q.Phase(int(lowestQbit), int(idxs[i]), deg)
			deg = deg / 2.0
		}
	}
}

/*
 To add non-operation (add a space)
*/
func (q *QBitsCircuit) OpSpace() {
	q.addOperation(OperationTypeSpace, q.GetRegister(1), 0, 0, 0, nil)
}

/*
Grover Algorithm
*/
func (q *QBitsCircuit) Grover(val int) {

	q.addOperation(OperationTypeSpace, q.GetRegister(val), 0, 0, 0, nil)

	q.Had(val, 0)
	q.Not(val, 0)

	controlVal := 0
	idxs := q.GetQBits(val)
	for i, idx := range idxs {
		if i != 0 {
			controlVal = controlVal | int(idx)
		}
	}
	//fmt.Println("Grover", idxs[0], controlVal)
	q.Phase(int(idxs[0]), controlVal, 180)
	q.Not(val, 0)
	q.Had(val, 0)
}

/*
Add val
*/
func (q *QBitsCircuit) Add(rangeValue int, val int, controlVal int) {
	if val < 0 {
		q.Subtract(rangeValue, -val, controlVal)
		return
	}

	controlIdxs := q.GetQBits(controlVal)

	idxes := q.GetQBits(val)

	for i := 0; i < len(idxes); i++ {
		idx := idxes[i]
		newRangeVal := 0
		for _, ridx := range q.GetQBits(rangeValue) {
			if ridx >= idx {
				newRangeVal = newRangeVal | int(ridx)
			}
		}
		newControlValue := 0
		if controlVal != 0 {
			newControlValue = int(controlIdxs[i])
		}
		q.addImpl(newRangeVal, newControlValue)
	}

}
func (q *QBitsCircuit) addImpl(rangeValue int, controlValue int) {
	cidxes := q.GetQBits(rangeValue)
	newControlVal := rangeValue | controlValue
	for i := len(cidxes) - 1; i >= 0; i-- {
		cidx := cidxes[i]
		newControlVal = newControlVal ^ int(cidx)
		q.Not(int(cidx), newControlVal)
	}
}

/*
Subtract val
*/
func (q *QBitsCircuit) Subtract(rangeValue int, val int, controlVal int) {
	if val < 0 {
		q.Add(rangeValue, -val, controlVal)
		return
	}

	controlIdxs := q.GetQBits(controlVal)

	idxes := q.GetQBits(val)

	for i := len(idxes) - 1; i >= 0; i-- {
		idx := idxes[i]
		newRangeVal := 0
		for _, ridx := range q.GetQBits(rangeValue) {
			if ridx >= idx {
				newRangeVal = newRangeVal | int(ridx)
			}
		}
		newControlValue := 0
		if controlVal != 0 {
			newControlValue = int(controlIdxs[i])
		}
		q.subtractImpl(newRangeVal, newControlValue)
	}

}
func (q *QBitsCircuit) subtractImpl(controlVal int, controlValue int) {

	cidxes := q.GetQBits(controlVal)
	newControlVal := controlValue
	for _, cidx := range cidxes {
		q.Not(int(cidx), newControlVal)
		newControlVal = newControlVal | int(cidx)
	}
}

func (q *QBitsCircuit) addOperation(opName string, reg *Register, target int, control int, swap int, options []float64) {
	controls := q.GetQBits(control)

	for _, t := range q.GetQBits(target) {
		var op Operation
		if len(controls) > 0 {
			op = Operation{OpName: opName, RegisterName: 1 << reg.shift, RegisterNameString: reg.Name, TargetQBit: t, ControlQBits: controls, SwapQBit: uint(swap), Options: options}
		} else {
			op = Operation{OpName: opName, RegisterName: 1 << reg.shift, RegisterNameString: reg.Name, TargetQBit: t, ControlQBits: nil, SwapQBit: uint(swap), Options: options}
		}
		q.operations = append(q.operations, op)
	}
	if opName == OperationTypeSpace {
		op := Operation{OpName: opName, RegisterName: 1 << reg.shift, RegisterNameString: reg.Name, TargetQBit: 0, ControlQBits: nil, SwapQBit: uint(swap), Options: options}
		q.operations = append(q.operations, op)
	}
}

func (q *QBitsCircuit) GetOperations() []Operation {
	return q.operations
}

func (q *QBitsCircuit) GetPrintBuffer() string {
	return q.printBuffer
}

/*
Put string into a buffer.
*/
func (q *QBitsCircuit) PrintBuffer(format string, a ...interface{}) {
	q.printBuffer += fmt.Sprintf(format, a...)
}

/*
Put string into a buffer with a newline.
*/
func (q *QBitsCircuit) PrintBufferln(format string, a ...interface{}) {
	q.PrintBuffer(format+"\n", a...)
}

/*
Dump a buffer which includes all operations, qbits states and a print buffer with json format.
*/
func (q QBitsCircuit) DumpAll() string {
	msg := q.GetPrintBuffer()
	ops := q.GetOperations()

	qbits := make([][]float64, 0)
	for _, qbit := range q.RawQBits.Data {
		tmp := make([]float64, 2)
		r, theta := cmplx.Polar(qbit)
		tmp[0] = r
		tmp[1] = theta
		qbits = append(qbits, tmp)
	}

	var registers []DumpFormatRegister
	for i, reg := range q.qBitRegisters {

		newReg := DumpFormatRegister{}
		newReg.Shift = reg.shift
		newReg.NumberOfQBits = reg.numberOfQBits
		newReg.QBits = q.GetQBits(int(reg.GetQBits()))
		if reg.Name != "" {
			newReg.Name = reg.Name
		} else {
			newReg.Name = fmt.Sprintf("Reg%d", i+1)
		}
		registers = append(registers, newReg)
	}

	df := DumpFormat{Message: msg, Operations: ops, Registers: registers, QBits: qbits}

	r, _ := json.Marshal(df)
	out := new(bytes.Buffer)
	// format json
	json.Indent(out, r, "", "    ")

	return out.String()
}

/*
Dump a buffer which includes all operations, qbits states and a print buffer to a json format file.
*/
func (q *QBitsCircuit) FileDumpAll(path string) {

	s := q.DumpAll()

	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Write(([]byte)(s))
}
