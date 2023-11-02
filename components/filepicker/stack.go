package filepicker

type stack struct {
	data []int
}

func (s *stack) Push(i int) {
	s.data = append(s.data, i)
}

func (s *stack) Pop() int {
	res := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return res
}

func (s stack) Length() int {
	return len(s.data)
}
