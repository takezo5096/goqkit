# GoQKit

GoQKit is the quantum circuit simulator Golang module.
You can write a code and run the quantum calculation in Golang on your PC.


## How to use
A sample code for entangled two qbits.
````golang
package main

import (
	"fmt"
	"github.com/takezo5096/goqkit"
)

func main() {
	n := 2 //number of qbits
	circuit := goqkit.MakeQBitsCircuit(n)
	reg1 := circuit.AssignQBits(1)
	reg2 := circuit.AssignQBits(1)

	reg1.Write(0)
	reg2.Write(0)

	reg1.HadAll()
	reg2.Not(0x01, reg1.ToGlobalQBits(0x01)) //entangle

	r1 := reg1.Read(0x01)
	r2 := reg2.Read(0x01)

	circuit.PrintQBits()
	fmt.Println(r1, r2)
}
````
You'll see r1 and r2 are same value running this code because those qbits are entangled virtually in this simulator.

First of all, You should make the instance of a quantum circuit by MakeQBitsCircuit(N),
N is a number of qbits you want to use.

Then assign a register which has some qbits by AssignQBits(n), n is a number of qbits which this register has.
n must be under N - (a number of qbits include other regsiter's qbits)
If a circuit has 5 qbits(N=5), the register "a" might assign 2, "b" might assign 3 for example.
Obviously total assigned qbits must be <=N.

Finally, apply functions to qbits as you like.
If you want to see all qbit's value as vector in the circuit, PrintQBits() or PrintQBitsComplex() is useful for it. 

## More details
I still haven't prepared anything for the documents.
See the source codes and comments of goqkit package by using godoc or something.

And also, I created the visualizing tool for GoQKit to help you visualize your own quantum circuit and qbit's states.
 See [GoQKit Visualizer](https://github.com/takezo5096/goqkit_visualizer).

## License
**MIT License**

## Background
Now I'm working for web the media service provider,
I,ve been writing tons of code for a quarter of a century as a software engineer though,
not familiar with the quantum physics and the quantum calculation,
but I'm interested in that, because I want to understand especially the quantum calculation.
In 10 years, I might have no choice but to write a quantum calculation code from necessity as a programmer.
That's why I started to learn quantum calculation.

At first, I began to read some text books about the quantum physics and quantum caluclation to learn,
but I hadn't understood deeply cuz of the lack of my quantum physics knowledge.
So I decided to make a quantum circuit simulator to understand that using my programming skills.
For me, it's easy to implement basic quantum bit and gates like a "Not", "Hadamard", "Z" gate, etc,
because the circuit and bit operations are very simple and similar to classic computer's one I've already known.
Then I figured out that advanced physics and mathematics aren't always necessarily to learn the quantum calculation logics.
(Of course, learning those academic stuffs are still important to understand though).
The most difficult thing for me to learn the quantum calculation 
is just assembling those primitive gates to build the beneficial quantum algorithm.
I hope that I would come up with and use those kinds of algorithm to solve a hard problem in the world on a real quantum hardware in the future.

Try this simulator ! I would be glad if I could help you learn the quantum circuit even just a little.

## Author
https://github.com/takezo5096