package goqkit

/*
QBits register which has same of qbits and can apply the many quantum gates
*/
type Register struct {
	//Number of QBits in this register.
	numberOfQBits int
	//QBits value. if this register has 0x01,0x02 qbits, then QBits will be 0x03
	qBits uint
	//How many shift from local to global value in circuit.
	shift int
	//Pointer of the circuit.
	circuit *QBitsCircuit
}

/*
Do nothing.
*/
func (reg *Register) Nothing() {

}

/*
Return number of qbits in this register.
*/
func (reg *Register) NumberOfQBits() int {
	return reg.numberOfQBits
}

/*
Return qbits value.

Example: if this register has 0x01, 0x02, 0x04 qbits, return 6.
*/
func (reg *Register) GetQBits() uint {
	return reg.qBits
}

/*
Return the qbits value from the point of global view.

The register use local qbits value on basis,
so depending on the situation, transfer from local value in this register to global one in global circuit.

Example: if shift is 4 and this register has 0x01 local value, this function will return 0x08 by shifting 4.
*/
func (reg *Register) ToGlobalQBits(val int) int {
	return val << reg.shift
}

/*
Read all qbits value in this register and return local integer value.
*/
func (reg *Register) ReadAll() int {
	r := reg.circuit.ReadQBits(int(reg.GetQBits()))
	//back to local
	return r >> reg.shift
}

/*
Read the qbits specified as val and return local integer value.

val: local qbits value
*/
func (reg *Register) Read(val int) int {
	qbits := reg.ToGlobalQBits(val)
	r := reg.circuit.ReadQBits(qbits)
	//back to local
	return r >> reg.shift
}

/*
Write the qbits specified as val.

val: local qbits value
*/
func (reg *Register) Write(val int) {
	qbits := reg.ToGlobalQBits(val)
	reg.circuit.Write(qbits)
}

/*
Apply Not Gate to all qbits in this register
*/
func (reg *Register) NotAll() {
	reg.circuit.Not(int(reg.qBits), 0)
}

/*
Appy Not Gate to the specified value with control

val: local qbits value

control: global control qbits value
*/
func (reg *Register) Not(val int, control int) {
	reg.circuit.Not(reg.ToGlobalQBits(val), control)
}

/*
Apply Hadamard Gate to all qbits in this register
*/
func (reg *Register) HadAll() {
	reg.circuit.Had(int(reg.qBits), 0)
}

/*
Appy Hadamard Gate to the specified value with control

val: local qbits value

control: global control qbits value
*/
func (reg *Register) Had(val int, control int) {
	qbits := reg.ToGlobalQBits(val)
	reg.circuit.Had(qbits, control)
}

/*
Apply Rotate Gate to all qbits in this register

deg: degree to rotate

*/
func (reg *Register) RotXAll(deg float64) {
	reg.circuit.RotX(int(reg.qBits), 0, deg)
}

/*
Apply Rotate Gate to all qbits in this register

deg: degree to rotate

*/
func (reg *Register) RotYAll(deg float64) {
	reg.circuit.RotY(int(reg.qBits), 0, deg)
}

/*
Apply Rotate Gate to all qbits in this register

deg: degree to rotate

*/
func (reg *Register) RotZAll(deg float64) {
	reg.circuit.RotZ(int(reg.qBits), 0, deg)
}

/*
Apply Rotate Gate to the value which specified val with control qbits.

val: local qbits value

control: global control qbits value

deg: degree to rotate
*/
func (reg *Register) RotX(val int, control int, deg float64) {
	qbits := reg.ToGlobalQBits(val)
	reg.circuit.RotX(qbits, control, deg)
}

/*
Apply Rotate Gate to the value which specified val with control qbits.

val: local qbits value

control: global control qbits value

deg: degree to rotate
*/
func (reg *Register) RotY(val int, control int, deg float64) {
	qbits := reg.ToGlobalQBits(val)
	reg.circuit.RotY(qbits, control, deg)
}

/*
Apply Rotate Gate to the value which specified val with control qbits.

val: local qbits value

control: global control qbits value

deg: degree to rotate
*/
func (reg *Register) RotZ(val int, control int, deg float64) {
	qbits := reg.ToGlobalQBits(val)
	reg.circuit.RotZ(qbits, control, deg)
}

/*
Apply X Gate to the value with control qbits.

val: local qbits value

control: global control qbits value
*/
func (reg *Register) XAll() {
	reg.circuit.X(int(reg.qBits), 0)
}
func (reg *Register) X(val int, control int) {
	reg.X(val, control)
}

/*
Apply Y Gate to the value with control qbits.

val: local qbits value

control: global control qbits value
*/
func (reg *Register) YAll() {
	reg.circuit.Y(int(reg.qBits), 0)
}
func (reg *Register) Y(val int, control int) {
	reg.Y(val, control)
}

/*
Apply Z Gate to the value with control qbits.

val: local qbits value

control: global control qbits value

*/
func (reg *Register) ZAll() {
	reg.circuit.Z(int(reg.qBits), 0)
}
func (reg *Register) Z(val int, control int) {
	reg.Z(val, control)
}

func (reg *Register) PhaseAll(deg float64) {
	reg.circuit.Phase(int(reg.qBits), 0, deg)
}
func (reg *Register) Phase(val int, control int, deg float64) {
	qbits := reg.ToGlobalQBits(val)
	reg.circuit.Phase(qbits, control, deg)
}

/*
Apply Swap Gate to all qbits in this register

targetVal: local target qbits value

swapVal: local swap target qbits value

control: global control qbits value

*/
func (reg *Register) Swap(targetVal int, swapVal int, control int) {
	tqbits := reg.ToGlobalQBits(targetVal)
	reg.circuit.Swap(tqbits, swapVal, control)

}

/*
Shift the qbits value in this register with shift number.

control: global control qbits value

numShift: number of shift
*/
func (reg *Register) ShiftLeft(control int, numShift int) {
	reg.circuit.ShiftLeft(int(reg.GetQBits()), control, numShift)
}

/*
Quantum version of Discrete Fourier transform(DFT)

Apply all qbits in this register.
*/
func (reg *Register) QFT() {
	reg.circuit.QFT(int(reg.GetQBits()))
}

/*
Quantum version of Inversed Discrete Fourier transform(InvDFT)

Apply all qbits in this register.
*/
func (reg *Register) InversedQFT() {
	reg.circuit.InversedQFT(int(reg.GetQBits()))
}

/*
Subtract the value with control qbits.

val: local qbits value

control: global control qbits value

Example: Subtract(3, 0)
*/
func (reg *Register) Subtract(val int, control int) {
	reg.circuit.Subtract(int(reg.GetQBits()), val, control)
}

/*
Subtract a register

registerB: the register to subtract to this register
*/
func (reg *Register) SubtractRegister(registerB Register) {
	reg.Subtract(int(registerB.GetQBits()>>registerB.shift), int(registerB.GetQBits()))
}

/*
Add the value with control qbits.

val: local qbits value

control: global control qbits value

Example: Add(3, 0)
*/
func (reg *Register) Add(val int, control int) {
	reg.circuit.Add(int(reg.GetQBits()), val, control)
}

/*
Add a register

registerB: the register to add to this register
*/
func (reg *Register) AddRegister(registerB *Register) {
	reg.Add(int(registerB.GetQBits()>>registerB.shift), int(registerB.GetQBits()))
}
