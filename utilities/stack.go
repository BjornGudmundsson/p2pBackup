package utilities

//GenericStack is a generic
//implementation of a stack
type GenericStack struct {
	stack []interface{}
	index int16
}

//Push pushes a value on top of the stack.
func (s *GenericStack) Push(v interface{}) {
	l := len(s.stack)
	if s.index == int16(l-1) {
		buf := make([]interface{}, 2*l)
		copy(buf[:l], s.stack)
		s.stack = buf
	}
	s.index++
	s.stack[s.index] = v
}

//Pop takes the top element of the stack.
func (s *GenericStack) Pop() interface{} {
	if s.index == -1 {
		return nil
	}
	t := s.stack[s.index]
	s.index--
	return t
}

//IsEmpty returns whether the given stack is empty or not
func (s *GenericStack) IsEmpty() bool {
	return s.index == -1
}

//NewStack gives a new, empty stack
func NewStack() Stack {
	return &GenericStack{
		stack: make([]interface{}, 2),
		index: -1,
	}
}
