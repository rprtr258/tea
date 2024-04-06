package hierachy

type Node[T any] struct {
	Value    T
	Children []Node[T]
}

type myNodeData[T any] struct {
	value           T
	depth           int
	lastChildrenIdx int
}

func remapTree[T any](node Node[T], depth int, dest *[]myNodeData[T]) {
	*dest = append(*dest, myNodeData[T]{
		value:           node.Value,
		depth:           depth,
		lastChildrenIdx: 0, // NOTE: filled later
	})
	cur := (*dest)[len(*dest)-1]
	for _, child := range node.Children {
		remapTree(child, depth+1, dest)
	}
	cur.lastChildrenIdx = len(*dest) - 1
}

// NOTE: readonly
type Hierachy[T any] struct {
	nodes    []myNodeData[T]
	selected int
}

func New[T any](tree Node[T]) *Hierachy[T] {
	var nodes []myNodeData[T]
	remapTree(tree, 0, &nodes)
	return &Hierachy[T]{
		nodes:    nodes,
		selected: 0,
	}
}

func (h *Hierachy[T]) sel() myNodeData[T] {
	return h.nodes[h.selected]
}

func (h *Hierachy[T]) Selected() T {
	return h.nodes[h.selected].value
}

func (h *Hierachy[T]) GoDown() {
	if h.sel().lastChildrenIdx > h.selected {
		h.selected++
	}
}

func (h *Hierachy[T]) GoUp() {
	depth := h.sel().depth
	if depth == 0 {
		return
	}

	// TODO: can be optimised by storing parent index
	for h.nodes[h.selected-1].depth >= depth {
		h.selected--
	}
}

func (h *Hierachy[T]) GoPrev() {
	if h.selected > 0 && h.nodes[h.selected-1].depth == h.sel().depth {
		h.selected--
	}
}

func (h *Hierachy[T]) GoNext() {
	if h.selected < len(h.nodes)-1 && h.nodes[h.selected+1].depth == h.sel().depth {
		h.selected++
	}
}

func (h *Hierachy[T]) GoPrevOrUp() {
	if h.selected > 0 {
		h.selected--
	}
}

func (h *Hierachy[T]) GoNextOrUp() {
	if h.selected < len(h.nodes)-1 {
		h.selected++
	}
}
