package stack

type Stack[V any] interface {
	Push(value V)
	Pop() (value V, ok bool)
	Peek() (value V, ok bool)
	Get(index int) (value V, ok bool)

	Size() int
	IsEmpty() bool
}

type arrayStack[V any] struct {
	values []V
}

func (s *arrayStack[V]) Push(value V) {
	s.values = append(s.values, value)
}

func (s *arrayStack[V]) Pop() (value V, ok bool) {
	if len(s.values) == 0 {
		ok = false
		return
	}
	value = s.values[len(s.values)-1]
	s.values = s.values[0 : len(s.values)-1]
	return value, true
}

func (s *arrayStack[V]) Peek() (value V, ok bool) {
	if len(s.values) == 0 {
		ok = false
		return
	}
	value = s.values[len(s.values)-1]
	return value, true
}

func (s *arrayStack[V]) Get(index int) (value V, ok bool) {
	if index >= len(s.values) {
		ok = false
		return
	}
	value = s.values[index]
	return value, true
}

func (s *arrayStack[V]) Size() int {
	return len(s.values)
}

func (s *arrayStack[V]) IsEmpty() bool {
	return len(s.values) == 0
}

func New[V any]() Stack[V] {
	return &arrayStack[V]{
		values: make([]V, 0),
	}
}
