package hierachy

type Node[T any] struct {
	Value    T
	Children []Node[T]
}

type myNodeData[T any] struct {
	value           T
	depth           int
	parentIdx       int // -1 for root
	lastChildrenIdx int
	collapsed       bool
}

func remapTree[T any](node Node[T], parentIdx int, dest *[]myNodeData[T]) int {
	depth := 0
	if parentIdx != -1 {
		depth = (*dest)[parentIdx].depth + 1
	}

	*dest = append(*dest, myNodeData[T]{
		value:           node.Value,
		depth:           depth,
		parentIdx:       parentIdx,
		lastChildrenIdx: -1,    // NOTE: filled later
		collapsed:       false, // NOTE: filled later
	})
	ownIdx := len(*dest) - 1
	count := 0 // TODO: count actually is just len(*dest)-ownIdx after the loop, but I couldn't make it work vahui
	for _, child := range node.Children {
		count += remapTree(child, ownIdx, dest)
	}
	(*dest)[ownIdx].lastChildrenIdx = ownIdx + count
	(*dest)[ownIdx].collapsed = count == 0
	return 1 + count
}

// NOTE: readonly
type Hierachy[T any] struct {
	nodes    []myNodeData[T]
	selected int
}

func New[T any](tree Node[T]) *Hierachy[T] {
	var nodes []myNodeData[T]
	remapTree(tree, -1, &nodes)
	return &Hierachy[T]{
		nodes:    nodes,
		selected: 0,
	}
}

func (h *Hierachy[T]) sel() *myNodeData[T] {
	return &h.nodes[h.selected]
}

func (h *Hierachy[T]) Selected() T {
	return h.sel().value
}

func (h *Hierachy[T]) IsCollapsed() bool {
	return h.sel().collapsed
}

func (h *Hierachy[T]) ToggleCollapsed() {
	if h.sel().lastChildrenIdx == h.selected {
		// leaf nodes cannot be collapsed
		return
	}

	h.sel().collapsed = !h.sel().collapsed
}

func (h *Hierachy[T]) GoDown() {
	if h.sel().lastChildrenIdx > h.selected {
		h.selected++
	}
}

func (h *Hierachy[T]) GoUp() {
	parentIdx := h.sel().parentIdx
	if parentIdx == -1 {
		return
	}

	h.selected = parentIdx
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
	if h.selected == 0 {
		return
	}

	h.selected--

	// if at least one parent is collapsed, go to highest one
	highestCollapsedParentIdx := -1
	for i := h.sel().parentIdx; i != -1; i = h.nodes[i].parentIdx {
		if h.nodes[i].collapsed {
			highestCollapsedParentIdx = i
		}
	}
	if highestCollapsedParentIdx != -1 {
		h.selected = highestCollapsedParentIdx
	}
}

func (h *Hierachy[T]) GoNextOrUp() {
	if h.sel().collapsed {
		if newSelected := h.sel().lastChildrenIdx + 1; newSelected < len(h.nodes) {
			h.selected = newSelected
		}
	} else if h.selected < len(h.nodes)-1 {
		h.selected++
	}
}

type IterItem[T any] struct {
	Value       T
	Depth       int
	IsSelected  bool
	HasChildren bool
	IsCollapsed bool
}

func (h *Hierachy[T]) Iter(yield func(IterItem[T]) bool) {
	for i := 0; i < len(h.nodes); i++ {
		node := h.nodes[i]
		if !yield(IterItem[T]{
			Value:       node.value,
			Depth:       node.depth,
			IsSelected:  i == h.selected,
			HasChildren: node.lastChildrenIdx != i,
			IsCollapsed: node.collapsed,
		}) {
			return
		}

		if node.collapsed {
			i = node.lastChildrenIdx
		}
	}
}
