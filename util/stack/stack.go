package stack

import "fmt"

// Stackは[]intのエイリアス
type Stack []int

// Push adds an element
func (s *Stack) Push(v int) {
	*s = append(*s, v)
}

// Pop removes the top element and return it
func (s *Stack) Pop() (int, error) {
	if s.Empty() {
		return 0, fmt.Errorf("stack is empty")
	}

	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v, nil
}

// Peek returns the top value
func (s *Stack) Peek() (int, error) {
	if s.Empty() {
		return 0, fmt.Errorf("stack is empty")
	}

	return (*s)[len(*s)-1], nil
}

// Size returns the length of stack
func (s *Stack) Size() int {
	return len(*s)
}

// Empty returns true when stack is empty
func (s *Stack) Empty() bool {
	return s.Size() == 0
}

// NewStack generates stack
func NewStack() *Stack {
	s := new(Stack)
	return s
}
