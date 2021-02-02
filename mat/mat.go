/*
mat is minimum complex vector and matrix package
 */
package mat

type Matrix struct {
	Rows uint
	Cols uint
	Data [][]complex128
}

type Vector struct {
	N uint
	Data []complex128
}

func NewMatrix(r, c uint) Matrix {
	var data = make([][]complex128, r)
	var i,j uint
	for i=0; i<r; i++ {
		data[i] = make([]complex128, c)
		for j=0; j<c; j++ {
			data[i][j] = complex128(0)
		}
	}
	return Matrix{Rows: r, Cols: c, Data: data}
}

func NewIMatrix(r, c uint) Matrix {
	m := NewMatrix(r, c)
	var i,j uint
	for i=0; i<m.Rows; i++ {
		for j=0; j<m.Cols; j++ {
			if i == j {
				m.Set(i, j, 1)
			}
		}
	}
	return m
}

func NewVector(n uint) Vector {
	var data []complex128
	var i uint
	for i=0; i<n; i++ {
		data = append(data, complex128(0))
	}
	return Vector{N:n, Data:data}
}

/*
func (v *Vector) Print() {
	var i uint
	fmt.Printf("%d: ", 0)
	for i=0; i<v.N; i++ {
		if i!=0 && i%16 == 0 {
			fmt.Println("")
			fmt.Printf("%d: ", i)
		}
		e := v.Data[i]
		fmt.Printf("%.1f", e)
	}
	fmt.Println("")
}

func (m *Matrix) Print() {
	var i,j uint
	for i=0; i<m.Rows; i++ {
		for j=0; j<m.Cols; j++ {
			fmt.Printf("%.1f ", m.At(i, j))
		}
		fmt.Println("")
	}
}
*/

func (v *Vector) Set(n uint, val complex128) {
	v.Data[n] = val
}
func (v *Vector) SetAll(val complex128) {
	var i uint
	for i=0; i<v.N; i++ {
		v.Set(i, val)
	}
}
func (v *Vector) At(n uint) complex128 {
	return v.Data[n]
}

func (m *Matrix) Set(r, c uint, val complex128) {
	m.Data[r][c] = val
}
func (m *Matrix) SetAll(val complex128) {
	var i,j uint
	for i=0; i<m.Rows; i++ {
		for j=0; j<m.Cols; j++ {
			m.Set(i, j, val)
		}
	}
}
func (m *Matrix) At(r, c uint) complex128 {
	return m.Data[r][c]
}

func (m *Matrix) Dot(x Vector) Vector {
	y := NewVector(m.Rows)

	var i,j uint
	for i=0; i<m.Rows; i++ {
		var tmp complex128 = 0
		for j=0; j<m.Cols; j++ {
			tmp += m.At(i, j) * x.At(j)
		}
		y.Set(i, tmp)
	}
	return y
}

func (m *Matrix) Copy() Matrix {
	var newM = NewMatrix(m.Rows, m.Cols)
	var i, j uint
	for i=0; i<m.Rows; i++ {
		for j=0; j<m.Cols; j++ {
			newM.Set(i, j, m.At(i, j))
		}
	}
	return newM
}

func (m *Matrix) RowExchange(r1, r2 uint) {

	var j uint
	for j=0; j<m.Cols; j++ {
		tmp1 := m.At(r1, j)
		tmp2 := m.At(r2, j)
		m.Set(r1, j, tmp2)
		m.Set(r2, j, tmp1)
	}
}
