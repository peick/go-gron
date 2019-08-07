package gron

import (
	"strings"
)

type stack struct {
	elements []element
}

func newStack() *stack {
	return &stack{[]element{}}
}

func (s *stack) Empty() bool {
	return len(s.elements) == 0
}

func (s *stack) Push(e element) {
	s.elements = append(s.elements, e)
}

func (s *stack) Pop() {
	s.elements = s.elements[0 : len(s.elements)-1]
}

func (s *stack) Peek() element {
	if len(s.elements) == 0 {
		return nil
	}

	return s.elements[len(s.elements)-1]
}

func (s *stack) String() string {
	vx := make([]string, 0)
	for _, item := range s.elements {
		vx = append(vx, item.String())
	}

	return strings.Join(vx, "")
}
