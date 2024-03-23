package filepicker

type stack struct {
	data []int
}

func Pop(data []int) ([]int, int) {
	return data[:len(data)-1], data[len(data)-1]
}
