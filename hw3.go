package main

import "fmt"

type Stack struct {
	elements []int
}

func (s *Stack) Push(value int) {
	s.elements = append(s.elements, value)
}

func (s *Stack) Pop() (int, bool) {
	if s.IsEmpty() {
		return 0, false
	}

	lastIndex := len(s.elements) - 1
	value := s.elements[lastIndex]
	s.elements = s.elements[:lastIndex]
	return value, true
}

func (s *Stack) IsEmpty() bool {
	return len(s.elements) == 0
}

func (s *Stack) Size() int {
	return len(s.elements)
}

func (s *Stack) Clear() {
	s.elements = make([]int, 0)
}

func main() {
	fmt.Println("My stack realization!")
}
