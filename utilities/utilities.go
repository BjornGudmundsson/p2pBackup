package utilities

//Stack is a generic LIFO queue
//Stack uses interfaces for more generality
type Stack interface {
	Pop() interface{}
	Push(v interface{})
	IsEmpty() bool
}
